package articles

import (
	"net/http"
	"strings"
	"testing"
)

// The breadcrumb is what turns a bare URL in search results into a labelled
// trail. Google requires the last crumb to be the page itself and unlinked.
func TestBreadcrumbLD(t *testing.T) {
	page := &ArticlePage{
		// ArticlePage declares its own Title, which shadows Base.Title — the
		// breadcrumb and the headline both read the outer one.
		Base:        Base{Lang: "ru", SiteURL: "https://shanraq.org"},
		Slug:        "test",
		Title:       "Заголовок",
		Category:    "world",
		Subcategory: "europe",
	}
	ld := breadcrumbLD(page)
	if ld["@type"] != "BreadcrumbList" {
		t.Fatalf("@type = %v", ld["@type"])
	}
	items, _ := ld["itemListElement"].([]map[string]any)
	if len(items) != 4 {
		t.Fatalf("crumbs = %d, want site → category → subcategory → article", len(items))
	}
	for i, it := range items {
		if it["position"] != i+1 {
			t.Errorf("crumb %d has position %v", i, it["position"])
		}
	}
	if items[0]["name"] != "Shanraq.org" {
		t.Errorf("first crumb = %v", items[0]["name"])
	}
	if _, linked := items[3]["item"]; linked {
		t.Error("the last crumb must not carry a link")
	}
	if items[3]["name"] != "Заголовок" {
		t.Errorf("last crumb = %v", items[3]["name"])
	}

	// An article filed without a subcategory gets a three-step trail, not one
	// with a hole where the missing level was.
	page.Subcategory = ""
	if got := len(breadcrumbLD(page)["itemListElement"].([]map[string]any)); got != 3 {
		t.Errorf("crumbs without a subcategory = %d, want 3", got)
	}
	page.Category = ""
	if got := len(breadcrumbLD(page)["itemListElement"].([]map[string]any)); got != 2 {
		t.Errorf("crumbs without a category = %d, want 2", got)
	}
}

// Новостная карта — это то, чем издание сообщает Google, что вышло за
// последние двое суток. Мы издаём одну статью на трёх языках по трём адресам,
// а карта перечисляла только язык оригинала: две трети написанного не попадали
// в окно вообще.
func TestNewsSitemapCarriesEveryFinishedLanguage(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	author := app.createUser("newsmap@example.com", "Parol123!")
	id, slug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() WHERE id = $1`, id)
	app.exec(`INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status)
	          VALUES ($1,'kz','Қазақша тақырып','Сипаттама','Мәтін','human','ready'),
	                 ($1,'en','English headline','Summary','Body','human','ready')`, id)
	// Пустая вкладка — это намерение, а не материал: анонсировать её нельзя.
	other, otherSlug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() WHERE id = $1`, other)
	app.exec(`INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status)
	          VALUES ($1,'en','','','','human','draft')`, other)

	body := app.do(http.MethodGet, "/sitemap-news.xml", nil).Body.String()
	for _, want := range []struct{ url, lang, title string }{
		{"/read/" + slug + "?lang=ru", "ru", "Тест заголовок"},
		{"/read/" + slug + "?lang=kz", "kk", "Қазақша тақырып"},
		{"/read/" + slug + "?lang=en", "en", "English headline"},
	} {
		if !strings.Contains(body, want.url) {
			t.Errorf("в карте нет адреса %s", want.url)
		}
		if !strings.Contains(body, "<news:language>"+want.lang+"</news:language>") {
			t.Errorf("в карте нет языка %s", want.lang)
		}
		if !strings.Contains(body, "<news:title>"+want.title+"</news:title>") {
			t.Errorf("в карте нет заголовка %q", want.title)
		}
	}
	if strings.Contains(body, "/read/"+otherSlug+"?lang=en") {
		t.Error("пустая языковая вкладка объявлена как новость")
	}

	// Казахский код языка для Google — kk, а не наш внутренний kz.
	if strings.Contains(body, "<news:language>kz</news:language>") {
		t.Error("в карту попал внутренний код kz вместо kk")
	}

	// И старое остаётся за окном: Google принимает только двое суток.
	old, oldSlug := app.seedArticle(author, "published")
	app.exec(`UPDATE articles SET published_at = NOW() - interval '5 days' WHERE id = $1`, old)
	if body := app.do(http.MethodGet, "/sitemap-news.xml", nil).Body.String(); strings.Contains(body, oldSlug) {
		t.Error("статья старше двух суток попала в новостную карту")
	}
}
