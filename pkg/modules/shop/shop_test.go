package shop

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"shanraq.org/pkg/site"
)

// renderer builds the shop's pages in the site's template set, exactly as the
// module does at boot: the frame from pkg/site, the shop's own pages on top.
func renderer(t *testing.T) *site.Renderer {
	t.Helper()
	r := site.NewRenderer()
	r.Add(nil, templateFiles, "templates/*.html")
	if err := r.Build(); err != nil {
		t.Fatalf("parse shop templates: %v", err)
	}
	return r
}

// Every page, in every language: a missing translation key or a helper the shop
// forgot to register shows up here rather than as a 500 on a live page.
func TestShopTemplatesExecute(t *testing.T) {
	r := renderer(t)
	for _, lang := range site.Langs {
		base := site.Base{Title: "T", Lang: lang, ShowLangs: true, LangLinks: site.LangLinks("/shop", "")}
		p := Product{
			Slug: "go-book", Status: StatusAnnounced, Edition: "0.4",
			Title:   map[string]string{lang: "Книга"},
			Summary: map[string]string{lang: "Кратко"},
			Body:    map[string]string{lang: "# Заголовок\n\nтекст"},
			Plans: []Plan{
				{Code: "book", Title: map[string]string{lang: "Книга"},
					Includes: map[string]string{lang: "- PDF и EPUB"}},
				{Code: "kit", Price: 19900, Title: map[string]string{lang: "Книга и код"},
					Includes: map[string]string{lang: "- PDF, EPUB, HTML\n- код магазина"}},
			},
		}
		cases := []struct {
			name string
			data any
		}{
			{"shop_index", indexPage{Base: base, Products: []productCard{{
				Slug: "go-book", Title: "Книга", Summary: "Кратко", Edition: "0.4", Status: StatusAnnounced, From: 9900, OnSale: true,
			}}}},
			{"shop_index", indexPage{Base: base}}, // nothing published yet
			{"shop_product", productPage{Base: base, P: p, Title: "Книга", Summary: "Кратко",
				Body: site.RenderMarkdown(p.BodyIn(lang)), Notice: "ok", Packages: []packageView{
					{Code: "book", Title: "Книга", Includes: site.RenderMarkdown("- PDF и EPUB")},
					{Code: "kit", Title: "Книга и код", Price: 12900, Buyable: true,
						Includes: site.RenderMarkdown("- PDF, EPUB, HTML")},
				}}},
			{"shop_admin", adminPage{Base: base, Statuses: Statuses, Items: []adminItem{{P: p, Leads: 3, URL: "/shop/go-book"}}}},
		}
		for _, c := range cases {
			if err := r.Execute(io.Discard, c.name, c.data); err != nil {
				t.Errorf("execute %q (lang %s): %v", c.name, lang, err)
			}
		}
	}
}

// A product that is not on sale must not show a price, and one that is must not
// hide it. The two are decided in one place so a page and a checkout cannot
// disagree about whether something can be bought.
func TestOnSale(t *testing.T) {
	cases := []struct {
		name   string
		status string
		prices []int64
		want   bool
	}{
		{"selling, both packages priced", StatusSelling, []int64{9900, 19900}, true},
		{"selling, one package priced", StatusSelling, []int64{0, 19900}, true},
		{"selling, nothing priced", StatusSelling, []int64{0, 0}, false},
		{"selling, no packages at all", StatusSelling, nil, false},
		{"announced", StatusAnnounced, []int64{9900}, false},
		{"paused", StatusPaused, []int64{9900}, false},
		{"draft", StatusDraft, []int64{9900}, false},
	}
	for _, c := range cases {
		p := Product{Status: c.status}
		for i, price := range c.prices {
			p.Plans = append(p.Plans, Plan{Code: fmt.Sprintf("p%d", i), Price: price})
		}
		if got := p.OnSale(); got != c.want {
			t.Errorf("%s: OnSale = %v, want %v", c.name, got, c.want)
		}
	}
}

// The card in the list quotes the cheapest package that has a price.
func TestFrom(t *testing.T) {
	p := Product{Status: StatusSelling, Plans: []Plan{
		{Code: "kit", Price: 19900}, {Code: "book", Price: 9900}, {Code: "unpriced"},
	}}
	if got := p.From(); got != 9900 {
		t.Errorf("From = %d, want the cheapest priced package, 9900", got)
	}
	if got := (Product{}).From(); got != 0 {
		t.Errorf("with nothing priced From = %d, want 0", got)
	}
}

// The page copy falls back the way the rest of the site does: the reader's
// language, then Russian, then whatever is filled in — never an empty page.
func TestTextFallback(t *testing.T) {
	p := Product{Title: map[string]string{site.LangRU: "Книга"}}
	if got := p.TitleIn(site.LangKZ); got != "Книга" {
		t.Errorf("missing Kazakh title fell back to %q, want the Russian one", got)
	}
	p = Product{Title: map[string]string{site.LangEN: "Book"}}
	if got := p.TitleIn(site.LangKZ); got != "Book" {
		t.Errorf("with only English filled in, got %q", got)
	}
	if got := (Product{}).TitleIn(site.LangRU); got != "" {
		t.Errorf("an empty product must give an empty title, got %q", got)
	}
}

func TestValidEmail(t *testing.T) {
	for _, ok := range []string{"a@b.kz", "reader.name@mail.example.com", "x+tag@shanraq.org"} {
		if !validEmail(ok) {
			t.Errorf("validEmail(%q) = false, want true", ok)
		}
	}
	for _, bad := range []string{"", "reader", "reader@", "@shanraq.org", "a@b", "a b@c.kz", "a@b.kz\nBcc: x@y.kz", "two@at@kz.kz"} {
		if validEmail(bad) {
			t.Errorf("validEmail(%q) = true, want false", bad)
		}
	}
}

// Nobody but leadership may change what the site sells or for how much. The
// pages are registered behind the session loader, so a request with no session
// is the ordinary case here.
func TestAdminNeedsLeadership(t *testing.T) {
	m := &Module{}
	r := chi.NewRouter()
	r.Get("/admin/shop", m.handleAdmin)
	r.Post("/admin/shop/{slug}", m.handleAdminSave)
	for _, c := range []struct{ method, path string }{
		{http.MethodGet, "/admin/shop"},
		{http.MethodPost, "/admin/shop/go-book"},
	} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s %s without a session = %d, want 403", c.method, c.path, rec.Code)
		}
	}
}

// Structured data is a claim a search engine will act on, so it may never
// offer a price the page cannot take money for.
func TestProductLD(t *testing.T) {
	p := Product{
		Slug: "go-book", Kind: "book", Status: StatusAnnounced, CoverURL: "/static/shop/c.jpg",
		Title:   map[string]string{site.LangRU: "Книга"},
		Summary: map[string]string{site.LangRU: "Кратко"},
		Plans:   []Plan{{Code: "book", Price: 9900, Title: map[string]string{site.LangRU: "Книга"}}},
	}
	announced := string(productLD("n0nce", "https://shanraq.org", p, site.LangRU))
	if !strings.Contains(announced, `"@type":"Product"`) || !strings.Contains(announced, "Книга") {
		t.Errorf("product markup missing: %s", announced)
	}
	if strings.Contains(announced, "offers") || strings.Contains(announced, "9900") {
		t.Errorf("a product that is not on sale must not publish an offer: %s", announced)
	}
	if !strings.Contains(announced, "https://shanraq.org/static/shop/c.jpg") {
		t.Errorf("image must be absolute: %s", announced)
	}

	p.Status = StatusSelling
	selling := string(productLD("n0nce", "https://shanraq.org", p, site.LangRU))
	if !strings.Contains(selling, `"price":9900`) || !strings.Contains(selling, `"priceCurrency":"KZT"`) {
		t.Errorf("a product on sale must publish its price: %s", selling)
	}
	if !strings.Contains(selling, "https://schema.org/InStock") {
		t.Errorf("availability missing: %s", selling)
	}

	// The block is a script tag carrying this page's nonce, or the policy stops
	// the crawler's browser from reading it.
	if !strings.Contains(selling, `nonce="n0nce"`) || !strings.Contains(selling, "application/ld+json") {
		t.Errorf("the block must be a nonced ld+json script: %s", selling)
	}
}
