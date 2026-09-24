package site

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// The frame must stand on its own: everything chrome.html calls has to come
// from BaseFuncs, or a module that renders a page without articles registered
// would fail to boot. This is the test that keeps the dependency pointing one
// way -- modules know about the site, the site knows nothing about a module.
func TestFrameParsesWithoutModules(t *testing.T) {
	r := NewRenderer()
	if err := r.Build(); err != nil {
		t.Fatalf("the site frame does not parse on its own: %v", err)
	}
	for _, name := range []string{"site_head", "site_header", "site_footer", "sidebar"} {
		if !r.Lookup(name) {
			t.Errorf("frame partial %q is missing", name)
		}
	}
}

// A page rendered before Build must say so, rather than serving an empty body
// with a 200.
func TestRenderBeforeBuildFails(t *testing.T) {
	r := NewRenderer()
	if err := r.Execute(io.Discard, "site_head", nil); err == nil {
		t.Fatal("execute before Build must fail")
	}
}

// A published course must be reachable from the footer in the reader's language.
func TestRustFooterLinkPreservesLanguage(t *testing.T) {
	r := NewRenderer()
	if err := r.Build(); err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"kz", "ru", "en"} {
		t.Run(lang, func(t *testing.T) {
			var body bytes.Buffer
			if err := r.Execute(&body, "site_footer", Base{Lang: lang}); err != nil {
				t.Fatal(err)
			}
			want := `href="/course/rust?lang=` + lang + `"`
			if !strings.Contains(body.String(), want) {
				t.Fatalf("footer missing Rust link %s", want)
			}
			newCourse := `href="/course/kazakh-ai?lang=` + lang + `"`
			if !strings.Contains(body.String(), newCourse) {
				t.Fatalf("footer missing Kazakh AI link %s", newCourse)
			}
			shortNames := map[string]string{"kz": "ЖИ курсы", "ru": "Курс ИИ", "en": "AI course"}
			if !strings.Contains(body.String(), `/static/brand/shanraq.svg`) {
				t.Fatal("footer missing Shanraq course mark")
			}
			if !strings.Contains(body.String(), shortNames[lang]) {
				t.Fatalf("footer missing short course name %q", shortNames[lang])
			}
			if strings.Contains(body.String(), `class="foot-course--soon"`) {
				t.Fatal("published courses still marked as coming soon")
			}
		})
	}
}
