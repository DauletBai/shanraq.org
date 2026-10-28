package media

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// handleArticleCover only exposes covers of published articles. The article ID
// is opaque, but knowing it is not sufficient to reveal a draft's image.
func (m *Module) handleArticleCover(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "article"))
	if err != nil || m.ledger == nil {
		http.NotFound(w, r)
		return
	}
	var data []byte
	var contentType string
	err = m.ledger.db.QueryRow(r.Context(), `SELECT b.body, b.content_type
		FROM article_cover_blobs b JOIN articles a ON a.id=b.article_id
		WHERE b.article_id=$1 AND a.status='published'`, id).Scan(&data, &contentType)
	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "cover unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, "cover", time.Time{}, bytes.NewReader(data))
}
