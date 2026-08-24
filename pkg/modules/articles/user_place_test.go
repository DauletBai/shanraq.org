package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// kacharID returns the reference node for the settlement of Kachar — the very
// example this was all built for: a housing service chair publishing for their own
// settlement.
func kacharID(t *testing.T, app *testApp) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id FROM geo_nodes WHERE name_ru = 'Качар' LIMIT 1`).Scan(&id); err != nil {
		t.Skipf("в тестовой базе нет справочника мест: %v", err)
	}
	return id
}

// The place is one field, not four. The republic and the region are derived from it
// by the tree, and they cannot drift apart from the city as four separate columns
// would.
func TestUserPlaceIsOneNodeAndItsAncestryGivesTheRest(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	userID := app.createUser("place1@example.com", "Parol123!")
	node := kacharID(t, app)
	geo := NewGeoStore(app.pool)

	if err := geo.SetUserPlace(context.Background(), userID, &node); err != nil {
		t.Fatalf("SetUserPlace: %v", err)
	}
	got, err := geo.UserPlace(context.Background(), userID)
	if err != nil || got == nil || *got != node {
		t.Fatalf("UserPlace вернул %v, %v", got, err)
	}

	chain, err := geo.Ancestry(context.Background(), node, LangRU)
	if err != nil {
		t.Fatalf("Ancestry: %v", err)
	}
	if len(chain) < 3 {
		t.Fatalf("цепочка предков короткая: %+v", chain)
	}
	if chain[0].Kind != "country" {
		t.Errorf("цепочка начинается не со страны: %+v", chain[0])
	}
	if chain[len(chain)-1].Name != "Качар" {
		t.Errorf("цепочка кончается не выбранным местом: %+v", chain[len(chain)-1])
	}

	// The reader's own label: the place and its region, so a namesake cannot be confused with it.
	label, err := geo.PlaceLabel(context.Background(), node, LangRU)
	if err != nil {
		t.Fatalf("PlaceLabel: %v", err)
	}
	if !strings.Contains(label, "Качар") || !strings.Contains(label, "область") {
		t.Errorf("подпись места не называет область: %q", label)
	}
}

// The place was chosen in haste, or the person has moved — they must be able to fix
// it themselves, without writing to support. And to erase it altogether.
func TestReaderCanChangeAndClearTheirPlace(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.createUser("place2@example.com", "Parol123!")
	cookie := app.login("place2@example.com", "Parol123!")
	node := kacharID(t, app)

	w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {node.String()}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("сохранение места: %d (%s)", w.Code, w.Body.String())
	}

	page := app.do(http.MethodGet, "/studio/profile", nil, withCookie(cookie))
	if !strings.Contains(page.Body.String(), "Качар") {
		t.Error("профиль не показывает выбранное место")
	}

	// An empty value erases the choice rather than breaking the form.
	if w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {""}}, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("очистка места: %d", w.Code)
	}
	var left int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM user_places WHERE user_id = (SELECT id FROM auth_users WHERE email = $1)`,
		"place2@example.com").Scan(&left); err != nil {
		t.Fatalf("проверка: %v", err)
	}
	if left != 0 {
		t.Errorf("после очистки осталось %d записей", left)
	}
}

// Rubbish in the field must neither crash nor be quietly saved.
func TestBadPlaceIsRefusedWithoutBreakingTheProfile(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.createUser("place3@example.com", "Parol123!")
	cookie := app.login("place3@example.com", "Parol123!")

	w := app.do(http.MethodPost, "/studio/place",
		url.Values{"geo_node_id": {"не-идентификатор"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("получили %d, ожидался переход обратно в профиль", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "place_bad") {
		t.Errorf("читателю не сказали, что не вышло: %s", loc)
	}
}

// The field is optional, and that is a point of principle: registration is the
// bottleneck, and demanding a settlement before an account exists would cost more
// than it gives. Whoever chose none sees everything, exactly as before places existed.
func TestRegistrationWorksWithoutAPlace(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	app.emails = append(app.emails, "placefree@example.com")
	w := app.do(http.MethodPost, "/studio/register", url.Values{
		"first_name": {"Тест"}, "last_name": {"Пользователь"},
		"email": {"placefree@example.com"}, "password": {"Parol123!"},
		"consent": {"on"},
	})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("регистрация без места: %d (%s)", w.Code, w.Body.String())
	}
	var n int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM auth_users WHERE email = $1`, "placefree@example.com").Scan(&n); err != nil {
		t.Fatalf("проверка: %v", err)
	}
	if n != 1 {
		t.Errorf("аккаунт не создан: %d", n)
	}
}

// Every field on the form has to explain why it is there. The place all the more so:
// a person is telling us where they live and is entitled to know what becomes of it.
func TestPlaceFieldExplainsItself(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	w := app.do(http.MethodGet, "/studio/register", nil)
	body := w.Body.String()
	if !strings.Contains(body, `data-geo`) {
		t.Fatal("в форме регистрации нет выбора места")
	}
	for _, must := range []string{"Где вы живёте", "Необязательно"} {
		if !strings.Contains(body, must) {
			t.Errorf("нет пояснения к полю места: %q", must)
		}
	}
	if !strings.Contains(body, "fhelp") {
		t.Error("у поля места нет подсказки «?»")
	}
}
