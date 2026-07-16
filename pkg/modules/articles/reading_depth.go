package articles

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// handleReadProgress records a reader reaching a depth milestone (25/50/75/100%)
// on an article. Called via navigator.sendBeacon, so it is fire-and-forget and
// always answers 204 — even for unknown slugs or bad input.
func (m *Module) handleReadProgress(w http.ResponseWriter, r *http.Request) {
	d, _ := strconv.Atoi(r.URL.Query().Get("d"))
	if !validDepth[d] {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := m.store.RecordDepth(r.Context(), a.ID, d); err != nil {
		m.rt.Logger.Warn("record reading depth", zap.Error(err))
	}
	w.WriteHeader(http.StatusNoContent)
}

// validDepth is the closed set of reading-depth milestones (percent).
var validDepth = map[int]bool{25: true, 50: true, 75: true, 100: true}

// pctOf returns v as a whole-percent share of total (0 when total is 0).
func pctOf(v, total int64) int {
	if total <= 0 {
		return 0
	}
	return int(v * 100 / total)
}

// RecordDepth increments the counter for how many readers reached a given depth
// milestone on an article.
func (s *Store) RecordDepth(ctx context.Context, articleID uuid.UUID, depth int) error {
	if !validDepth[depth] {
		return fmt.Errorf("invalid depth %d", depth)
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO reading_depth (article_id, depth, count) VALUES ($1,$2,1)
		ON CONFLICT (article_id, depth) DO UPDATE SET count = reading_depth.count + 1`,
		articleID, depth)
	return err
}

// AuthorReadingDepth returns, per article of the author, the reader counts at
// each depth milestone: map[articleID]{25:.., 50:.., 75:.., 100:..}.
func (s *Store) AuthorReadingDepth(ctx context.Context, authorID uuid.UUID) (map[string]map[int]int64, error) {
	rows, err := s.db.Query(ctx, `
		SELECT rd.article_id, rd.depth, rd.count
		FROM reading_depth rd
		JOIN articles a ON a.id = rd.article_id
		WHERE a.author_id = $1`, authorID)
	if err != nil {
		return nil, fmt.Errorf("author reading depth: %w", err)
	}
	defer rows.Close()
	out := map[string]map[int]int64{}
	for rows.Next() {
		var id uuid.UUID
		var depth int
		var count int64
		if err := rows.Scan(&id, &depth, &count); err != nil {
			return nil, err
		}
		k := id.String()
		if out[k] == nil {
			out[k] = map[int]int64{}
		}
		out[k][depth] = count
	}
	return out, rows.Err()
}
