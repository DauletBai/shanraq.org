package syndicate

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func testModule() *Module {
	return &Module{baseURL: "https://shanraq.org", log: zap.NewNop()}
}

func TestRenderRSS(t *testing.T) {
	m := testModule()
	when := time.Date(2026, 7, 13, 10, 0, 0, 0, time.UTC)
	out, err := m.renderRSS("ru", []feedEntry{
		{Slug: "ekonomika", Title: "Экономика 2026", Summary: "Разбор цифр", Lang: "ru", Modified: when},
	})
	if err != nil {
		t.Fatalf("renderRSS: %v", err)
	}
	s := string(out)
	for _, want := range []string{
		`<?xml version="1.0"`,
		`<rss version="2.0">`,
		"<title>Экономика 2026</title>",
		"https://shanraq.org/read/ekonomika?lang=ru",
		"<language>ru</language>",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("RSS missing %q\n---\n%s", want, s)
		}
	}
}

func TestBuildTelegramMessageEscapes(t *testing.T) {
	msg := buildTelegramMessage("Цены <выросли> & упали", "Кратко про <тэги>", "https://shanraq.org/read/x?lang=ru")
	if strings.Contains(msg, "<выросли>") {
		t.Errorf("title not HTML-escaped: %s", msg)
	}
	for _, want := range []string{"&lt;выросли&gt;", "<b>", "🔗", "https://shanraq.org/read/x?lang=ru"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q: %s", want, msg)
		}
	}
}

func TestTelegramDisabledByDefault(t *testing.T) {
	m := testModule()
	if m.TelegramEnabled() {
		t.Fatal("telegram should be disabled without config")
	}
	// EnqueuePublish must be a safe no-op when disabled (nil store tolerated).
	if err := m.EnqueuePublish(context.Background(), nil, uuid.New()); err != nil {
		t.Fatalf("EnqueuePublish no-op: %v", err)
	}
}

// TestFetchFeedIntegration checks the RSS query against a real DB (schema from
// migrations). Skipped unless SHANRAQ_TEST_DB is set.
func TestFetchFeedIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the RSS integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	authorID := uuid.New()
	articleID := uuid.New()
	slug := "rss-itest-" + articleID.String()[:8]
	_, _ = pool.Exec(ctx, `INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`, authorID, "rss-"+authorID.String()+"@t.test")
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, authorID) })
	if _, err := pool.Exec(ctx, `INSERT INTO articles (id, author_id, slug, original_lang, status, published_at) VALUES ($1,$2,$3,'ru','published',NOW())`, articleID, authorID, slug); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES ($1,'ru','РСС Тест','Аннотация','Тело','human','ready')`, articleID); err != nil {
		t.Fatalf("insert translation: %v", err)
	}

	m := &Module{db: pool, baseURL: "https://shanraq.org", log: zap.NewNop()}
	entries, err := m.fetchFeed(ctx, "ru", 30)
	if err != nil {
		t.Fatalf("fetchFeed: %v", err)
	}
	var found bool
	for _, e := range entries {
		if e.Slug == slug {
			found = true
			if e.Title != "РСС Тест" {
				t.Errorf("title = %q", e.Title)
			}
		}
	}
	if !found {
		t.Fatalf("published article %s not present in feed", slug)
	}
}
