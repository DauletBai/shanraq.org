package articles

import (
	"strings"
	"testing"
)

// A table in an article has to survive a narrow screen without becoming
// something else. The old stylesheet scrolled it by making the table a block,
// which costs it its table semantics and leaves a region only a mouse can
// scroll.
func TestRenderedTablesScrollFromTheKeyboard(t *testing.T) {
	src := "| | Цена | Срок |\n|---|---|---|\n| Аренда | 250 000 | 12 мес |\n"
	out := string(RenderMarkdown(src))

	if !strings.Contains(out, `<div class="tscroll tscroll--prose" tabindex="0"><table>`) {
		t.Errorf("table is not in a focusable scroll wrapper:\n%s", out)
	}
	if !strings.Contains(out, "</table></div>") {
		t.Errorf("scroll wrapper is not closed:\n%s", out)
	}
	if strings.Contains(out, "<th></th>") {
		t.Errorf("an empty header cell survived; it labels every cell under it with a blank:\n%s", out)
	}
	if !strings.Contains(out, "<th>Цена</th>") {
		t.Errorf("a real header cell was demoted:\n%s", out)
	}
	// Balanced: one wrapper per table, and nothing wrapped when there is none.
	if got, want := strings.Count(out, `<div class="tscroll`), 1; got != want {
		t.Errorf("wrapped %d times, want %d", got, want)
	}
	if got, want := strings.Count(out, "</table></div>"), 1; got != want {
		t.Errorf("closed %d wrappers, want %d", got, want)
	}
	if plain := string(RenderMarkdown("просто абзац")); strings.Contains(plain, "tscroll") {
		t.Errorf("prose without a table was wrapped: %s", plain)
	}
}

// The table of contents pass rewrites the same HTML; both renderers must end up
// with the same wrapper, or an article gets it and a page does not.
func TestTOCRenderAlsoWrapsTables(t *testing.T) {
	src := "## Раздел\n\n| A | B |\n|---|---|\n| 1 | 2 |\n"
	out, toc := RenderMarkdownTOC(src)
	if len(toc) != 1 {
		t.Fatalf("toc has %d entries, want 1", len(toc))
	}
	if !strings.Contains(string(out), `tabindex="0"`) {
		t.Errorf("table rendered through the TOC path is not focusable:\n%s", out)
	}
}
