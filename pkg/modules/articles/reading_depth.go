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
	// How the article was taken in. Anything unrecognised is reading: the
	// parameter arrived later than the beacon that sends it, and an old cached
	// page must keep counting rather than start dropping its reports.
	mode := ModeRead
	if r.URL.Query().Get("m") == ModeListen {
		mode = ModeListen
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := m.store.RecordDepth(r.Context(), a.ID, mode, d); err != nil {
		m.rt.Logger.Warn("record reading depth", zap.Error(err))
	}
	w.WriteHeader(http.StatusNoContent)
}

// validDepth is the closed set of reading-depth milestones (percent).
var validDepth = map[int]bool{25: true, 50: true, 75: true, 100: true}

// How a reader took the article in. The two are counted apart because they do
// not mean the same thing: someone who scrolls away halfway has stopped, while
// someone listening halfway may be driving and will hear it out.
const (
	ModeRead   = "read"
	ModeListen = "listen"
)

// pctOf returns v as a whole-percent share of total (0 when total is 0).
func pctOf(v, total int64) int {
	if total <= 0 {
		return 0
	}
	return int(v * 100 / total)
}

// RecordDepth increments the counter for how many readers reached a given depth
// milestone on an article.
func (s *Store) RecordDepth(ctx context.Context, articleID uuid.UUID, mode string, depth int) error {
	if !validDepth[depth] {
		return fmt.Errorf("invalid depth %d", depth)
	}
	if mode != ModeRead && mode != ModeListen {
		return fmt.Errorf("invalid mode %q", mode)
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO reading_depth (article_id, mode, depth, count) VALUES ($1,$2,$3,1)
		ON CONFLICT (article_id, mode, depth) DO UPDATE SET count = reading_depth.count + 1`,
		articleID, mode, depth)
	return err
}

// AuthorReadingDepth returns, per article of the author, the reader counts at
// each depth milestone, split by how the article was taken in:
// map[articleID][mode][depth].
//
// The split is not decoration. Scrolling away at half is giving up; stopping a
// recording at half may be arriving at work. Summed together the two hide each
// other, so they are kept apart all the way to the screen and only added where
// a total is what is wanted.
func (s *Store) AuthorReadingDepth(ctx context.Context, authorID uuid.UUID) (map[string]map[string]map[int]int64, error) {
	rows, err := s.db.Query(ctx, `
		SELECT rd.article_id, rd.mode, rd.depth, rd.count
		FROM reading_depth rd
		JOIN articles a ON a.id = rd.article_id
		WHERE a.author_id = $1`, authorID)
	if err != nil {
		return nil, fmt.Errorf("author reading depth: %w", err)
	}
	defer rows.Close()
	out := map[string]map[string]map[int]int64{}
	for rows.Next() {
		var id uuid.UUID
		var mode string
		var depth int
		var count int64
		if err := rows.Scan(&id, &mode, &depth, &count); err != nil {
			return nil, err
		}
		k := id.String()
		if out[k] == nil {
			out[k] = map[string]map[int]int64{}
		}
		if out[k][mode] == nil {
			out[k][mode] = map[int]int64{}
		}
		out[k][mode][depth] = count
	}
	return out, rows.Err()
}

// depthTotal adds the milestone across every mode: how many reached this far,
// whichever way they got there.
func depthTotal(byMode map[string]map[int]int64, depth int) int64 {
	var n int64
	for _, d := range byMode {
		n += d[depth]
	}
	return n
}
