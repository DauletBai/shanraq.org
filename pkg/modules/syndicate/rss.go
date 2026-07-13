package syndicate

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

var rssLangs = map[string]bool{"kz": true, "ru": true, "en": true}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// feedEntry is one article resolved for a language.
type feedEntry struct {
	Slug     string
	Title    string
	Summary  string
	Lang     string
	Modified time.Time
}

func (m *Module) handleRSS(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if !rssLangs[lang] {
		lang = "ru"
	}

	entries, err := m.fetchFeed(r.Context(), lang, 30)
	if err != nil {
		m.log.Error("rss fetch failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	xmlBytes, err := m.renderRSS(lang, entries)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	_, _ = w.Write(xmlBytes)
}

func (m *Module) renderRSS(lang string, entries []feedEntry) ([]byte, error) {
	items := make([]rssItem, 0, len(entries))
	for _, e := range entries {
		items = append(items, rssItem{
			Title:       e.Title,
			Link:        m.articleURL(e.Slug, e.Lang),
			GUID:        m.articleURL(e.Slug, e.Lang),
			Description: e.Summary,
			PubDate:     e.Modified.UTC().Format(time.RFC1123Z),
		})
	}
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       "Shanraq — үй, мұнда еркін дауыстар тоғысады",
			Link:        m.baseURL + "/read?lang=" + lang,
			Description: "Аналитика, ой-пікірлер және оқиғалар үш тілде. Аналитика, мнения и истории на трёх языках.",
			Language:    lang,
			Items:       items,
		},
	}
	body, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal rss: %w", err)
	}
	return append([]byte(xml.Header), body...), nil
}

func (m *Module) articleURL(slug, lang string) string {
	return fmt.Sprintf("%s/read/%s?lang=%s", m.baseURL, slug, lang)
}

// fetchFeed loads published articles resolved to the requested language, falling
// back to the original-language title/summary when a translation is missing.
func (m *Module) fetchFeed(ctx context.Context, lang string, limit int) ([]feedEntry, error) {
	if limit <= 0 || limit > 60 {
		limit = 30
	}
	rows, err := m.db.Query(ctx, `
		SELECT a.slug,
		       COALESCE(NULLIF(tl.title, ''), torig.title)     AS title,
		       COALESCE(NULLIF(tl.summary, ''), torig.summary)  AS summary,
		       CASE WHEN tl.title IS NOT NULL AND tl.title <> '' THEN $1 ELSE a.original_lang END AS lang,
		       COALESCE(a.published_at, a.updated_at)           AS modified
		FROM articles a
		JOIN article_translations torig
		     ON torig.article_id = a.id AND torig.lang = a.original_lang
		LEFT JOIN article_translations tl
		     ON tl.article_id = a.id AND tl.lang = $1 AND tl.title <> '' AND tl.body_md <> ''
		WHERE a.status = 'published'
		ORDER BY a.published_at DESC NULLS LAST
		LIMIT $2
	`, lang, limit)
	if err != nil {
		return nil, fmt.Errorf("query feed: %w", err)
	}
	defer rows.Close()

	var entries []feedEntry
	for rows.Next() {
		var e feedEntry
		if err := rows.Scan(&e.Slug, &e.Title, &e.Summary, &e.Lang, &e.Modified); err != nil {
			return nil, fmt.Errorf("scan feed row: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
