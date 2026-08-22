package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// place возвращает узел справочника и его адрес.
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

// Каждому месту нужен адрес, иначе страницы у него быть не может. Справочник
// приходит с колонкой code, но она неоднородна: у Качара kz-kostanay-kachar, у
// Костанайской области g65, и /place/g65 — шифр, а не адрес.
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

	// Тёзки разводятся, а не затирают друг друга: Алматы — и город, и районы.
	var dupes int
	if err := app.pool.QueryRow(context.Background(),
		`SELECT count(*) FROM (SELECT slug FROM geo_nodes GROUP BY slug HAVING count(*) > 1) d`).Scan(&dupes); err != nil {
		t.Fatalf("проверка тёзок: %v", err)
	}
	if dupes != 0 {
		t.Errorf("одинаковых адресов: %d", dupes)
	}
}

// Лента места смотрит вверх, а не вниз. Место, которое выбрал автор, — это
// адресат: объявление для Качара написано качарцам, и на странице области оно
// попало бы к сотне тысяч человек, которым не предназначалось. Наоборот —
// можно: областное сообщение адресовано и жителям Качара.
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

	// Страница посёлка несёт своё и адресованное всей области.
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

	// Страница области несёт только адресованное области, но не то, что
	// написано для одного посёлка внутри неё.
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

// Пустое место — обычное состояние, а не ошибка: страница существует, чтобы
// было куда опубликовать первым.
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

// Автор указывает место при публикации, и оно переживает правку статьи.
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

	// И возвращается в форму, чтобы правка не стирала выбор.
	page := app.do(http.MethodGet, "/studio/a/"+id.String(), nil, withCookie(cookie))
	if !strings.Contains(page.Body.String(), kacharID.String()) {
		t.Error("редактор не вернул выбранное место в форму")
	}
}

// В sitemap попадают места, для которых что-то написано, — и только они.
// Областное сообщение видно и на странице каждого посёлка внутри области, но
// предлагать поисковику тридцать страниц с одной и той же статьёй значит
// тридцать раз показать пустоту, а не расширить охват.
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

// Путь наверх должен вести наверх. Ancestry писалась для объявлений, где нужны
// только названия, и вернувшись на страницу места без адреса рисовала
// «Костанайская область» ссылкой на /place/ — то есть в никуда.
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

// Карточка в ленте должна называть организацию. Страница места существует ради
// сообщений ЖКХ и акимата; карточка, подписанная «А. Смағұлова», прячет ровно
// то, ради чего заведены организации-авторы.
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

	// И один запрос на всю ленту, а не по одному на карточку.
	names, err := store.VerifiedNames(ctx, []uuid.UUID{authorID})
	if err != nil || names[authorID] != "КСК «Качарец»" {
		t.Errorf("пакетный поиск не нашёл организацию: %v, %v", names, err)
	}
}

// Разделитель стоит между звеньями, а не после последнего.
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

// Матрица видимости в общей ленте. Место, которое поставил автор, — это круг
// читателей: качарское объявление об отключении света жителю Алматы не новость,
// а мусор в ленте, и гостю, о котором мы ничего не знаем, — тоже.
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

	// Гость и всякий, кто не назвал места, видит только общее.
	check("гость", sees(nil), map[string]bool{
		forEveryone: true, forKachar: false, forOblast: false, forAlmaty: false,
	})

	// Житель Качара: своё, областное, общее — но не чужой город.
	kacharReader, err := geo.AddressedTo(ctx, kacharID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель Качара", sees(kacharReader), map[string]bool{
		forEveryone: true, forKachar: true, forOblast: true, forAlmaty: false,
	})

	// Житель областного центра — не житель Качара: качарское его не касается.
	oblastReader, err := geo.AddressedTo(ctx, oblastID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель области вне Качара", sees(oblastReader), map[string]bool{
		forEveryone: true, forOblast: true, forKachar: false, forAlmaty: false,
	})

	// Житель Алматы: своё и общее.
	almatyReader, err := geo.AddressedTo(ctx, almatyID)
	if err != nil {
		t.Fatalf("AddressedTo: %v", err)
	}
	check("житель Алматы", sees(almatyReader), map[string]bool{
		forEveryone: true, forAlmaty: true, forKachar: false, forOblast: false,
	})
}

// То же правило — на самой странице: гость на главной не должен видеть местное.
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

	// А на своей странице места оно есть — там его и ищут.
	place := app.do(http.MethodGet, "/place/kachar", nil).Body.String()
	if !strings.Contains(place, localSlug) {
		t.Error("местное объявление пропало и со страницы места")
	}
}
