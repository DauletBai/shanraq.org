package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestExportPublishedDatabaseCover(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the export integration test")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatal("refusing to use a non-test database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	owner, id := uuid.New(), uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO auth_users(id,email,password_hash,role) VALUES($1,$2,'x','user')`, owner, owner.String()+"@test.invalid"); err != nil {
		t.Fatal(err)
	}
	defer pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, owner)
	if _, err := pool.Exec(ctx, `INSERT INTO articles(id,author_id,slug,status) VALUES($1,$2,$3,'published')`, id, owner, "export-cover-"+id.String()); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_cover_blobs(article_id,content_type,body) VALUES($1,'image/png',$2)`, id, []byte("image data")); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	count, err := copyDatabaseCovers(ctx, pool, []article{{Slug: "test", CoverURL: "/article-covers/" + id.String()}}, dir)
	if err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "article-covers", id.String()))
	if err != nil || string(b) != "image data" {
		t.Fatalf("cover=%q err=%v", b, err)
	}
	if _, err := copyDatabaseCovers(ctx, pool, []article{{Slug: "test", CoverURL: "/article-covers/../../escape"}}, dir); err == nil {
		t.Fatal("accepted invalid cover path")
	}
}
