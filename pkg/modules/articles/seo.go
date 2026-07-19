package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// htmlLang maps our internal UI code to the BCP-47 language subtag used in
// HTML lang / hreflang / schema.org. Kazakh's ISO 639-1 code is "kk"; we keep
// "kz" internally (routing, ?lang=) but must present "kk" to browsers/crawlers.
func htmlLang(lang string) string {
	if lang == LangKZ {
		return "kk"
	}
	return lang
}

// ogLocale maps a UI language to an Open Graph locale.
func ogLocale(lang string) string {
	switch lang {
	case LangKZ:
		return "kk_KZ"
	case LangEN:
		return "en_US"
	default:
		return "ru_RU"
	}
}

// jsonLD renders a schema.org value as an inline ld+json script. encoding/json
// escapes <, > and & so the payload is safe inside <script>.
func jsonLD(v any) template.HTML {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return template.HTML(`<script type="application/ld+json">` + string(b) + `</script>`)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (m *Module) siteURL() string {
	return strings.TrimRight(m.rt.Config.PublicBase(), "/")
}

// handleRobots serves robots.txt: crawl everything except the author cabinet,
// and point crawlers at the sitemap.
func (m *Module) handleRobots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /studio\nDisallow: /studio/\n\nSitemap: %s/sitemap.xml\n", m.siteURL())
}

func seoURL(site, path, lang string) string {
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return xmlEscape(site + path + sep + "lang=" + lang)
}

// handleSitemap emits a trilingual sitemap: every page lists its language
// variants as hreflang alternates so search engines index all three.
func (m *Module) handleSitemap(w http.ResponseWriter, r *http.Request) {
	site := m.siteURL()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")

	emit := func(path string, mod time.Time) {
		b.WriteString("<url><loc>")
		b.WriteString(seoURL(site, path, LangRU))
		b.WriteString("</loc>")
		for _, l := range Langs {
			b.WriteString(`<xhtml:link rel="alternate" hreflang="` + htmlLang(l) + `" href="` + seoURL(site, path, l) + `"/>`)
		}
		b.WriteString(`<xhtml:link rel="alternate" hreflang="x-default" href="` + seoURL(site, path, LangRU) + `"/>`)
		if !mod.IsZero() {
			b.WriteString("<lastmod>" + mod.UTC().Format("2006-01-02") + "</lastmod>")
		}
		b.WriteString("</url>\n")
	}

	now := time.Now()
	emit("/", now)
	for _, p := range []string{"/about", "/guide", "/formatting", "/pricing", "/support", "/listings", "/author/sana"} {
		emit(p, time.Time{})
	}
	for _, c := range Categories {
		emit("/?cat="+c, time.Time{})
	}
	if arts, err := m.store.SitemapArticles(r.Context()); err != nil {
		m.rt.Logger.Error("sitemap articles", zap.Error(err))
	} else {
		for _, a := range arts {
			emit("/read/"+a.Slug, a.Updated)
		}
	}
	if lst, err := m.listings.SitemapListings(r.Context()); err != nil {
		m.rt.Logger.Error("sitemap listings", zap.Error(err))
	} else {
		for _, l := range lst {
			emit("/listings/"+l.Slug, l.Updated)
		}
	}

	b.WriteString("</urlset>\n")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(b.String()))
}

func absURL(site, u string) string {
	if u == "" {
		return ""
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return site + u
}

func clip(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return strings.TrimSpace(s[:n]) + "…"
}

// applyArticleSEO fills the page's meta description, social image, and a
// NewsArticle JSON-LD block from the already-populated article fields.
func (m *Module) applyArticleSEO(page *ArticlePage) {
	page.OGType = "article"
	if page.Summary != "" {
		page.Desc = clip(page.Summary, 200)
	}
	if img := absURL(page.SiteURL, page.CoverURL); img != "" {
		page.OGImage = img
	}
	canonical := page.SiteURL + page.Path + "?lang=" + page.Lang
	ld := map[string]any{
		"@context":         "https://schema.org",
		"@type":            "NewsArticle",
		"headline":         page.Title,
		"description":      page.Desc,
		"inLanguage":       htmlLang(page.ServedLang),
		"author":           map[string]any{"@type": "Person", "name": page.AuthorName},
		"image":            page.OGImage,
		"mainEntityOfPage": canonical,
		"publisher": map[string]any{
			"@type": "Organization", "name": "Shanraq",
			"logo": map[string]any{"@type": "ImageObject", "url": page.SiteURL + "/static/brand/shanraq.svg"},
		},
	}
	if page.Category != "" {
		ld["articleSection"] = T(page.Lang, "cat."+page.Category)
	}
	if page.Published != nil {
		ld["datePublished"] = page.Published.UTC().Format(time.RFC3339)
	}
	page.JSONLD = jsonLD(ld)
}

// applyListingSEO fills meta description, social image, and a Product JSON-LD
// block (with a price offer) for a listing detail page.
func (m *Module) applyListingSEO(page *ListingViewPage) {
	l := page.L
	page.OGType = "article"
	desc := l.Description
	if desc == "" {
		desc = l.Title + " — " + l.Location()
	}
	page.Desc = clip(desc, 200)
	img := ""
	if len(l.Images) > 0 {
		img = l.Images[0]
	} else {
		img = l.CoverURL
	}
	if a := absURL(page.SiteURL, img); a != "" {
		page.OGImage = a
	}
	ld := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Product",
		"name":        l.Title,
		"description": page.Desc,
		"image":       page.OGImage,
		"category":    T(page.Lang, "re.type_"+l.PropertyType),
		"offers": map[string]any{
			"@type":         "Offer",
			"price":         l.Price,
			"priceCurrency": "KZT",
			"availability":  "https://schema.org/InStock",
			"url":           page.SiteURL + page.Path + "?lang=" + page.Lang,
		},
	}
	page.JSONLD = jsonLD(ld)
}

// SitemapItem is a URL entry for the sitemap.
type SitemapItem struct {
	Slug    string
	Updated time.Time
}

// SitemapArticles returns published article slugs with their last update.
func (s *Store) SitemapArticles(ctx context.Context) ([]SitemapItem, error) {
	rows, err := s.db.Query(ctx, `SELECT slug, updated_at FROM articles
		WHERE status = 'published' ORDER BY updated_at DESC LIMIT 5000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SitemapItem{}
	for rows.Next() {
		var it SitemapItem
		if err := rows.Scan(&it.Slug, &it.Updated); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// SitemapListings returns active listing ids with their creation time.
func (s *ListingStore) SitemapListings(ctx context.Context) ([]SitemapItem, error) {
	rows, err := s.db.Query(ctx, `SELECT id, created_at FROM listings
		WHERE status = 'published' AND expires_at > NOW() ORDER BY created_at DESC LIMIT 5000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SitemapItem{}
	for rows.Next() {
		var it SitemapItem
		if err := rows.Scan(&it.Slug, &it.Updated); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
