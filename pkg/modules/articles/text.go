package articles

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// md is a shared, safe Markdown renderer. Raw inline HTML is NOT enabled
// (no WithUnsafe), so user-supplied HTML is escaped — our first XSS guard.
var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Typographer),
	goldmark.WithRendererOptions(gmhtml.WithHardWraps()),
)

// RenderMarkdown converts Markdown to sanitized HTML for templates.
func RenderMarkdown(source string) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return template.HTML(template.HTMLEscapeString(source))
	}
	return template.HTML(buf.String()) //nolint:gosec // goldmark configured without raw HTML
}

// stripMD removes the most common Markdown markup to produce plain text
// suitable for a feed summary.
func stripMD(s string) string {
	repl := strings.NewReplacer(
		"#", "", "*", "", "_", "", "`", "", ">", "", "~", "",
		"![", "", "](", " ", "]", "", "[", "",
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
	return strings.TrimSpace(b.String())
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
