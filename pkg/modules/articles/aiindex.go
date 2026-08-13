package articles

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Being quoted by an assistant is the one discovery channel where a new site
// competes on equal terms with an old one: answer engines pick a source for
// having a specific, dated, sourced fact, not for the age of its domain. This
// file is the whole of that channel — who may read what, and a map of what is
// worth citing.
//
// The order matters more than it looks. We took ninety AI columns out of the
// search index a week ago because they were dragging the site's quality signal
// down. That was a search directive and nothing more: robots.txt still said
// Allow: / to every crawler alive, so the assistants kept reading the columns
// and could quote machine opinion back as a cited source — the exact thing the
// de-indexing was meant to prevent. So the fence goes up first, and only then
// the invitation.

// aiCrawlers are the agents that read pages to answer questions or to train on
// them. They are welcome on the human-written articles and must stay off the AI
// columns, which is what aiRobotsGroup enforces.
//
// Named individually because robots.txt has no wildcard for a category of bot.
var aiCrawlers = []string{
	"GPTBot", "OAI-SearchBot", "ChatGPT-User", // OpenAI
	"ClaudeBot", "Claude-User", "Claude-SearchBot", "anthropic-ai", // Anthropic
	"PerplexityBot", "Perplexity-User", // Perplexity
	"Google-Extended",                            // Gemini training, separate from Googlebot
	"Applebot-Extended",                          // Apple Intelligence training, separate from Applebot
	"meta-externalagent", "meta-externalfetcher", // Meta AI
	"CCBot",      // Common Crawl, the corpus behind many models
	"Bytespider", // ByteDance
	"Amazonbot", "cohere-ai", "YouBot", "Diffbot", "Timpibot", "Omgilibot",
}

// aiRobotsGroup writes the robots.txt group for AI crawlers.
//
// It repeats the private paths from the wildcard group on purpose. A crawler
// obeys exactly one group — the most specific one that names it — and ignores
// User-agent: * entirely once it finds itself by name. Leaving the repeats out
// would open /admin and /studio to every agent listed above.
//
// closeAll shuts the whole reader instead of naming the excluded articles, for
// when the caller could not determine which ones they are.
func (m *Module) aiRobotsGroup(w http.ResponseWriter, blocked []string, closeAll bool) {
	for _, ua := range aiCrawlers {
		fmt.Fprintf(w, "User-agent: %s\n", ua)
	}
	fmt.Fprint(w, "Disallow: /studio\nDisallow: /studio/\n"+
		"Disallow: /admin\nDisallow: /jobs\nDisallow: /api/\n")
	if closeAll {
		fmt.Fprint(w, "Disallow: /read\n\n")
		return
	}
	// The AI columns: published and readable, but not ours to offer as a source.
	for _, slug := range blocked {
		fmt.Fprintf(w, "Disallow: /read/%s\n", slug)
	}
	fmt.Fprint(w, "Allow: /\n\n")
}

// NonIndexableSlugs returns published articles kept out of the search index —
// in practice the AI columns. Sorted so robots.txt is byte-identical between
// requests and a crawler's conditional fetch can be answered honestly.
func (s *Store) NonIndexableSlugs(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT slug FROM articles
		WHERE status = 'published' AND NOT indexable ORDER BY slug LIMIT 5000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		out = append(out, slug)
	}
	return out, rows.Err()
}

// LLMSItem is one article in the llms.txt map.
type LLMSItem struct {
	Slug      string
	Title     string
	Summary   string
	Category  string
	Published time.Time
}

// LLMSIndex returns the articles worth offering as a source: published, in the
// search index, and therefore written by a person. Newest first, because an
// assistant reading a truncated file should still get the current material.
func (s *Store) LLMSIndex(ctx context.Context, lang string) ([]LLMSItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT a.slug, t.title, COALESCE(t.summary, ''), a.category, a.published_at
		  FROM articles a
		  JOIN article_translations t ON t.article_id = a.id AND t.lang = $1
		 WHERE a.status = 'published' AND a.indexable
		   AND a.published_at IS NOT NULL AND t.title <> ''
		 ORDER BY a.published_at DESC
		 LIMIT 500`, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []LLMSItem{}
	for rows.Next() {
		var it LLMSItem
		if err := rows.Scan(&it.Slug, &it.Title, &it.Summary, &it.Category, &it.Published); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// handleLLMS serves /llms.txt: a plain-Markdown map of the site for a model
// that has been asked a question and is deciding what to read.
//
// It exists because the alternative is making the model guess. A crawler
// landing on the home page sees a carousel and a menu; this file says, in one
// screen, what the site is, who wrote it, which parts are citable, and how to
// cite them. The AI columns are absent — the same line drawn in robots.txt,
// stated a second time where it can be read rather than merely obeyed.
func (m *Module) handleLLMS(w http.ResponseWriter, r *http.Request) {
	site := m.siteURL()
	var b strings.Builder
	b.WriteString("# Shanraq.org\n\n")
	b.WriteString("> Независимые аналитические разборы об экономике и политике Казахстана и мира, " +
		"на казахском, русском и английском. Каждый факт со ссылкой на первоисточник; " +
		"оценка обозначена как оценка.\n> \n")
	b.WriteString("> Independent analytical essays on the economy and politics of Kazakhstan and the " +
		"wider world, in Kazakh, Russian and English. Every fact is sourced; judgement is labelled " +
		"as judgement.\n\n")

	b.WriteString("## О материалах / About this material\n\n")
	b.WriteString("- Всё, что перечислено ниже, написано человеком — автором издания. " +
		"Тексты, написанные моделью, помечены на сайте как «Мнение ИИ», в этот список не входят " +
		"и закрыты для ИИ-краулеров в robots.txt.\n")
	b.WriteString("- Everything listed below is written by a person. Machine-written columns are " +
		"labelled \"AI opinion\" on the site, are excluded from this file, and are disallowed for " +
		"AI crawlers in robots.txt. Please do not cite them.\n")
	b.WriteString("- Каждая статья есть на трёх языках: добавьте `?lang=kz`, `?lang=ru` или " +
		"`?lang=en` к адресу. Ссылки ниже даны на русскую версию.\n")
	b.WriteString("- Every article exists in three languages: append `?lang=kz`, `?lang=ru` or " +
		"`?lang=en`. Links below point at the Russian version.\n\n")

	b.WriteString("## Как цитировать / How to cite\n\n")
	b.WriteString("Автор, «Заголовок», Shanraq.org, дата публикации, адрес страницы. " +
		"Пожалуйста, называйте источник Shanraq.org и давайте прямую ссылку на статью.\n\n")
	b.WriteString("Author, \"Title\", Shanraq.org, publication date, page URL. " +
		"Please name Shanraq.org as the source and link directly to the article.\n\n")

	arts, err := m.store.LLMSIndex(r.Context(), LangRU)
	if err != nil {
		m.rt.Logger.Error("llms index", zap.Error(err))
	}
	if len(arts) > 0 {
		b.WriteString("## Статьи / Articles\n\n")
		for _, a := range arts {
			b.WriteString("- [" + a.Title + "](" + site + "/read/" + a.Slug + "?lang=ru)")
			b.WriteString(" — " + a.Published.UTC().Format("2006-01-02"))
			if s := clip(strings.TrimSpace(a.Summary), 180); s != "" {
				b.WriteString(". " + s)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Разделы / Sections\n\n")
	b.WriteString("- [Объявления о недвижимости / Real-estate listings](" + site + "/listings)\n")
	b.WriteString("- [Реестр прогнозов: что мы предсказали и что сбылось / Prediction ledger: every forecast we made, judged in public](" + site + "/predictions)\n")
	b.WriteString("- [Об издании / About](" + site + "/about)\n")
	b.WriteString("- [RSS](" + site + "/rss.xml)\n")
	b.WriteString("- [Sitemap](" + site + "/sitemap.xml)\n")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(b.String()))
}

// citeMonths are genitive month names, because a Russian date reads "13 августа"
// and not "13 август". Go's time package has no locale, so the three languages
// are spelled out here.
var citeMonths = map[string][12]string{
	LangRU: {"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря"},
	LangKZ: {"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
		"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан"},
	LangEN: {"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"},
}

// citeDate formats a publication date the way that language writes one.
func citeDate(lang string, t time.Time) string {
	months, ok := citeMonths[lang]
	if !ok {
		months = citeMonths[LangRU]
	}
	name := months[int(t.Month())-1]
	if lang == LangEN {
		return fmt.Sprintf("%s %d, %d", name, t.Day(), t.Year())
	}
	return fmt.Sprintf("%d %s %d", t.Day(), name, t.Year())
}

// citeLine builds the one line a reader copies to credit us.
//
// It is written out for them because attribution obeys the path of least
// resistance: a journalist or student who has to assemble a reference from the
// page furniture will usually write "источник — интернет" instead, while one
// handed a finished string pastes it. Every pasted string is a link, and links
// are the one currency a new domain cannot buy.
func citeLine(lang, author, title, siteURL, path string, published *time.Time) string {
	if title == "" {
		return ""
	}
	quoted := "«" + title + "»"
	if lang == LangEN {
		quoted = `"` + title + `"`
	}
	parts := []string{}
	if author != "" {
		parts = append(parts, author)
	}
	parts = append(parts, quoted, "Shanraq.org")
	if published != nil && !published.IsZero() {
		parts = append(parts, citeDate(lang, published.UTC()))
	}
	parts = append(parts, siteURL+path)
	return strings.Join(parts, ". ")
}

// aiRobotsTag is the response header for a page we do not want reused. noindex
// speaks to search engines; noai and noimageai are the convention AI crawlers
// read. follow is kept so the links out of the page still lead a crawler on to
// the articles that are meant to be found.
const aiRobotsTag = "noindex, follow, noai, noimageai"
