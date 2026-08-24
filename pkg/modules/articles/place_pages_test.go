package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// place returns a reference node and its address.
func place(t *testing.T, app *testApp, name string) (uuid.UUID, string) {
	t.Helper()
	var id uuid.UUID
	var slug string
	if err := app.pool.QueryRow(context.Background(),
		`SELECT id, COALESCE(slug,'') FROM geo_nodes WHERE name_ru = $1 ORDER BY level LIMIT 1`, name).
		Scan(&id, &slug); err != nil {
		t.Skipf("в тестовой базе нет места %q: %v", name, err)
	}
	if slug == "" {
		t.Skipf("у места %q нет адреса", name)
	}
	return id, slug
}

// Every place needs an address, or it can have no page. The reference arrives
// with a code column, but that column is not uniform: Kachar has
// kz-kostanay-kachar, the Kostanay region has g65, and /place/g65 is a cipher,
// not an address.
func TestEveryPlaceHasAReadableAddress(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	var empty int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM geo_nodes WHERE slug IS NULL OR slug = ''`).Scan(&empty); err != nil {
		t.Fatalf("проверка: %v", err)
	}
	if empty != 0 {
		t.Errorf("без адреса осталось мест: %d", empty)
	}

	_, slug := place(t, app, "Качар")
	if slug != "kachar" {
		t.Errorf("адрес Качара %q, ожидался kachar", slug)
	}

	// Namesakes are separated rather than overwriting each other: Almaty is both a city and its districts.
	var dupes int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM (SELECT slug FROM geo_nodes GROUP BY slug HAVING count(*) > 1) d`).Scan(&dupes); err != nil {
		t.Fatalf("проверка тёзок: %v", err)
	}
	if dupes != 0 {
		t.Errorf("одинаковых адресов: %d", dupes)
	}
}

// A place's feed looks upward, not down. The place an author chose is the
// addressee: a notice for Kachar is written for the people of Kachar, and on the
// region's page it would reach a hundred thousand it was not meant for. The other
// way round is fine: a regional message is addressed to Kachar's residents too.
func TestPlaceFeedCarriesWhatIsAddressedToThePlace(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("placefeed@example.com", "Parol123!")
	oblastID, oblastSlug := place(t, app, "Костанайская область")
	kacharID, kacharSlug := place(t, app, "Качар")

	forOblast, oblastArticle := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, forOblast, oblastID)
	forKachar, kacharArticle := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, forKachar, kacharID)

	// A settlement's page carries its own and what is addressed to the whole region.
	w := app.do(http.MethodGet, "/place/"+kacharSlug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница посёлка: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, kacharArticle) {
		t.Error("страница посёлка не показывает свой материал")
	}
	if !strings.Contains(body, oblastArticle) {
		t.Error("страница посёлка не показывает областное сообщение, адресованное и ей")
	}

	// A region's page carries only what is addressed to the region, and not what was
	// written for one settlement inside it.
	w = app.do(http.MethodGet, "/place/"+oblastSlug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница области: %d", w.Code)
	}
	body = w.Body.String()
	if !strings.Contains(body, oblastArticle) {
		t.Error("страница области не показывает свой же материал")
	}
	if strings.Contains(body, kacharArticle) {
		t.Error("написанное для одного посёлка попало на страницу всей области")
	}
}

// An empty place is the ordinary state, not an error: the page exists so that
// there is somewhere for the first person to publish.
func TestEmptyPlaceStillHasAPage(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	_, slug := place(t, app, "Качар")
	w := app.do(http.MethodGet, "/place/"+slug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("получили %d, страница места должна открываться всегда", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Качар") {
		t.Error("страница не называет место")
	}
	if w := app.do(http.MethodGet, "/place/takogo-mesta-net", nil); w.Code != http.StatusNotFound {
		t.Errorf("несуществующее место вернуло %d, ожидался 404", w.Code)
	}
}

// The author sets the place on publication, and it survives an edit of the article.
func TestAuthorPicksThePlaceAndItSticks(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("placeauthor@example.com", "Parol123!")
	app.exec(`UPDATE auth_users SET email_verified_at = NOW() WHERE id = $1`, authorID)
	cookie := app.login("placeauthor@example.com", "Parol123!")
	id, _ := app.seedArticle(authorID, "draft")
	kacharID, _ := place(t, app, "Качар")

	form := url.Values{
		"original_lang": {"ru"}, "category": {"society"}, "subcategory": {""},
		"title_ru": {"Заголовок"}, "summary_ru": {"Описание"}, "body_ru": {"Текст"},
		"geo_node_id": {kacharID.String()},
	}
	if w := app.do(http.MethodPost, "/studio/a/"+id.String(), form, withCookie(cookie)); w.Code != http.StatusSeeOther {
		t.Fatalf("сохранение статьи: %d (%s)", w.Code, w.Body.String())
	}

	got, err := NewStore(app.pool).ArticlePlace(context.Background(), id)
	if err != nil || got == nil || *got != kacharID {
		t.Fatalf("место статьи не сохранилось: %v, %v", got, err)
	}

	// And it comes back into the form, so an edit does not erase the choice.
	page := app.do(http.MethodGet, "/studio/a/"+id.String(), nil, withCookie(cookie))
	if !strings.Contains(page.Body.String(), kacharID.String()) {
		t.Error("редактор не вернул выбранное место в форму")
	}
}

// The sitemap holds the places something has been written for, and only those. A
// regional message shows on every settlement page inside the region, but offering
// a search engine thirty pages carrying the same article means showing an
// emptiness thirty times, not widening the reach.
func TestSitemapCarriesOnlyPlacesSomethingWasWrittenFor(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("placemap@example.com", "Parol123!")
	kacharID, kacharSlug := place(t, app, "Качар")
	_, oblastSlug := place(t, app, "Костанайская область")

	id, _ := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, id, kacharID)

	places, err := NewStore(app.pool).PlacesWithArticles(context.Background())
	if err != nil {
		t.Fatalf("PlacesWithArticles: %v", err)
	}
	has := func(s string) bool {
		for _, p := range places {
			if p == s {
				return true
			}
		}
		return false
	}
	if !has(kacharSlug) {
		t.Error("места со статьёй нет в списке для sitemap")
	}
	if has(oblastSlug) {
		t.Error("область без собственных материалов попала в sitemap")
	}
}

// The path upward has to lead upward. Ancestry was written for listings, which
// need only the names, and coming back to a place page without addresses it drew
// "Kostanay region" as a link to /place/ — that is, to nowhere.
func TestTheWayUpActuallyLeadsSomewhere(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	kacharID, _ := place(t, app, "Качар")
	_, oblastSlug := place(t, app, "Костанайская область")

	chain, err := NewGeoStore(app.pool).Ancestry(context.Background(), kacharID, "ru")
	if err != nil {
		t.Fatalf("Ancestry: %v", err)
	}
	for _, n := range chain {
		if n.Slug == "" {
			t.Errorf("у предка %q нет адреса", n.Name)
		}
	}

	body := app.do(http.MethodGet, "/place/kachar", nil).Body.String()
	if strings.Contains(body, `href="/place/?`) {
		t.Error("хлебная крошка ведёт в никуда")
	}
	if !strings.Contains(body, "/place/"+oblastSlug) {
		t.Error("со страницы посёлка нельзя подняться в область")
	}
}

// A card in the feed has to name the organisation. A place page exists for
// notices from the utilities and the mayor's office; a card signed "A. Smagulova"
// hides the very thing organisation authors were created for.
func TestAFeedCardNamesTheOrganisation(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("orgfeed@example.com", "Parol123!")
	kacharID, kacharSlug := place(t, app, "Качар")
	ctx := context.Background()

	store := NewOrgStore(app.pool)
	if err := store.Apply(ctx, authorID, OrgAuthor{Name: "КСК «Качарец»", Kind: "utility"}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if err := store.SetStatus(ctx, authorID, orgVerified, "", nil); err != nil {
		t.Fatalf("SetStatus: %v", err)
	}

	id, _ := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, id, kacharID)

	body := app.do(http.MethodGet, "/place/"+kacharSlug, nil).Body.String()
	if !strings.Contains(body, "КСК «Качарец»") {
		t.Error("карточка в ленте места не называет организацию")
	}

	// And one query for the whole feed, not one per card.
	names, err := store.VerifiedNames(ctx, []uuid.UUID{authorID})
	if err != nil || names[authorID] != "КСК «Качарец»" {
		t.Errorf("пакетный поиск не нашёл организацию: %v, %v", names, err)
	}
}

// The separator goes between the links, not after the last one.
func TestBreadcrumbHasNoTrailingSeparator(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	place(t, app, "Качар")
	body := app.do(http.MethodGet, "/place/kachar", nil).Body.String()
	i := strings.Index(body, "place-up")
	if i < 0 {
		t.Fatal("на странице нет пути наверх")
	}
	nav := body[i:]
	nav = nav[:strings.Index(nav, "</nav>")]
	if strings.Contains(nav[strings.LastIndex(nav, "</a>"):], "·") {
		t.Errorf("после последнего звена висит разделитель: %q", nav[strings.LastIndex(nav, "</a>"):])
	}
}

// The visibility matrix for the common feed. The place an author set is a circle
// of readers: a Kachar notice about a power cut is not news to someone in Almaty
// but litter in their feed — and the same for a guest we know nothing about.
func TestTheFeedCarriesOnlyWhatIsAddressedToTheReader(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	store := NewStore(app.pool)
	geo := NewGeoStore(app.pool)

	kacharID, _ := place(t, app, "Качар")
	oblastID, _ := place(t, app, "Костанайская область")
	almatyID, _ := place(t, app, "Алматы")

	author := app.createUser("feedscope@example.com", "Parol123!")
	mark := func(status string, node *uuid.UUID) string {
		id, slug := app.seedArticle(author, status)
		app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, id, node)
		return slug
	}
	forEveryone := mark("published", nil)
	forKachar := mark("published", &kacharID)
	forOblast := mark("published", &oblastID)
	forAlmaty := mark("published", &almatyID)

	sees := func(addressed []uuid.UUID) map[string]bool {
		arts, err := store.ListPublished(ctx, "", "", "", 60, 0, addressed)
		if err != nil {
			t.Fatalf("ListPublished: %v", err)
		}
		out := map[string]bool{}
		for _, a := range arts {
			out[a.Slug] = true
		}
		return out
	}
	check := func(who string, got map[string]bool, want map[string]bool) {
		t.Helper()
		for slug, expected := range want {
			if got[slug] != expected {
				verb := "не увидел"
				if !expected {
					verb = "увидел лишнее:"
				}
				t.Errorf("%s %s %s", who, verb, slug)
			}
		}
	}

	// A guest, and anyone who named no place, sees only the general.
	check("гость", sees(nil), map[string]bool{
		forEveryone: true, forKachar: false, forOblast: false, forAlmaty: false,
	})

	// A Kachar resident: their own, the region's, the general — but not another town's.
	kacharReader, err := geo.AddressedTo(ctx, kacharID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель Качара", sees(kacharReader), map[string]bool{
		forEveryone: true, forKachar: true, forOblast: true, forAlmaty: false,
	})

	// A resident of the regional centre is not a Kachar resident: Kachar's does not concern them.
	oblastReader, err := geo.AddressedTo(ctx, oblastID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель области вне Качара", sees(oblastReader), map[string]bool{
		forEveryone: true, forOblast: true, forKachar: false, forAlmaty: false,
	})

	// An Almaty resident: their own and the general.
	almatyReader, err := geo.AddressedTo(ctx, almatyID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель Алматы", sees(almatyReader), map[string]bool{
		forEveryone: true, forAlmaty: true, forKachar: false, forOblast: false,
	})
}

// The same rule on the page itself: a guest on the front page must not see the local.
func TestAGuestOnTheHomePageSeesNoLocalNotices(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	kacharID, _ := place(t, app, "Качар")
	author := app.createUser("feedguest@example.com", "Parol123!")

	localID, localSlug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, localID, kacharID)
	commonID, commonSlug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() WHERE id = $1`, commonID)

	body := app.do(http.MethodGet, "/", nil).Body.String()
	if !strings.Contains(body, commonSlug) {
		t.Error("гость не увидел материал, написанный для всех")
	}
	if strings.Contains(body, localSlug) {
		t.Error("местное объявление попало в ленту гостя")
	}

	// And on its own place page it is there — which is where people look for it.
	place := app.do(http.MethodGet, "/place/kachar", nil).Body.String()
	if !strings.Contains(place, localSlug) {
		t.Error("местное объявление пропало и со страницы места")
	}
}
