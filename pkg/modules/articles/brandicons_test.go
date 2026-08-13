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
		if !strings.HasPrefix(icon, `<svg class="brandico"`) || !strings.HasSuffix(icon, "</svg>") {
			t.Errorf("%q produced something that is not a complete svg: %s", b, icon)
		}
		if !strings.Contains(icon, "fill=\"#") {
			t.Errorf("%q has no brand colour", b)
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
		if strings.Count(mark, "<") == 1 && strings.Contains(mark, "<circle") {
			t.Errorf("%q is a bare coloured circle with nothing on it", slug)
		}
	}
}
