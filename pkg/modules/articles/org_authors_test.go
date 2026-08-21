package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Главное свойство всей функции: имя, которое можно вписать самому, не
// показывается нигде, пока его не подтвердил человек. Иначе площадка
// превращается в машину для выдачи себя за акимат, а одна поддельная
// публикация стоит дороже сотни настоящих.
func TestUnverifiedOrganisationIsShownNowhere(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("orgfake@example.com", "Parol123!")
	store := NewOrgStore(app.pool)
	if err := store.Apply(context.Background(), authorID,
		OrgAuthor{Name: "Акимат г. Костанай", Kind: "akimat"}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Заявка есть, но подтверждения нет.
	org, err := store.ByUser(context.Background(), authorID)
	if err != nil || org == nil {
		t.Fatalf("ByUser: %v, %v", org, err)
	}
	if org.Verified() {
		t.Fatal("свежая заявка сразу оказалась подтверждённой")
	}
	got, err := store.VerifiedByUser(context.Background(), authorID)
	if err != nil {
		t.Fatalf("VerifiedByUser: %v", err)
	}
	if got != nil {
		t.Fatal("непроверенная организация вернулась как подтверждённая")
	}

	// И на странице статьи её тоже нет.
	_, slug := app.seedArticle(authorID, "published")
	w := app.do(http.MethodGet, "/read/"+slug, nil)
	if strings.Contains(w.Body.String(), "Акимат г. Костанай") {
		t.Error("непроверенное название организации попало в подпись статьи")
	}

	// После подтверждения — появляется, и рядом стоит человек.
	if err := store.SetStatus(context.Background(), authorID, orgVerified, "", nil); err != nil {
		t.Fatalf("SetStatus: %v", err)
	}
	w = app.do(http.MethodGet, "/read/"+slug, nil)
	body := w.Body.String()
	if !strings.Contains(body, "Акимат г. Костанай") {
		t.Error("подтверждённая организация не попала в подпись")
	}
	if !strings.Contains(body, "опубликовал") {
		t.Error("под организацией не назван человек, который опубликовал")
	}
}

// Правка заявки сбрасывает подтверждение: изменённое название никто ещё не
// смотрел, и оно не должно донашивать значок, выданный под другое имя.
func TestEditingAnApplicationDropsTheBadge(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("orgedit@example.com", "Parol123!")
	store := NewOrgStore(app.pool)
	ctx := context.Background()

	if err := store.Apply(ctx, authorID, OrgAuthor{Name: "ТОО Водоканал", Kind: "company"}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if err := store.SetStatus(ctx, authorID, orgVerified, "", nil); err != nil {
		t.Fatalf("SetStatus: %v", err)
	}
	if org, _ := store.VerifiedByUser(ctx, authorID); org == nil {
		t.Fatal("подтверждение не сработало")
	}

	if err := store.Apply(ctx, authorID, OrgAuthor{Name: "Акимат г. Костанай", Kind: "akimat"}); err != nil {
		t.Fatalf("повторная заявка: %v", err)
	}
	if org, _ := store.VerifiedByUser(ctx, authorID); org != nil {
		t.Errorf("после смены названия значок остался: %q", org.Name)
	}
}

// Неизвестный вид не должен тихо превращаться в государственный орган.
func TestUnknownKindNeverBecomesAnOfficialBody(t *testing.T) {
	for _, in := range []string{"", "акимат", "government", "  ", "official"} {
		got := NormalizeOrgKind(in)
		if got == "akimat" {
			t.Errorf("вид %q превратился в акимат", in)
		}
		if !IsOrgKind(got) {
			t.Errorf("вид %q дал неизвестное значение %q", in, got)
		}
	}
	if NormalizeOrgKind("akimat") != "akimat" {
		t.Error("настоящий акимат не распознан")
	}
	if !(OrgAuthor{Kind: "akimat"}).Official() {
		t.Error("акимат не помечен как государственный орган")
	}
	if (OrgAuthor{Kind: "company"}).Official() {
		t.Error("компания помечена как государственный орган")
	}
}

// Организация публикует для своей территории и того, что внутри неё. Без
// этого подтверждённый значок становится пропуском на всю страну.
func TestOrganisationCannotPublishOutsideItsTerritory(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	oblastID, _ := place(t, app, "Костанайская область")
	kacharID, _ := place(t, app, "Качар")
	almatyID, _ := place(t, app, "Алматы")

	store := NewOrgStore(app.pool)
	ctx := context.Background()
	akimat := &OrgAuthor{Name: "Акимат Костанайской области", Kind: "akimat", PlaceID: &oblastID}

	if ok, err := store.MayPublishFor(ctx, akimat, &kacharID); err != nil || !ok {
		t.Errorf("акимату области запрещено публиковать для посёлка внутри неё: %v, %v", ok, err)
	}
	if ok, err := store.MayPublishFor(ctx, akimat, &oblastID); err != nil || !ok {
		t.Errorf("акимату запрещено публиковать для своей же области: %v, %v", ok, err)
	}
	if ok, err := store.MayPublishFor(ctx, akimat, &almatyID); err != nil || ok {
		t.Error("акимат Костаная смог опубликовать для Алматы")
	}

	// Организация без территории не ограничена: так выглядит республиканское
	// ведомство, и заявку на него смотрел человек.
	national := &OrgAuthor{Name: "Министерство", Kind: "akimat"}
	if ok, err := store.MayPublishFor(ctx, national, &almatyID); err != nil || !ok {
		t.Errorf("организация без территории оказалась ограничена: %v, %v", ok, err)
	}
}

// Кривой БИН отсекается до модератора, чтобы не тратить его время.
func TestMalformedBINIsRefusedAtTheForm(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.createUser("orgbin@example.com", "Parol123!")
	cookie := app.login("orgbin@example.com", "Parol123!")

	w := app.do(http.MethodPost, "/studio/org", url.Values{
		"name": {"ТОО Водоканал"}, "kind": {"company"}, "bin": {"12345"},
	}, withCookie(cookie))
	if w.Code == http.StatusSeeOther {
		t.Fatal("заявка с коротким БИН была принята")
	}
	if !strings.Contains(w.Body.String(), "двенадцати") {
		t.Errorf("заявителю не объяснили, что не так с БИН")
	}

	// А правильный проходит.
	if w := app.do(http.MethodPost, "/studio/org", url.Values{
		"name": {"ТОО Водоканал"}, "kind": {"company"}, "bin": {"123456789012"},
	}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("заявка с верным БИН отклонена: %d (%s)", w.Code, w.Body.String())
	}
}
