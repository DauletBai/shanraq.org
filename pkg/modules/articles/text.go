package articles

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// md is a shared, safe Markdown renderer. Raw inline HTML is NOT enabled
// (no WithUnsafe), so user-supplied HTML is escaped — our first XSS guard.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Typographer, codeHighlighting),
	goldmark.WithRendererOptions(gmhtml.WithHardWraps()),
)

// codeHighlighting colours fenced code the way an editor does. It runs on the
// server, so there is no script to load and no CDN to be cut off from, and a
// crawler sees exactly what a reader sees.
//
// Classes rather than inline styles: inline colours are fixed at render time
// and could not follow the site's light/dark switch. The palette lives in the
// stylesheet beside every other colour we use. Chroma's class names are very
// short (k, s, nf), so every rule for them is scoped under .chroma.
//
// Guessing is off on purpose. A fence with no language is terminal output in
// these articles, and a guesser paints it as though it were code.
var codeHighlighting = highlighting.NewHighlighting(
	highlighting.WithFormatOptions(
		chromahtml.WithClasses(true),
	),
	highlighting.WithGuessLanguage(false),
)

// tableOpen matches the opening tag goldmark emits for a GFM table, and
// emptyHeader an unlabelled header cell — the corner cell a Markdown table
// leaves blank when its first column is the row label.
var (
	tableOpen   = regexp.MustCompile(`<table>`)
	emptyHeader = regexp.MustCompile(`<th[^>]*>\s*</th>`)
	bareImage   = regexp.MustCompile(`<img `)
)

// deferImages keeps a page with several screenshots from paying for all of them
// before the first paragraph is readable. Goldmark emits a bare <img>; every
// other image on this site is already lazy, and a guide page carries six.
func deferImages(html string) string {
	if !strings.Contains(html, "<img ") {
		return html
	}
	return bareImage.ReplaceAllString(html, `<img loading="lazy" decoding="async" `)
}

// wrapTables makes a rendered table survive a narrow screen without losing what
// it is.
//
// The stylesheet used to scroll the table by making it display:block. That
// works visually and costs the element its table semantics — rows and headers
// stop being announced as a table — and it leaves a scrollable region that no
// keyboard can reach, because a div that scrolls is only scrollable by someone
// holding a mouse. The scroll belongs to a focusable wrapper; the table stays a
// table.
//
// An empty header cell becomes a data cell in the same pass: a <th> with
// nothing in it labels nothing, and a screen reader reads the blank as the
// heading of every cell beneath it.
func wrapTables(html string) string {
	if !strings.Contains(html, "<table>") {
		return html
	}
	html = emptyHeader.ReplaceAllString(html, "<td></td>")
	html = tableOpen.ReplaceAllString(html, `<div class="tscroll tscroll--prose" tabindex="0"><table>`)
	return strings.ReplaceAll(html, "</table>", "</table></div>")
}

// RenderMarkdown converts Markdown to sanitized HTML for templates.
func RenderMarkdown(source string) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(source))
	}
	return template.HTML(deferImages(wrapTables(buf.String()))) //nolint:gosec // goldmark configured without raw HTML
}

// TOCItem is one entry in an article's table of contents.
type TOCItem struct {
	ID   string
	Text string
}

// RenderMarkdownTOC renders Markdown and, in the same pass, extracts a table of
// contents from the level-2 (##) headings. Each rendered <h2> gets a sequential
// id ("sec-1", "sec-2", …) so the TOC anchors line up exactly with the body.
func RenderMarkdownTOC(source string) (template.HTML, []TOCItem) {
	// Collect ## headings in document order (ignore ### and deeper).
	var toc []TOCItem
	for _, line := range strings.Split(source, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "## ") && !strings.HasPrefix(t, "### ") {
			text := strings.TrimSpace(strings.TrimPrefix(t, "##"))
			toc = append(toc, TOCItem{ID: fmt.Sprintf("sec-%d", len(toc)+1), Text: stripInline(text)})
		}
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(source)), nil //nolint:gosec
	}
	out := buf.String()
	// Inject the sequential ids into the rendered <h2> tags, in order.
	for _, it := range toc {
		out = strings.Replace(out, "<h2>", `<h2 id="`+it.ID+`">`, 1)
	}
	return template.HTML(foldAnswers(deferImages(wrapTables(out)))), toc //nolint:gosec // goldmark configured without raw HTML
}

// answersHead matches a lesson's answers heading, which the course writes with
// the same three words in every lesson.
var answersHead = regexp.MustCompile(`<h2 id="sec-\d+">(Ответы|Жауаптар|Answers)</h2>`)

// foldShow is the label on the fold, by the language the heading is written in.
var foldShow = map[string]string{
	"Ответы":   "Показать ответы",
	"Жауаптар": "Жауаптарды көрсету",
	"Answers":  "Show the answers",
}

// foldAnswers hides a lesson's answers behind a disclosure.
//
// The course asks three questions and answers them at the bottom of the same
// page. Read straight through, the answers arrive before the reader has tried
// to produce one, and recall that costs nothing teaches nothing. A <details>
// puts the effort back: the answer is one click away, but the click is a
// decision the reader makes.
//
// Only lessons are folded -- the guard is an exercise heading in the same
// document -- so an ordinary article with a section called "Ответы" is left
// alone. The heading itself stays outside the fold: it carries the table of
// contents anchor, and a link into a closed disclosure would land nowhere.
func foldAnswers(html string) string {
	if !strings.Contains(html, ">Задание</h2>") &&
		!strings.Contains(html, ">Тапсырма</h2>") &&
		!strings.Contains(html, ">Exercise</h2>") {
		return html
	}
	loc := answersHead.FindStringSubmatchIndex(html)
	if loc == nil {
		return html
	}
	head := html[loc[0]:loc[1]]
	word := html[loc[2]:loc[3]]
	body := html[loc[1]:]
	// The fold ends where the next section begins; the sources block below the
	// answers is a section of its own and stays visible.
	if i := strings.Index(body, "<h2"); i >= 0 {
		body = body[:i]
	}
	folded := head + `<details class="fold"><summary class="fold__s">` +
		template.HTMLEscapeString(foldShow[word]) + `</summary><div class="fold__b">` +
		body + `</div></details>`
	return html[:loc[0]] + folded + html[loc[1]+len(body):]
}

// stripInline removes inline Markdown emphasis/link markup from a heading so the
// TOC label reads as plain text.
func stripInline(s string) string {
	r := strings.NewReplacer("**", "", "__", "", "*", "", "`", "", "_", "")
	return strings.TrimSpace(r.Replace(s))
}

// readingMinutes estimates reading time from the plain-text word count
// (~180 wpm), with a one-minute floor.
func readingMinutes(source string) int {
	words := len(strings.Fields(stripMD(source)))
	m := (words + 179) / 180
	if m < 1 {
		m = 1
	}
	return m
}

// stripMD removes the most common Markdown markup to produce plain text
// suitable for a feed summary.
// htmlTag matches a complete tag, so an author who pasted HTML gets it removed
// rather than mangled. The old code only dropped the ">" of each tag, turning
// "<script>alert(1)</script>" into "<scriptalert(1)</script" — which then went
// into card excerpts and, worse, into the meta description of the page. Never an
// XSS (templates escape it), but it is rubbish where a summary should be.
var htmlTag = regexp.MustCompile(`(?s)<[^>]*>`)

func stripMD(s string) string {
	s = htmlTag.ReplaceAllString(s, " ")
	repl := strings.NewReplacer(
		"#", "", "*", "", "_", "", "`", "", ">", "", "~", "",
		"![", "", "](", " ", "]", "", "[", "",
		"<", "", // a stray "<" that was not part of a tag
	)
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		b.WriteString(repl.Replace(line))
		b.WriteByte(' ')
	}
	// Tag removal leaves runs of spaces where the markup used to be.
	return strings.Join(strings.Fields(b.String()), " ")
}

// excerpt trims text to a plain-text summary of at most n runes.
func excerpt(s string, n int) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	cut := string(runes[:n])
	if idx := strings.LastIndex(cut, " "); idx > n/2 {
		cut = cut[:idx]
	}
	return strings.TrimSpace(cut) + "…"
}

// displayName derives a friendly author name from an email address.
func displayName(email string) string {
	at := strings.IndexByte(email, '@')
	if at <= 0 {
		if email == "" {
			return "Автор"
		}
		return email
	}
	local := email[:at]
	local = strings.NewReplacer(".", " ", "_", " ", "-", " ").Replace(local)
	fields := strings.Fields(local)
	for i, f := range fields {
		r := []rune(f)
		r[0] = []rune(strings.ToUpper(string(r[0])))[0]
		fields[i] = string(r)
	}
	if len(fields) == 0 {
		return "Автор"
	}
	return strings.Join(fields, " ")
}

// shortAuthor compresses a byline to initial + family name ("Daulet Baimurza"
// → "D. Baimurza") so the card's meta row has room for the view counter. The
// given name is what gets abbreviated — that is the convention everywhere a
// byline has to fit (news credits, academic citations, wire services), and the
// family name is the half readers actually recognise and search for. A
// single-word name (a persona, an email-derived fallback) is left alone. The
// initial is glued to the surname with a non-breaking space so "D." can never
// end up alone at the end of a line.
func shortAuthor(name string) string {
	fields := strings.Fields(name)
	if len(fields) < 2 {
		return strings.TrimSpace(name)
	}
	var b strings.Builder
	for _, f := range fields[:len(fields)-1] {
		r := []rune(f)
		b.WriteString(strings.ToUpper(string(r[0])))
		b.WriteString(". ")
	}
	b.WriteString(fields[len(fields)-1])
	return b.String()
}

// compactNum shortens a counter that has to live in a tight meta row: under a
// thousand it prints in full, above that it collapses to a localized short form
// (1,2 мың / 1,2 тыс. / 1.2k). One decimal is kept below ten so 1500 doesn't
// flatten to a bare "1 тыс.", and the value is truncated rather than rounded —
// a counter should never claim more views than there were.
func compactNum(lang string, n int64) string {
	if n < 1000 {
		return strconv.FormatInt(n, 10)
	}
	unit, div := "num.thousand", int64(1000)
	if n >= 1_000_000 {
		unit, div = "num.million", 1_000_000
	}
	whole := n / div
	num := strconv.FormatInt(whole, 10)
	if tenth := (n % div) * 10 / div; whole < 10 && tenth > 0 {
		sep := ","
		if lang == LangEN {
			sep = "."
		}
		num += sep + strconv.FormatInt(tenth, 10)
	}
	// "1.2k" is set solid in English; the Cyrillic word forms take a
	// non-breaking space, which also keeps number and unit on one line.
	gap := " "
	if lang == LangEN {
		gap = ""
	}
	return num + gap + T(lang, unit)
}

// translitMap transliterates Kazakh + Russian Cyrillic to latin for slugs.
var translitMap = map[rune]string{
	'а': "a", 'ә': "a", 'б': "b", 'в': "v", 'г': "g", 'ғ': "g", 'д': "d",
	'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i", 'й': "i", 'к': "k",
	'қ': "q", 'л': "l", 'м': "m", 'н': "n", 'ң': "n", 'о': "o", 'ө': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ұ': "u", 'ү': "u",
	'ф': "f", 'х': "h", 'һ': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'і': "i", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// Slugify produces a URL-safe latin slug from a possibly-Cyrillic title.
func Slugify(title string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case translitMap[r] != "":
			b.WriteString(translitMap[r])
			prevDash = false
		default:
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if len(slug) > 70 {
		slug = strings.Trim(slug[:70], "-")
	}
	if slug == "" {
		slug = "article"
	}
	return slug
}
