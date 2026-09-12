package shop

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/site"
)

// indexPage is the shop's front page: everything the house sells.
type indexPage struct {
	site.Base
	Products []productCard
}

// productCard is one product as a card in the list.
type productCard struct {
	Slug     string
	Title    string
	Summary  string
	CoverURL string
	Edition  string
	Status   string
	// From is the cheapest package's price, 0 while nothing is priced.
	From   int64
	OnSale bool
	Draft  bool
}

// productPage is one product's own page.
type productPage struct {
	site.Base
	P       Product
	Title   string
	Summary string
	Body    template.HTML
	// Packages are what the reader chooses between: the book, or the book and
	// the code. Rendered in the product's own order, cheapest first.
	Packages []packageView
	// Notice is the answer to a just-submitted notify form: "you are on the
	// list", "you already were", or "that address does not look right".
	Notice string
	Bad    bool
}

// packageView is one package as the page shows it.
type packageView struct {
	Code     string
	Title    string
	Includes template.HTML
	Price    int64
	Buyable  bool
}

func (m *Module) handleIndex(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	staff := m.isStaff(r)
	items, err := m.store.List(r.Context(), staff)
	if err != nil {
		m.logger().Error("shop list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := indexPage{Base: m.rt.Site.Base(r, site.T(lang, "shop.title"), lang)}
	page.Desc = site.T(lang, "shop.seo_desc")
	for _, p := range items {
		page.Products = append(page.Products, productCard{
			Slug: p.Slug, Title: p.TitleIn(lang), Summary: p.SummaryIn(lang),
			CoverURL: p.CoverURL, Edition: p.Edition, Status: p.Status,
			From: p.From(), OnSale: p.OnSale(), Draft: p.Status == StatusDraft,
		})
	}
	m.render(w, "shop_index", page)
}

func (m *Module) handleProduct(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	p, err := m.store.BySlug(r.Context(), chi.URLParam(r, "slug"), m.isStaff(r))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	title := p.TitleIn(lang)
	page := productPage{
		Base:    m.rt.Site.Base(r, title, lang),
		P:       p,
		Title:   title,
		Summary: p.SummaryIn(lang),
		Body:    site.RenderMarkdown(p.BodyIn(lang)),
	}
	for _, pl := range p.Plans {
		page.Packages = append(page.Packages, packageView{
			Code:  pl.Code,
			Title: pl.TitleIn(lang),
			// The list is Markdown so the owner writes it in the panel the way
			// the rest of the site is written, without markup in a form field.
			Includes: site.RenderMarkdown(pl.IncludesIn(lang)),
			Price:    pl.Price,
			Buyable:  pl.Buyable(),
		})
	}
	if s := strings.TrimSpace(page.Summary); s != "" {
		page.Desc = s
	}
	if p.CoverURL != "" {
		page.OGImage = absolute(page.SiteURL, p.CoverURL)
	}
	page.JSONLD = productLD(page.Nonce, page.SiteURL, p, lang)
	switch r.URL.Query().Get("notified") {
	case "1":
		page.Notice = site.T(lang, "shop.notify_ok")
	case "again":
		page.Notice = site.T(lang, "shop.notify_again")
	case "bad":
		page.Notice, page.Bad = site.T(lang, "shop.notify_bad"), true
	}
	m.render(w, "shop_product", page)
}

// handleNotify records an address that wants to hear when a product is ready.
// It answers with a redirect so a refresh does not re-submit, and it never says
// whether the address was already on the list to somebody who is not that
// address — "you are on the list" is the answer either way.
func (m *Module) handleNotify(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	slug := chi.URLParam(r, "slug")
	p, err := m.store.BySlug(r.Context(), slug, false)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	back := "/shop/" + slug + "?lang=" + lang
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, back+"&notified=bad", http.StatusSeeOther)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	if !validEmail(email) {
		http.Redirect(w, r, back+"&notified=bad", http.StatusSeeOther)
		return
	}
	added, err := m.store.AddLead(r.Context(), p.ID, email, lang)
	if err != nil {
		m.logger().Error("shop lead", zap.Error(err))
		http.Redirect(w, r, back+"&notified=bad", http.StatusSeeOther)
		return
	}
	if added {
		http.Redirect(w, r, back+"&notified=1", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, back+"&notified=again", http.StatusSeeOther)
}

// validEmail is the same shallow check the rest of the site uses: one @, no
// spaces, something on each side. The address is proved by the letter that is
// sent to it, not by a regular expression.
func validEmail(s string) bool {
	if len(s) < 5 || len(s) > 254 || strings.ContainsAny(s, " \t\r\n") {
		return false
	}
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-3 && strings.LastIndexByte(s, '@') == at &&
		strings.Contains(s[at:], ".")
}

// productLD describes the product to a search engine: what it is, what it looks
// like, who publishes it — and what it costs, but only while it can actually be
// bought. A price in structured data is a promise that the page can take money
// for it; until an acquirer is connected, the honest markup is a product with
// no offers rather than an offer nobody can accept.
func productLD(nonce, origin string, p Product, lang string) template.HTML {
	ld := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Product",
		"name":        p.TitleIn(lang),
		"description": p.SummaryIn(lang),
		"url":         absolute(origin, "/shop/"+p.Slug+"?lang="+lang),
		"brand":       map[string]any{"@type": "Organization", "name": "Shanraq.org"},
		"category":    p.Kind,
	}
	if p.CoverURL != "" {
		ld["image"] = absolute(origin, p.CoverURL)
	}
	if !p.OnSale() {
		return site.JSONLD(nonce, ld)
	}
	offers := make([]map[string]any, 0, len(p.Plans))
	for _, pl := range p.Plans {
		if !pl.Buyable() {
			continue
		}
		offers = append(offers, map[string]any{
			"@type":         "Offer",
			"name":          pl.TitleIn(lang),
			"price":         pl.Price,
			"priceCurrency": currencyOr(p.Currency),
			"availability":  "https://schema.org/InStock",
			"url":           absolute(origin, "/shop/"+p.Slug+"?lang="+lang),
		})
	}
	ld["offers"] = offers
	return site.JSONLD(nonce, ld)
}

// currencyOr keeps the markup valid for a row written before the column had a
// default: an empty currency would publish a price in no currency at all.
func currencyOr(c string) string {
	if c = strings.TrimSpace(c); c != "" {
		return c
	}
	return "KZT"
}

// absolute turns a site-relative URL into an absolute one for a social preview.
func absolute(origin, u string) string {
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return strings.TrimSuffix(origin, "/") + u
}

// isStaff reports whether the viewer may see drafts and the admin pages.
func (m *Module) isStaff(r *http.Request) bool {
	c, ok := auth.ClaimsFromContext(r.Context())
	return ok && c != nil && c.HasAnyRole("admin", "director", "manager", "editor")
}

// canEdit is the narrower right: changing what the site sells, and for how
// much, is leadership's call.
func canEdit(c *auth.Claims) bool { return c != nil && c.HasAnyRole("admin", "director") }

func (m *Module) render(w http.ResponseWriter, name string, data any) {
	if err := m.rt.Site.Render(w, name, data); err != nil {
		m.logger().Error("render shop template", zap.String("template", name), zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
