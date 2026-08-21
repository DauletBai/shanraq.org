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

// Лента места смотрит вниз, а не вверх: страница области несёт написанное для
// посёлков внутри неё. Обратное было бы неверно — страница Качара, забитая
// областными сообщениями, похоронила бы то немногое, что про сам Качар.
func TestPlaceFeedLooksDownwardsNotUpwards(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()

	authorID := app.createUser("placefeed@example.com", "Parol123!")
	oblastID, oblastSlug := place(t, app, "Костанайская область")
	kacharID, kacharSlug := place(t, app, "Качар")

	forOblast, oblastArticle := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, forOblast, oblastID)
	forKachar, kacharArticle := app.seedArticle(authorID, "published")
	app.exec(`UPDATE articles SET geo_node_id = $2, published_at = NOW() WHERE id = $1`, forKachar, kacharID)

	// Страница области несёт обе.
	w := app.do(http.MethodGet, "/place/"+oblastSlug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница области: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, kacharArticle) {
		t.Error("страница области не показывает материал посёлка внутри неё")
	}
	if !strings.Contains(body, oblastArticle) {
		t.Error("страница области не показывает свой же материал")
	}

	// Страница посёлка несёт только своё.
	w = app.do(http.MethodGet, "/place/"+kacharSlug, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("страница посёлка: %d", w.Code)
	}
	body = w.Body.String()
	if !strings.Contains(body, kacharArticle) {
		t.Error("страница посёлка не показывает свой материал")
	}
	if strings.Contains(body, oblastArticle) {
		t.Error("страница посёлка засыпана областными материалами")
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

// В sitemap попадают места, о которых есть что читать, — вместе с их
// родителями, чтобы страница области предлагалась поисковику даже когда
// написано пока только про посёлок.
func TestSitemapCarriesPlacesThatHaveArticles(t *testing.T) {
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
	if !has(oblastSlug) {
		t.Error("области над этим местом нет в списке для sitemap")
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
