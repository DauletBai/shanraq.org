package shop

import (
	"io"
	"net/http"
	"net/http/httptest"
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
			Slug: "go-book", Status: StatusAnnounced, Price: 0, Edition: "0.2",
			Title:   map[string]string{lang: "Книга"},
			Summary: map[string]string{lang: "Кратко"},
			Body:    map[string]string{lang: "# Заголовок\n\nтекст"},
		}
		cases := []struct {
			name string
			data any
		}{
			{"shop_index", indexPage{Base: base, Products: []productCard{{
				Slug: "go-book", Title: "Книга", Summary: "Кратко", Edition: "0.2", Status: StatusAnnounced,
			}}}},
			{"shop_index", indexPage{Base: base}}, // nothing published yet
			{"shop_product", productPage{Base: base, P: p, Title: "Книга", Summary: "Кратко",
				Body: site.RenderMarkdown(p.BodyIn(lang)), Notice: "ok"}},
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
		status string
		price  int64
		want   bool
	}{
		{StatusSelling, 12000, true},
		{StatusSelling, 0, false}, // priced at nothing is not for sale
		{StatusAnnounced, 12000, false},
		{StatusPaused, 12000, false},
		{StatusDraft, 12000, false},
	}
	for _, c := range cases {
		p := Product{Status: c.status, Price: c.price}
		if got := p.OnSale(); got != c.want {
			t.Errorf("status %q price %d: OnSale = %v, want %v", c.status, c.price, got, c.want)
		}
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
