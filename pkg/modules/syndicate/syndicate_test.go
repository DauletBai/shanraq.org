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

func TestRenderDigest(t *testing.T) {
	m := testModule()
	when := time.Date(2026, 7, 13, 10, 0, 0, 0, time.UTC)
	subject, body := m.renderDigest("ru", []feedEntry{
		{Slug: "ekonomika", Title: "Экономика 2026", Lang: "ru", Modified: when},
	}, "tok123")
	if subject != "Shanraq: обзор недели" {
		t.Errorf("subject = %q", subject)
	}
	for _, want := range []string{"Экономика 2026", "https://shanraq.org/read/ekonomika?lang=ru", "Отписаться", "/unsubscribe?token=tok123"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q:\n%s", want, body)
		}
	}
}

type fakeMailer struct{ sent []string }

func (f *fakeMailer) Send(_ context.Context, to, subject, body string) error {
	f.sent = append(f.sent, to+"|"+subject)
	return nil
}

// TestDigestIntegration exercises subscribe → SendDigest → unsubscribe against a
// real DB with a fake mailer. Skipped unless SHANRAQ_TEST_DB is set.
func TestDigestIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the digest integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	// A published article within the last 7 days, and a subscriber.
	authorID := uuid.New()
	articleID := uuid.New()
	email := "digest-" + articleID.String()[:8] + "@t.test"
	_, _ = pool.Exec(ctx, `INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`, authorID, "dg-"+authorID.String()+"@t.test")
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, authorID)
		_, _ = pool.Exec(ctx, `DELETE FROM subscribers WHERE email=$1`, email)
	})
	if _, err := pool.Exec(ctx, `INSERT INTO articles (id, author_id, slug, original_lang, status, published_at) VALUES ($1,$2,$3,'ru','published',NOW())`, articleID, authorID, "dg-"+articleID.String()[:8]); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES ($1,'ru','Дайджест тест','Аннотация','Тело','human','ready')`, articleID); err != nil {
		t.Fatalf("insert translation: %v", err)
	}

	fm := &fakeMailer{}
	m := &Module{db: pool, baseURL: "https://shanraq.org", log: zap.NewNop(), mailer: fm, emailEnabled: true}

	if err := m.subscribe(ctx, email, "ru"); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	sent, err := m.SendDigest(ctx)
	if err != nil {
		t.Fatalf("SendDigest: %v", err)
	}
	if sent < 1 {
		t.Fatalf("expected ≥1 sent, got %d", sent)
	}
	var found bool
	for _, s := range fm.sent {
		if strings.HasPrefix(s, email+"|") {
			found = true
		}
	}
	if !found {
		t.Fatalf("digest not sent to %s (sent: %v)", email, fm.sent)
	}

	// Unsubscribe by token removes the subscriber.
	var token string
	if err := pool.QueryRow(ctx, `SELECT unsubscribe_token FROM subscribers WHERE email=$1`, email).Scan(&token); err != nil {
		t.Fatalf("read token: %v", err)
	}
	if _, err := m.unsubscribe(ctx, token); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}
	var cnt int
	_ = pool.QueryRow(ctx, `SELECT COUNT(*) FROM subscribers WHERE email=$1`, email).Scan(&cnt)
	if cnt != 0 {
		t.Fatalf("subscriber not removed after unsubscribe")
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
