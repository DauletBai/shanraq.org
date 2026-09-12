package site

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// Markdown is how this site turns text into pages: an article, a lesson, a
// static page and a product description all go through the one renderer, so
// prose written in one part of the site cannot render differently in another.

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
	return template.HTML(deferImages(wrapTables(out))), toc //nolint:gosec // goldmark configured without raw HTML
}

// stripInline removes inline Markdown emphasis/link markup from a heading so the
// TOC label reads as plain text.
func stripInline(s string) string {
	r := strings.NewReplacer("**", "", "__", "", "*", "", "`", "", "_", "")
	return strings.TrimSpace(r.Replace(s))
}
