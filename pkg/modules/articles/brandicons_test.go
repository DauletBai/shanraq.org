package articles

import (
	"strings"
	"testing"
)

// Every brand the analytics cards can name must have a mark, or the column is
// ragged: some rows indented by an icon, some not.
func TestEveryAnalyticsBrandHasAMark(t *testing.T) {
	// Slugs the trackers emit that name a company. Non-brands — direct, other,
	// mobile, seo, ai — are deliberately absent and render no icon.
	brands := []string{
		"google", "yandex", "bing", "duckduckgo", "facebook", "linkedin", "twitter",
		"youtube", "telegram", "whatsapp", "instagram",
		"chrome", "safari", "firefox", "edge", "opera", "samsung",
		"windows", "android", "linux", "chromeos", "apple", "ios", "macos",
	}
	for _, b := range brands {
		icon := string(brandIcon(b))
		if icon == "" {
			t.Errorf("%q has no brand mark", b)
			continue
		}
		if !strings.HasPrefix(icon, `<svg class="brandico`) || !strings.HasSuffix(icon, "</svg>") {
			t.Errorf("%q produced something that is not a complete svg: %s", b, icon)
		}
		// Either the brand's own hex, or currentColor for the marks whose real
		// colour is near-black — those carry the themed class instead.
		if !strings.Contains(icon, `fill="#`) && !strings.Contains(icon, "brandico--ink") {
			t.Errorf("%q is neither in its brand colour nor themed ink", b)
		}
	}
}

// Rows that name no company must render nothing, so the label starts where the
// icons do rather than being pushed by a placeholder.
func TestNonBrandRowsHaveNoMark(t *testing.T) {
	for _, s := range []string{"direct", "other", "mobile", "tablet", "desktop", "seo", "ai", "email", "share", ""} {
		if got := brandIcon(s); got != "" {
			t.Errorf("%q should have no mark, got %s", s, got)
		}
	}
}

// A black dot tells the reader nothing; every mark carries a glyph or a shape.
func TestMarksAreNotBareDiscs(t *testing.T) {
	for slug, mark := range brandMarks {
		if strings.Count(mark.body, "<") == 1 && strings.Contains(mark.body, "<circle") {
			t.Errorf("%q is a bare coloured circle with nothing on it", slug)
		}
	}
}

// A mark drawn on the 24×24 grid inside a 16×16 viewBox would render as a
// quarter of the logo, cropped to the top-left corner. Since the two grids are
// mixed here on purpose, check every mark declares the grid it was drawn on.
func TestMarkGridMatchesArtwork(t *testing.T) {
	for slug, mark := range brandMarks {
		fetched := strings.Contains(mark.body, `<path fill=`)
		switch {
		case fetched && mark.box != "0 0 24 24":
			t.Errorf("%q is vendored artwork but declares viewBox %q", slug, mark.box)
		case !fetched && mark.box != "0 0 16 16":
			t.Errorf("%q is drawn here but declares viewBox %q", slug, mark.box)
		}
	}
}

// An ink mark that also carries a hex would ignore the theme and disappear on
// the half of it the hex was not chosen for.
func TestInkMarksCarryNoHex(t *testing.T) {
	for slug, mark := range brandMarks {
		if mark.ink && strings.Contains(mark.body, `fill="#`) {
			t.Errorf("%q is themed ink yet hard-codes a colour", slug)
		}
		if !mark.ink && strings.Contains(mark.body, "currentColor") {
			t.Errorf("%q draws in currentColor without asking for the ink class", slug)
		}
	}
}
