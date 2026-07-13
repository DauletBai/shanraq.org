package ai

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
)

// TestTranslateJobIntegration exercises the full translate-job DB path against a
// real PostgreSQL (schema created by the migrations module), using a fake
// completer so no API key or spend is involved. Skipped unless SHANRAQ_TEST_DB
// points at a database with the articles schema.
func TestTranslateJobIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the AI translate integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	// Fresh author + article with a Russian original.
	authorID := uuid.New()
	articleID := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`, authorID, "ai-itest-"+authorID.String()+"@t.test"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, authorID) })

	if _, err := pool.Exec(ctx, `INSERT INTO articles (id, author_id, slug, original_lang, status) VALUES ($1,$2,$3,'ru','draft')`, articleID, authorID, "ai-itest-"+articleID.String()[:8]); err != nil {
		t.Fatalf("insert article: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status) VALUES ($1,'ru','Заголовок','Кратко','Тело статьи','human','ready')`, articleID); err != nil {
		t.Fatalf("insert original: %v", err)
	}

	m := New()
	m.db = pool
	m.log = zap.NewNop()
	m.setCompleter(&fakeCompleter{reply: func(r Request) string { return "TRANSLATED:" + r.User }})

	payload, _ := EnqueuePayload(articleID)
	job := jobs.Job{ID: uuid.New(), Name: JobTranslate, Payload: payload}
	if err := m.handleTranslateJob(ctx, nil, job); err != nil {
		t.Fatalf("handleTranslateJob: %v", err)
	}

	// KZ and EN AI versions must now exist; RU original stays human & untouched.
	assertTranslation(t, ctx, pool, articleID, "kz", "ai", "TRANSLATED:Заголовок")
	assertTranslation(t, ctx, pool, articleID, "en", "ai", "TRANSLATED:Заголовок")
	assertTranslation(t, ctx, pool, articleID, "ru", "human", "Заголовок")

	// Now add a HUMAN Kazakh version and re-run: it must NOT be overwritten.
	if _, err := pool.Exec(ctx, `
		INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status)
		VALUES ($1,'kz','Қолмен','','Қолмен жазылған','human','ready')
		ON CONFLICT (article_id, lang) DO UPDATE SET title=EXCLUDED.title, body_md=EXCLUDED.body_md, source='human'`, articleID); err != nil {
		t.Fatalf("insert human kz: %v", err)
	}
	if err := m.handleTranslateJob(ctx, nil, job); err != nil {
		t.Fatalf("handleTranslateJob rerun: %v", err)
	}
	assertTranslation(t, ctx, pool, articleID, "kz", "human", "Қолмен")
}

func assertTranslation(t *testing.T, ctx context.Context, pool *pgxpool.Pool, articleID uuid.UUID, lang, wantSource, wantTitle string) {
	t.Helper()
	var source, title string
	err := pool.QueryRow(ctx, `SELECT source, title FROM article_translations WHERE article_id=$1 AND lang=$2`, articleID, lang).Scan(&source, &title)
	if err != nil {
		t.Fatalf("query %s: %v", lang, err)
	}
	if source != wantSource {
		t.Errorf("%s source = %q, want %q", lang, source, wantSource)
	}
	if title != wantTitle {
		t.Errorf("%s title = %q, want %q", lang, title, wantTitle)
	}
}
