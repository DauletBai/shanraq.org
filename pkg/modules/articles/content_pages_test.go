package articles

import "testing"

// Every editable page must resolve to real built-in content in all three
// languages — that content is both the seed for content_pages and the fallback
// the reader serves if a row is ever missing, so a gap here would ship a blank
// legal page.
func TestEditablePageKeysHaveDefaults(t *testing.T) {
	for _, key := range editablePageKeys {
		for _, lang := range []string{"kz", "ru", "en"} {
			c := staticContent(key, lang)
			if c.Title == "" || c.Body == "" {
				t.Errorf("built-in page %q/%s is empty (title=%q, body len=%d)", key, lang, c.Title, len(c.Body))
			}
		}
	}
}

func TestIsEditablePage(t *testing.T) {
	for _, key := range editablePageKeys {
		if !isEditablePage(key) {
			t.Errorf("isEditablePage(%q) = false, want true", key)
		}
	}
	for _, bad := range []string{"", "home", "admin", "listings", "../terms"} {
		if isEditablePage(bad) {
			t.Errorf("isEditablePage(%q) = true, want false", bad)
		}
	}
}
