package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

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

// handleRobots serves robots.txt: search engines crawl everything except the
// author cabinet; commercial SEO scanners (which crawl heavily and return
// nothing to us) are turned away; the sitemaps are advertised.
//
// AI crawlers get a group of their own, welcome everywhere except the columns
// written by a model — see aiRobotsGroup. Being quoted in an AI answer is a
// real discovery channel; being quoted out of a machine-written column would
// put our name behind a machine's opinion, which is the one thing this site
// says it does not do.
func (m *Module) handleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	site := m.siteURL()
	// Commercial SEO scanners: pure parasites (feed third-party SEO databases,
	// send us no traffic) and the heaviest crawlers we see — refuse them.
	seoBots := []string{"AhrefsBot", "SemrushBot", "MJ12bot", "DotBot", "BLEXBot", "DataForSeoBot", "PetalBot", "MegaIndex"}
	// /jobs and /admin answer crawlers with 401 and 303 — nothing to index, and
	// Search Console's crawl report showed Googlebot spending requests on them.
	// On a site Google visits about ten times a day, that is worth reclaiming.
	fmt.Fprint(w, "User-agent: *\nAllow: /\n"+
		"Disallow: /studio\nDisallow: /studio/\n"+
		"Disallow: /admin\nDisallow: /jobs\nDisallow: /api/\n\n")
	for _, b := range seoBots {
		fmt.Fprintf(w, "User-agent: %s\nDisallow: /\n\n", b)
	}
	// Read from the same indexable flag that took these out of the search index,
	// so the two can never drift apart.
	//
	// A failed query closes /read entirely rather than emitting a short list.
	// The two mistakes are not equal: hiding the human articles from assistants
	// for one crawl costs a little reach and is undone by the next fetch, while
	// exposing one machine-written column puts our name behind a machine's
	// opinion inside a system with no retraction.
	blocked, err := m.store.NonIndexableSlugs(r.Context())
	if err != nil {
		m.rt.Logger.Error("robots: non-indexable slugs, closing /read to AI crawlers", zap.Error(err))
		blocked = nil
	}
	m.aiRobotsGroup(w, blocked, err != nil)
	fmt.Fprintf(w, "Sitemap: %s/sitemap.xml\nSitemap: %s/sitemap-listings.xml\nSitemap: %s/sitemap-news.xml\n", site, site, site)
}

func seoURL(site, path, lang string) string {
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return xmlEscape(site + path + sep + "lang=" + lang)
}

// sitemapDoc wraps a caller-supplied body in the <urlset> envelope. The emit
// closure writes ONE <url> PER LANGUAGE, each carrying the full set of hreflang
// alternates including itself, so every sitemap we serve shares identical
// formatting.
//
// One entry per language, not one entry per page with alternates hanging off
// it: Google's specification is explicit that you "create a separate <url>
// element for each URL" and that each must list "every alternate version of the
// page, including itself". We used to emit only the Russian <loc>, which left
// the Kazakh and English versions never actually submitted for crawling — on a
// trilingual site that is two thirds of the catalogue relying on being noticed
// sideways. The self-reference also satisfies the reciprocity rule; without it
// Google is entitled to ignore the annotations altogether.
func (m *Module) sitemapDoc(build func(emit func(path string, mod time.Time))) []byte {
	site := m.siteURL()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	emit := func(path string, mod time.Time) {
		for _, self := range Langs {
			b.WriteString("<url><loc>")
			b.WriteString(seoURL(site, path, self))
			b.WriteString("</loc>")
			for _, l := range Langs {
				b.WriteString(`<xhtml:link rel="alternate" hreflang="`)
				b.WriteString(htmlLang(l))
				b.WriteString(`" href="`)
				b.WriteString(seoURL(site, path, l))
				b.WriteString(`"/>`)
			}
			b.WriteString(`<xhtml:link rel="alternate" hreflang="x-default" href="`)
			b.WriteString(seoURL(site, path, LangRU))
			b.WriteString(`"/>`)
			if !mod.IsZero() {
				b.WriteString("<lastmod>")
				b.WriteString(mod.UTC().Format("2006-01-02"))
				b.WriteString("</lastmod>")
			}
			b.WriteString("</url>\n")
		}
	}
	build(emit)
	b.WriteString("</urlset>\n")
	return []byte(b.String())
}

// handleSitemap emits the main trilingual sitemap: home, static pages, category
// feeds and articles. Real-estate listings live in their own sitemap
// (handleSitemapListings) so Search Console tracks the classifieds separately.
func (m *Module) handleSitemap(w http.ResponseWriter, r *http.Request) {
	doc := m.sitemapDoc(func(emit func(path string, mod time.Time)) {
		emit("/", time.Now())
		for _, p := range []string{"/about", "/guide", "/formatting", "/pricing", "/support", "/listings", "/predictions", "/analytics", "/rates", "/author/sana"} {
			emit(p, time.Time{})
		}
		fresh, ferr := m.store.CategoryFreshness(r.Context())
		if ferr != nil {
			m.rt.Logger.Warn("sitemap category freshness", zap.Error(ferr))
		}
		for _, c := range Categories {
			emit("/?cat="+c, fresh[c])
		}
		// Pages of places that actually have something published. Eight hundred
		// empty ones would be eight hundred thin pages offered to a crawler
		// that visits us twenty times a day; the ones with articles are worth
		// its time. Ancestors are included, so the oblast page is offered even
		// when only its village has been written about.
		if places, perr := m.store.PlacesWithArticles(r.Context()); perr != nil {
			m.rt.Logger.Warn("sitemap places", zap.Error(perr))
		} else {
			for _, slug := range places {
				emit("/place/"+slug, time.Time{})
			}
		}
		if arts, err := m.store.SitemapArticles(r.Context()); err != nil {
			m.rt.Logger.Error("sitemap articles", zap.Error(err))
		} else {
			for _, a := range arts {
				emit("/read/"+a.Slug, a.Updated)
			}
		}
	})
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write(doc)
}

// handleSitemapNews emits the Google News sitemap: only articles published in
// the last 48 hours, in the language each was written in.
//
// A separate format because it answers a different question. The main sitemap
// says "these pages exist"; the news sitemap says "this went up in the last two
// days", which is the window Google News and Discover read. Google requires
// entries be dropped after 48 hours, so this file is usually short and often
// empty — an empty <urlset> is valid and is the correct answer on a quiet day.
func (m *Module) handleSitemapNews(w http.ResponseWriter, r *http.Request) {
	site := m.siteURL()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" ` +
		`xmlns:news="http://www.google.com/schemas/sitemap-news/0.9">` + "\n")
	if arts, err := m.store.NewsSitemapArticles(r.Context()); err != nil {
		m.rt.Logger.Error("news sitemap", zap.Error(err))
	} else {
		for _, a := range arts {
			b.WriteString("<url><loc>")
			b.WriteString(seoURL(site, "/read/"+a.Slug, a.Lang))
			b.WriteString("</loc><news:news><news:publication><news:name>Shanraq.org</news:name>")
			b.WriteString("<news:language>" + htmlLang(a.Lang) + "</news:language></news:publication>")
			b.WriteString("<news:publication_date>" + a.Published.UTC().Format(time.RFC3339) + "</news:publication_date>")
			b.WriteString("<news:title>" + xmlEscape(a.Title) + "</news:title>")
			b.WriteString("</news:news></url>\n")
		}
	}
	b.WriteString("</urlset>\n")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(b.String()))
}

// handleSitemapListings emits a sitemap of only real-estate listing detail
// pages, so the classifieds can be submitted and tracked as their own section
// in Google Search Console, separate from editorial content.
func (m *Module) handleSitemapListings(w http.ResponseWriter, r *http.Request) {
	doc := m.sitemapDoc(func(emit func(path string, mod time.Time)) {
		if lst, err := m.listings.SitemapListings(r.Context()); err != nil {
			m.rt.Logger.Error("sitemap listings", zap.Error(err))
		} else {
			for _, l := range lst {
				emit("/listings/"+l.Slug, l.Updated)
			}
		}
	})
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write(doc)
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

// clip shortens a string to at most n runes. It counts runes, not bytes:
// Russian and Kazakh text is two bytes per letter, so a byte slice would cut a
// letter in half and emit invalid UTF-8 into the meta description. It also
// backs up to the last space so the description does not end mid-word.
func clip(s string, n int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	cut, i := len(s), 0
	for j := range s { // ranging a string yields rune start offsets
		if i == n {
			cut = j
			break
		}
		i++
	}
	out := strings.TrimSpace(s[:cut])
	if sp := strings.LastIndexAny(out, "  "); sp > n/2 {
		out = strings.TrimSpace(out[:sp])
	}
	return out + "…"
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
		"author":           authorLD(page),
		"image":            page.OGImage,
		"mainEntityOfPage": canonical,
		"publisher": map[string]any{
			"@type": "Organization", "name": "Shanraq.org",
			"logo": map[string]any{"@type": "ImageObject", "url": page.SiteURL + "/static/brand/shanraq.svg"},
		},
	}
	if page.Category != "" {
		ld["articleSection"] = T(page.Lang, "cat."+page.Category)
	}
	if page.Published != nil {
		ld["datePublished"] = page.Published.UTC().Format(time.RFC3339)
		// Without dateModified an updated article looks as old as the day it was
		// published, and freshness is one of the few levers a small publisher has.
		mod := page.Published
		if page.Updated != nil && page.Updated.After(*page.Published) {
			mod = page.Updated
		}
		ld["dateModified"] = mod.UTC().Format(time.RFC3339)
	}
	// Two blocks in one array: schema.org allows it and Google reads it. The
	// breadcrumb is what puts a "Shanraq.org › Мир › Европа" trail under the
	// result instead of a bare URL.
	page.JSONLD = jsonLD([]any{ld, breadcrumbLD(page)})
}

// breadcrumbLD builds the trail shown under a search result: site → category →
// subcategory → this article. Only the levels that exist are emitted, so an
// article filed without a subcategory produces a three-step trail rather than
// one with a hole in it.
func breadcrumbLD(page *ArticlePage) map[string]any {
	items := []map[string]any{{
		"@type": "ListItem", "position": 1,
		"name": "Shanraq.org", "item": page.SiteURL + "/?lang=" + page.Lang,
	}}
	if page.Category != "" {
		items = append(items, map[string]any{
			"@type": "ListItem", "position": len(items) + 1,
			"name": T(page.Lang, "cat."+page.Category),
			"item": page.SiteURL + "/?cat=" + page.Category + "&lang=" + page.Lang,
		})
	}
	if page.Subcategory != "" {
		items = append(items, map[string]any{
			"@type": "ListItem", "position": len(items) + 1,
			"name": T(page.Lang, "sub."+page.Subcategory),
			"item": page.SiteURL + "/?sub=" + page.Subcategory + "&lang=" + page.Lang,
		})
	}
	// The last crumb is the page itself and carries no link, per Google's guidance.
	items = append(items, map[string]any{
		"@type": "ListItem", "position": len(items) + 1, "name": page.Title,
	})
	return map[string]any{
		"@context": "https://schema.org", "@type": "BreadcrumbList", "itemListElement": items,
	}
}

// applyListingSEO fills meta description, social image, and a Product JSON-LD
// block (with a price offer) for a listing detail page.
func (m *Module) applyListingSEO(page *ListingViewPage) {
	l := page.L
	page.OGType = "article"
	desc := l.DescriptionIn(page.Lang)
	if desc == "" {
		desc = l.TitleIn(page.Lang) + " — " + l.Location()
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
	priceCurrency := l.Currency
	if priceCurrency == "" {
		priceCurrency = "KZT"
	}
	ld := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "Product",
		"name":        l.TitleIn(page.Lang),
		"description": page.Desc,
		"image":       page.OGImage,
		"category":    T(page.Lang, "re.type_"+l.PropertyType),
		"offers": map[string]any{
			"@type":         "Offer",
			"price":         l.Price,
			"priceCurrency": priceCurrency,
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

// NewsItem is one entry in the Google News sitemap: the article in the language
// it was written in, with the title and date that format requires.
type NewsItem struct {
	Slug      string
	Lang      string
	Title     string
	Published time.Time
}

// NewsSitemapArticles returns articles published in the last 48 hours — the
// window Google News accepts — one entry per finished language version, each
// with that version's own title.
//
// One row per translation, not one per article. The site publishes the same
// piece in Kazakh, Russian and English at three addresses, and listing only the
// original left two thirds of everything we write outside the window entirely.
//
// A version counts as finished when it has a title and a body: an empty tab is
// an intention, and an intention announced as news is a page a reader opens for
// nothing. Non-indexable articles are excluded on the same grounds as
// everywhere else — a page kept out of the index has no business being pushed
// as news.
func (s *Store) NewsSitemapArticles(ctx context.Context) ([]NewsItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT a.slug, t.lang, t.title, a.published_at
		  FROM articles a
		  JOIN article_translations t ON t.article_id = a.id
		 WHERE a.status = 'published' AND a.indexable
		   AND a.published_at IS NOT NULL
		   AND a.published_at >= NOW() - INTERVAL '48 hours'
		   AND t.title <> '' AND t.body_md <> ''
		   AND (t.status = 'ready' OR t.lang = a.original_lang)
		 ORDER BY a.published_at DESC, t.lang
		 LIMIT 1000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []NewsItem{}
	for rows.Next() {
		var it NewsItem
		if err := rows.Scan(&it.Slug, &it.Lang, &it.Title, &it.Published); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// SitemapArticles returns published article slugs with their last update.
func (s *Store) SitemapArticles(ctx context.Context) ([]SitemapItem, error) {
	rows, err := s.db.Query(ctx, `SELECT slug, updated_at FROM articles
		WHERE status = 'published' AND indexable ORDER BY updated_at DESC LIMIT 5000`)
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

// authorLD builds the author block. The url matters: it is what lets a search
// engine tie a person's articles together instead of treating every piece as
// written by a stranger with the same name.
func authorLD(page *ArticlePage) map[string]any {
	a := map[string]any{"@type": "Person", "name": page.AuthorName}
	if page.AuthorID != "" {
		a["url"] = page.SiteURL + "/author/" + page.AuthorID
	}
	return a
}
