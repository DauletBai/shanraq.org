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

// Правка документа в коде должна доходить до сайта. Страницы засеиваются в базу
// на первом запуске, и при вставке «только если нет» дальнейшие правки в коде
// на сайт не попадали никогда: посетитель годами читал текст первого запуска.
func TestCorrectingADocumentInCodeReachesTheSite(t *testing.T) {
	app := newTestApp(t)
	defer app.cleanup()
	ctx := context.Background()
	store := NewContentStore(app.pool)
	// Тестовая база живёт между прогонами: строка, оставшаяся с прошлого раза,
	// прошла бы этот тест и без исправления.
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

// А правку человека посев не трогает: иначе перезапуск стирал бы работу
// редактора, и админка превратилась бы в ловушку.
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
