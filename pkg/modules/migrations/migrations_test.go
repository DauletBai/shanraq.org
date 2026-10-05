package migrations

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Production cleanup runs on every restart. A published column may have real
// views and votes by then, so its counters must survive the next startup.
func TestProductionCleanupPreservesArticleCounters(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the migration integration test")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	id := uuid.New()
	_, err = db.ExecContext(ctx, `INSERT INTO articles
		(id,author_id,slug,original_lang,status,published_at,views_count,score)
		VALUES ($1,'5a2a0000-0000-0000-0000-000000000001',$2,'ru','published',NOW(),5,3)`,
		id, "cleanup-"+id.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id=$1`, id)
		_ = db.Close()
	})
	if err := stripDemoFixtures(ctx, db); err != nil {
		t.Fatal(err)
	}
	var views, score int
	if err := db.QueryRowContext(ctx, `SELECT views_count,score FROM articles WHERE id=$1`, id).Scan(&views, &score); err != nil {
		t.Fatal(err)
	}
	if views != 5 || score != 3 {
		t.Fatalf("production restart changed article counters: views=%d score=%d", views, score)
	}
}
