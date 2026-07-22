package articles

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"shanraq.org/internal/config"
	"shanraq.org/pkg/shanraq"
)

// TestDecideArticleIntegration exercises the human review decision against a
// real PostgreSQL: an article stuck in 'review' (the launch blocker) is cleared
// by a moderator, and the status change and the ledger entry must both land, in
// one transaction. Skipped unless SHANRAQ_TEST_DB points at a migrated database.
func TestDecideArticleIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the review-decision integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	author, admin, article := uuid.New(), uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO auth_users (id,email,password_hash,role) VALUES ($1,$2,'x','user'),($3,$4,'x','admin')`,
		author, "rev-a-"+author.String()+"@t.test", admin, "rev-m-"+admin.String()+"@t.test"); err != nil {
		t.Fatalf("users: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id IN ($1,$2)`, author, admin) })
	if _, err := pool.Exec(ctx, `INSERT INTO articles (id,author_id,slug,original_lang,category,status,submitted_at)
		VALUES ($1,$2,$3,'ru','economy','review',NOW())`, article, author, "decide-itest-"+article.String()[:8]); err != nil {
		t.Fatalf("article: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM articles WHERE id=$1`, article) })
	if _, err := pool.Exec(ctx, `INSERT INTO article_translations (article_id,lang,title,summary,body_md,source,status)
		VALUES ($1,'ru','Заголовок','саммари','тело','human','ready')`, article); err != nil {
		t.Fatalf("translation: %v", err)
	}

	m := &Module{
		rt:    &shanraq.Runtime{DB: pool, Logger: zap.NewNop(), Config: config.Config{}},
		store: NewStore(pool),
		mods:  NewModStore(pool),
	}

	// The queue must surface the stuck article.
	q, err := m.store.ReviewQueue(ctx, 100)
	if err != nil {
		t.Fatalf("queue: %v", err)
	}
	found := false
	for _, it := range q {
		if it.ID == article.String() {
			found = true
		}
	}
	if !found {
		t.Fatal("stuck article did not appear in the review queue")
	}

	// Approve it. Status and ledger must both reflect the decision.
	if err := m.DecideArticle(ctx, article, "approve", admin, "ok"); err != nil {
		t.Fatalf("decide: %v", err)
	}

	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM articles WHERE id=$1`, article).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != "published" {
		t.Fatalf("want published, got %q", status)
	}
	var logged int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM moderation_actions
		WHERE target_id=$1 AND action='approve' AND actor_kind='human'`, article.String()).Scan(&logged); err != nil {
		t.Fatalf("read ledger: %v", err)
	}
	if logged != 1 {
		t.Fatalf("want exactly one human approval in the ledger, got %d", logged)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM moderation_actions WHERE target_id=$1`, article.String()) })
}
