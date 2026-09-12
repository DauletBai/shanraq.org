package site

import (
	"io"
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
