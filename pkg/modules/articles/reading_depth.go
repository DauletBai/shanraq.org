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
// each depth milestone, split by how the article was taken in:
// map[articleID][mode][depth].
//
// The split is not decoration. Scrolling away at half is giving up; stopping a
// recording at half may be arriving at work. Summed together the two hide each
// other, so they are kept apart all the way to the screen and only added where
// a total is what is wanted.
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

// readFinishShare is how much of the estimated reading time a session must have
// spent engaged before it counts as a read. Half, because the estimate is a
// blunt 180 words a minute and real readers scatter widely around it -- a fast
// one who genuinely finished should not be thrown away to exclude a flick, and
// a flick takes seconds, nowhere near half.
const readFinishShare = 0.5

// readMaxSeconds caps one session's contribution. Engaged time already stops
// for a hidden tab and for idleness, but a page open all day in front of
// somebody doing something else would still dominate an average it has no
// business being in.
const readMaxSeconds = 3 * 60 * 60

// handleReadDone records a finished reading session: how far the reader got and
// how long they were actually engaged with the text.
//
// It arrives on page-hide, from navigator.sendBeacon, so it is fire-and-forget
// and always answers 204. Whether the session counts as a read is decided here
// rather than in the browser: the expected time is computed from the text the
// server itself holds, so a client cannot report having read something faster
// than the words in it allow.
func (m *Module) handleReadDone(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)

	secs, _ := strconv.Atoi(r.URL.Query().Get("t"))
	depth, _ := strconv.Atoi(r.URL.Query().Get("d"))
	if secs <= 0 {
		return
	}
	if secs > readMaxSeconds {
		secs = readMaxSeconds
	}
	// A crawler that runs JavaScript is rare; one that also idles on the page
	// for minutes is rarer still. Even so, the same filter the view counter uses
	// applies, so the two numbers are drawn from the same population.
	if botLabel(r.UserAgent()) != "" {
		return
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		return
	}
	// The page tags itself "kk" for Kazakh, which is the correct BCP 47 subtag
	// and not the code the translations are stored under. Read back what we
	// wrote rather than trusting the two to be the same string.
	lang := r.URL.Query().Get("l")
	if lang == "kk" {
		lang = LangKZ
	}
	expect := m.store.ReadingSeconds(r.Context(), a.ID, lang)
	finished := readCounts(depth, secs, expect)
	if err := m.store.RecordRead(r.Context(), a.ID, secs, finished); err != nil {
		m.rt.Logger.Warn("record read", zap.Error(err))
	}
}

// readCounts decides whether one session was a read.
//
// Both tests, not either: the end of the prose reached, and enough of the
// article's own estimate spent with it. Depth alone passes a two-second flick
// to the bottom; time alone passes a tab left open on the first paragraph.
func readCounts(depth, secs, expect int) bool {
	if depth < 100 || expect <= 0 {
		return false
	}
	return float64(secs) >= float64(expect)*readFinishShare
}
