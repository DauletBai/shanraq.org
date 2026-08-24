package articles

import (
	"context"
	"testing"
)

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

// An edit to a document in the code has to reach the site. The pages are seeded into
// the database on first start, and with an insert-only-if-absent later edits in the
// code never reached the site: a visitor read the first start's text for years.
func TestCorrectingADocumentInCodeReachesTheSite(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	store := NewContentStore(app.pool)
	// The test database lives between runs: a row left over from last time would pass
	// this test even without the fix.
	app.exec(`DELETE FROM content_pages WHERE page_key = 'guide-refresh-probe'`)

	if err := store.seed(ctx, "guide-refresh-probe", "ru", "Руководство", "Старый текст"); err != nil {
		t.Fatalf("первый посев: %v", err)
	}
	if err := store.seed(ctx, "guide-refresh-probe", "ru", "Руководство", "Новый текст"); err != nil {
		t.Fatalf("повторный посев: %v", err)
	}
	_, body, _, err := store.Get(ctx, "guide-refresh-probe", "ru")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if body != "Новый текст" {
		t.Errorf("на сайте остался текст первого запуска: %q", body)
	}
}

// And the seeding does not touch a person's edit: otherwise a restart would erase an
// editor's work, and the admin panel would become a trap.
func TestSeedingNeverOverwritesAHumanEdit(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	store := NewContentStore(app.pool)
	editor := app.createUser("pageeditor@example.com", "Parol123!")
	app.exec(`DELETE FROM content_pages WHERE page_key = 'guide-edit-probe'`)

	if err := store.seed(ctx, "guide-edit-probe", "ru", "Руководство", "Встроенный текст"); err != nil {
		t.Fatalf("посев: %v", err)
	}
	if err := store.Upsert(ctx, "guide-edit-probe", "ru", "Руководство", "Текст редактора", editor); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if err := store.seed(ctx, "guide-edit-probe", "ru", "Руководство", "Другой встроенный текст"); err != nil {
		t.Fatalf("посев после правки: %v", err)
	}
	_, body, _, err := store.Get(ctx, "guide-edit-probe", "ru")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if body != "Текст редактора" {
		t.Errorf("посев затёр правку человека: %q", body)
	}
}
