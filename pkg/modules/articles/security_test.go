package articles

import (
	"strings"
	"testing"
)

// Everything below guards an invariant whose failure is an incident, not a bug:
// stored XSS, a leaked phone number, a gate that stops gating. They are cheap to
// keep and they fail loudly the day someone changes the code beneath them.

// RenderMarkdown is the single most dangerous function in the project: it takes
// text written by any author or commenter and turns it into HTML that every
// reader executes. goldmark is configured WITHOUT WithUnsafe, so raw HTML must
// come out escaped. That flag is one word long and looks harmless — someone will
// eventually reach for it to embed a video, and stored XSS ships the same day.
func TestRenderMarkdownEscapesHTML(t *testing.T) {
	payloads := []string{
		`<script>alert(1)</script>`,
		`<img src=x onerror="alert(1)">`,
		`<iframe src="https://evil.example"></iframe>`,
		`<svg/onload=alert(1)>`,
		`<a href="javascript:alert(1)">click</a>`,
		`<style>body{display:none}</style>`,
		`<div onclick="steal()">text</div>`,
		"Обычный текст <b>с тегом</b> внутри",
	}
	for _, in := range payloads {
		out := string(RenderMarkdown(in))
		low := strings.ToLower(out)
		// The invariant is that no tag from the input survives as a tag. Testing
		// for the substring "onload=" would be wrong: goldmark escapes the whole
		// thing to "&lt;svg/onload=alert(1)&gt;", where those characters are
		// text and cannot execute. What must never appear is an opening angle
		// bracket followed by one of these names.
		for _, tag := range []string{"<script", "<iframe", "<style", "<svg", "<img", "<a ", "<div"} {
			if strings.Contains(low, tag) {
				t.Errorf("RenderMarkdown(%q) emitted a live %q tag:\n%s", in, tag, out)
			}
		}
		if strings.Contains(low, "javascript:") {
			t.Errorf("RenderMarkdown(%q) kept a javascript: URL:\n%s", in, out)
		}
	}

	// goldmark deals with the two cases differently — inline HTML comes back
	// escaped, a block of it is dropped outright — and both are safe. What must
	// not happen is the renderer becoming so blunt it eats ordinary writing, so
	// check that real markdown still works.
	out := string(RenderMarkdown("# Заголовок\n\nАбзац с **выделением** и [ссылкой](/guide)."))
	for _, want := range []string{"<h1", "Заголовок", "<strong>", `href="/guide"`} {
		if !strings.Contains(out, want) {
			t.Errorf("ordinary markdown lost %q:\n%s", want, out)
		}
	}
}

// Markdown links are written by users too. A javascript: or data: destination
// in an ordinary-looking [text](url) is the quiet version of the same attack.
func TestRenderMarkdownRejectsDangerousLinks(t *testing.T) {
	for _, in := range []string{
		`[клик](javascript:alert(1))`,
		`[клик](JaVaScRiPt:alert(1))`,
		`![картинка](javascript:alert(1))`,
	} {
		out := strings.ToLower(string(RenderMarkdown(in)))
		if strings.Contains(out, "javascript:") {
			t.Errorf("RenderMarkdown(%q) kept a javascript: destination:\n%s", in, out)
		}
	}
	// A legitimate link must survive, or the guard is useless in practice.
	if out := string(RenderMarkdown(`[справка](/guide)`)); !strings.Contains(out, `href="/guide"`) {
		t.Errorf("an ordinary link did not survive: %s", out)
	}
}

// The contact is hidden until a reader asks for it, and revealing it is counted.
// If the mask leaks enough characters the counter measures nothing, because the
// number can be read straight off the page.
func TestMaskContactHidesTheNumber(t *testing.T) {
	const phone = "+7 701 234 56 78"
	masked := maskContact(phone)

	if masked == phone {
		t.Fatal("the contact was returned unmasked")
	}
	// The digits in the middle are the ones worth hiding; count how many of the
	// original digits survive. A handful of leading and trailing ones are kept
	// on purpose so the seller recognises their own number.
	var kept, total int
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			total++
		}
	}
	for _, c := range masked {
		if c >= '0' && c <= '9' {
			kept++
		}
	}
	if kept*2 > total {
		t.Errorf("mask kept %d of %d digits (%q) — more than half is not a mask", kept, total, masked)
	}
	// Short strings must not be partially revealed at all.
	for _, short := range []string{"", "1", "12345", "123456"} {
		if got := maskContact(short); strings.ContainsAny(got, "0123456789") {
			t.Errorf("maskContact(%q) = %q — a short contact must be fully hidden", short, got)
		}
	}
}

// The staged-launch gates are the difference between a closed beta and an open
// site. Each status has to behave exactly as the admin panel promises.
func TestServiceGateStatuses(t *testing.T) {
	flags := NewServiceFlags(nil)
	for _, code := range []string{SvcRegistration, SvcArticleSubmit, SvcListingSubmit, SvcComments} {
		// Missing row defaults to open, so a cold start never locks the site.
		flags.cache = map[string]ServiceFlag{}
		if !flags.Flag(code).Available() {
			t.Errorf("%s: an unset flag must default to available", code)
		}
		for _, st := range []string{svcOff, svcMaintenance, svcInviteOnly} {
			flags.cache = map[string]ServiceFlag{code: {Code: code, Status: st}}
			if flags.Flag(code).Available() {
				t.Errorf("%s: status %q must not report available", code, st)
			}
		}
		flags.cache = map[string]ServiceFlag{code: {Code: code, Status: svcOn}}
		if !flags.Flag(code).Available() {
			t.Errorf("%s: status on must report available", code)
		}
	}
}

// The global switch decides whether the site serves at all. Defaulting it the
// wrong way takes everything down on a cold start or a database hiccup.
func TestSiteSwitchFailsOpen(t *testing.T) {
	flags := NewServiceFlags(nil)
	flags.cache = map[string]ServiceFlag{}
	if !flags.SiteUp() {
		t.Error("an unloaded site flag must default to up")
	}
	flags.cache = map[string]ServiceFlag{SvcSite: {Code: SvcSite, Status: svcOff}}
	if flags.SiteUp() {
		t.Error("status off must take the site down")
	}
}

// Category and subcategory slugs reach SQL and URLs. Anything unrecognised must
// be refused rather than normalised into something that looks valid.
func TestCategorySlugsAreClosed(t *testing.T) {
	bad := []string{"", " ", "sport; DROP TABLE articles", "../admin", "SPORT",
		"politics'", "<script>", "general"}
	for _, s := range bad {
		if IsCategory(s) && s != "general" {
			t.Errorf("IsCategory(%q) accepted an unknown slug", s)
		}
		if IsSubcategory(s) {
			t.Errorf("IsSubcategory(%q) accepted an unknown slug", s)
		}
	}
	if NormalizeCategory("sport; DROP TABLE articles") != CategoryGeneral {
		t.Error("an unknown category must normalise to the general bucket")
	}
	for _, c := range Categories {
		if !IsCategory(c) {
			t.Errorf("IsCategory(%q) rejected a registered category", c)
		}
	}
}

// Slugs end up in URLs and in a unique index. Cyrillic must transliterate, and
// nothing that could break a path may survive.
func TestSlugifyIsURLSafe(t *testing.T) {
	for _, in := range []string{
		"Привет мир", "../../etc/passwd", "a/b?c=d#e", "<script>alert(1)</script>",
		"   ", "Ұлттық қор", "100% рост", "a\nb\tc",
	} {
		got := Slugify(in)
		if got == "" {
			t.Errorf("Slugify(%q) produced an empty slug", in)
		}
		for _, c := range got {
			ok := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-'
			if !ok {
				t.Errorf("Slugify(%q) = %q — contains %q, which is not URL-safe", in, got, c)
			}
		}
		if strings.Contains(got, "--") || strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
			t.Errorf("Slugify(%q) = %q — malformed dashes", in, got)
		}
	}
}

// Excerpts are built from article bodies and shown on cards and in meta tags.
// They must be plain text: markdown that survives stripping would render as
// markup in a description tag.
func TestExcerptIsPlainAndBounded(t *testing.T) {
	body := "# Заголовок\n\n**Жирный** текст со [ссылкой](/x) и `кодом`.\n\n<script>alert(1)</script>"
	got := excerpt(stripMD(body), 100)
	for _, marker := range []string{"#", "**", "`", "](", "<script"} {
		if strings.Contains(got, marker) {
			t.Errorf("excerpt kept %q: %q", marker, got)
		}
	}
	if n := len([]rune(got)); n > 101 { // 100 + the ellipsis
		t.Errorf("excerpt is %d runes, want at most 101: %q", n, got)
	}
}
