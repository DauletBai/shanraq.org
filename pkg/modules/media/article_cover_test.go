package media

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestArticleCoverOnlyServesPublishedData(t *testing.T) {
	pool := ledgerPool(t)
	owner := seedOwner(t, pool)
	id := uuid.New()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO articles(id,author_id,slug,status)
		VALUES($1,$2,$3,'draft')`, id, owner, "cover-test-"+id.String()); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO article_cover_blobs(article_id,content_type,body)
		VALUES($1,'image/png',$2)`, id, []byte("cover bytes")); err != nil {
		t.Fatal(err)
	}
	m := &Module{ledger: &ledger{db: pool}}
	router := chi.NewRouter()
	router.Get("/article-covers/{article}", m.handleArticleCover)
	request := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}
	path := "/article-covers/" + id.String()
	if rec := request(path); rec.Code != http.StatusNotFound {
		t.Fatalf("draft returned %d", rec.Code)
	}
	if _, err := pool.Exec(ctx, `UPDATE articles SET status='published' WHERE id=$1`, id); err != nil {
		t.Fatal(err)
	}
	if rec := request(path); rec.Code != http.StatusOK || rec.Body.String() != "cover bytes" || rec.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("published cover: status=%d type=%q body=%q", rec.Code, rec.Header().Get("Content-Type"), rec.Body.String())
	}
	if _, err := pool.Exec(ctx, `UPDATE articles SET status='archived' WHERE id=$1`, id); err != nil {
		t.Fatal(err)
	}
	if rec := request(path); rec.Code != http.StatusNotFound {
		t.Fatalf("archived returned %d", rec.Code)
	}
}
