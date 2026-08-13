package articles

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The prediction ledger. Every forecast the site makes goes in with the date it
// was made, and is later marked hit, miss or partial — publicly, the misses
// included.
//
// This is a credibility instrument, not a feature. Anyone can write a forecast;
// what nobody does is keep the receipt. A scoreboard that can go down is a
// costly signal — it would hurt to publish if the record were bad — and that is
// exactly why a reader believes it. It also earns return visits, because a
// resolved forecast is a reason to come back that no newsletter can manufacture.

// Prediction statuses. Partial exists because a forecast is usually neither
// wholly right nor wholly wrong, and forcing it into one of two boxes would
// make the score a fiction in whichever direction flattered us.
const (
	PredOpen    = "open"
	PredHit     = "hit"
	PredMiss    = "miss"
	PredPartial = "partial"
)

// PredStatuses lists the statuses in the order the admin form offers them.
var PredStatuses = []string{PredOpen, PredHit, PredMiss, PredPartial}

func validPredStatus(s string) bool {
	for _, v := range PredStatuses {
		if v == s {
			return true
		}
	}
	return false
}

// Prediction is one forecast with its per-language text already loaded.
type Prediction struct {
	ID         uuid.UUID
	ArticleID  *uuid.UUID
	MadeOn     time.Time
	Horizon    *time.Time
	Status     string
	ResolvedOn *time.Time
	SourceURL  string

	// Statement and Verdict keyed by language code.
	Statement map[string]string
	Verdict   map[string]string

	// ArticleSlug and ArticleTitle are filled for display when the forecast is
	// tied to a piece; empty when it is not, or when the article was deleted.
	ArticleSlug  string
	ArticleTitle string
}

// StatementIn returns the forecast in the reader's language, falling back to
// any language that has text. A forecast shown as an empty line would be worse
// than one shown in the wrong language.
func (p *Prediction) StatementIn(lang string) string { return pickLang(p.Statement, lang) }

// VerdictIn returns what actually happened, in the reader's language.
func (p *Prediction) VerdictIn(lang string) string { return pickLang(p.Verdict, lang) }

func pickLang(m map[string]string, lang string) string {
	if v := strings.TrimSpace(m[lang]); v != "" {
		return v
	}
	for _, l := range Langs {
		if v := strings.TrimSpace(m[l]); v != "" {
			return v
		}
	}
	return ""
}

// Resolved reports whether the forecast has been judged.
func (p *Prediction) Resolved() bool { return p.Status != PredOpen }

// Overdue reports whether an open forecast has passed its own deadline. Shown
// in the admin so a prediction cannot quietly stay "open" forever, which is the
// obvious way a ledger like this gets gamed.
func (p *Prediction) Overdue() bool {
	return p.Status == PredOpen && p.Horizon != nil && p.Horizon.Before(time.Now())
}

// PredictionScore is the public tally.
type PredictionScore struct {
	Hit, Miss, Partial, Open int
}

// Resolved is the number of forecasts that have been judged.
func (s PredictionScore) Resolved() int { return s.Hit + s.Miss + s.Partial }

// Total is every forecast on the ledger.
func (s PredictionScore) Total() int { return s.Resolved() + s.Open }

// Accuracy is the share of judged forecasts that came true outright, rounded to
// a whole percent. Partial hits count as half: calling them successes would
// flatter the number and calling them failures would punish honesty about the
// middle ground, and the middle ground is where most forecasts land.
func (s PredictionScore) Accuracy() int {
	n := s.Resolved()
	if n == 0 {
		return 0
	}
	return int((float64(s.Hit) + 0.5*float64(s.Partial)) / float64(n) * 100.0)
}

// PredictionStore persists the ledger.
type PredictionStore struct{ db *pgxpool.Pool }

// NewPredictionStore builds a PredictionStore over the shared pgx pool.
func NewPredictionStore(db *pgxpool.Pool) *PredictionStore { return &PredictionStore{db: db} }

const predSelect = `
	SELECT p.id, p.article_id, p.made_on, p.horizon, p.status, p.resolved_on, p.source_url,
	       COALESCE(a.slug, ''), COALESCE(t.title, '')
	  FROM predictions p
	  LEFT JOIN articles a ON a.id = p.article_id
	  LEFT JOIN article_translations t ON t.article_id = a.id AND t.lang = $1`

// List returns the whole ledger, open forecasts first (soonest deadline first,
// because those are the ones a reader can still watch), then the settled ones
// newest first.
func (s *PredictionStore) List(ctx context.Context, lang string) ([]*Prediction, error) {
	rows, err := s.db.Query(ctx, predSelect+`
		 ORDER BY (p.status <> 'open'), p.horizon NULLS LAST, p.made_on DESC
		 LIMIT 1000`, lang)
	if err != nil {
		return nil, err
	}
	out, err := scanPredictions(rows)
	if err != nil {
		return nil, err
	}
	return out, s.loadTexts(ctx, out)
}

// ForArticle returns the forecasts made in one article, oldest first, so the
// block under the piece reads in the order they were written.
func (s *PredictionStore) ForArticle(ctx context.Context, lang string, id uuid.UUID) ([]*Prediction, error) {
	rows, err := s.db.Query(ctx, predSelect+`
		 WHERE p.article_id = $2 ORDER BY p.made_on, p.id LIMIT 100`, lang, id)
	if err != nil {
		return nil, err
	}
	out, err := scanPredictions(rows)
	if err != nil {
		return nil, err
	}
	return out, s.loadTexts(ctx, out)
}

// Get returns one forecast for the admin editor.
func (s *PredictionStore) Get(ctx context.Context, lang string, id uuid.UUID) (*Prediction, error) {
	rows, err := s.db.Query(ctx, predSelect+` WHERE p.id = $2`, lang, id)
	if err != nil {
		return nil, err
	}
	out, err := scanPredictions(rows)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, pgx.ErrNoRows
	}
	return out[0], s.loadTexts(ctx, out)
}

func scanPredictions(rows pgx.Rows) ([]*Prediction, error) {
	defer rows.Close()
	out := []*Prediction{}
	for rows.Next() {
		p := &Prediction{Statement: map[string]string{}, Verdict: map[string]string{}}
		if err := rows.Scan(&p.ID, &p.ArticleID, &p.MadeOn, &p.Horizon, &p.Status,
			&p.ResolvedOn, &p.SourceURL, &p.ArticleSlug, &p.ArticleTitle); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// loadTexts fills every language of every prediction in one query, so a ledger
// of a hundred forecasts costs two round trips rather than a hundred and one.
func (s *PredictionStore) loadTexts(ctx context.Context, list []*Prediction) error {
	if len(list) == 0 {
		return nil
	}
	byID := make(map[uuid.UUID]*Prediction, len(list))
	ids := make([]uuid.UUID, 0, len(list))
	for _, p := range list {
		byID[p.ID] = p
		ids = append(ids, p.ID)
	}
	rows, err := s.db.Query(ctx,
		`SELECT prediction_id, lang, statement, verdict FROM prediction_texts
		  WHERE prediction_id = ANY($1)`, ids)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var lang, statement, verdict string
		if err := rows.Scan(&id, &lang, &statement, &verdict); err != nil {
			return err
		}
		if p := byID[id]; p != nil {
			p.Statement[lang] = statement
			p.Verdict[lang] = verdict
		}
	}
	return rows.Err()
}

// Score tallies the ledger in one query.
func (s *PredictionStore) Score(ctx context.Context) (PredictionScore, error) {
	var sc PredictionScore
	rows, err := s.db.Query(ctx, `SELECT status, COUNT(*) FROM predictions GROUP BY status`)
	if err != nil {
		return sc, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			return sc, err
		}
		switch status {
		case PredHit:
			sc.Hit = n
		case PredMiss:
			sc.Miss = n
		case PredPartial:
			sc.Partial = n
		case PredOpen:
			sc.Open = n
		}
	}
	return sc, rows.Err()
}

// PredictionInput is one forecast as the admin form submits it.
type PredictionInput struct {
	ArticleID  *uuid.UUID
	MadeOn     time.Time
	Horizon    *time.Time
	Status     string
	ResolvedOn *time.Time
	SourceURL  string
	Statement  map[string]string
	Verdict    map[string]string
}

// ErrPredictionEmpty is returned when a forecast has no text in any language.
var ErrPredictionEmpty = errors.New("a prediction needs a statement in at least one language")

// Validate normalizes the input and rejects the states the ledger must not
// hold. A resolved forecast is stamped with today's date if the operator did
// not supply one, so the database CHECK can never be the thing that reports a
// missing date to a person filling in a form.
func (in *PredictionInput) Validate() error {
	if !validPredStatus(in.Status) {
		return fmt.Errorf("unknown prediction status %q", in.Status)
	}
	if in.MadeOn.IsZero() {
		in.MadeOn = time.Now()
	}
	empty := true
	for _, l := range Langs {
		in.Statement[l] = strings.TrimSpace(in.Statement[l])
		in.Verdict[l] = strings.TrimSpace(in.Verdict[l])
		if in.Statement[l] != "" {
			empty = false
		}
	}
	if empty {
		return ErrPredictionEmpty
	}
	if in.Status == PredOpen {
		// Reopening a forecast must clear its verdict date, or the row violates
		// the table's own invariant.
		in.ResolvedOn = nil
	} else if in.ResolvedOn == nil {
		now := time.Now()
		in.ResolvedOn = &now
	}
	return nil
}

// Save writes a forecast and all its languages in one transaction: a partial
// write would leave a prediction whose statement exists in Russian and whose
// verdict does not, which is precisely the kind of gap that makes a ledger
// look edited.
func (s *PredictionStore) Save(ctx context.Context, id uuid.UUID, in PredictionInput) (uuid.UUID, error) {
	if err := in.Validate(); err != nil {
		return uuid.Nil, err
	}
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin predictions tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if id == uuid.Nil {
		err = tx.QueryRow(ctx,
			`INSERT INTO predictions (article_id, made_on, horizon, status, resolved_on, source_url)
			 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
			in.ArticleID, in.MadeOn, in.Horizon, in.Status, in.ResolvedOn, in.SourceURL).Scan(&id)
	} else {
		_, err = tx.Exec(ctx,
			`UPDATE predictions SET article_id=$2, made_on=$3, horizon=$4, status=$5,
			        resolved_on=$6, source_url=$7, updated_at=NOW() WHERE id=$1`,
			id, in.ArticleID, in.MadeOn, in.Horizon, in.Status, in.ResolvedOn, in.SourceURL)
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("write prediction: %w", err)
	}
	for _, l := range Langs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO prediction_texts (prediction_id, lang, statement, verdict)
			 VALUES ($1,$2,$3,$4)
			 ON CONFLICT (prediction_id, lang) DO UPDATE
			   SET statement = EXCLUDED.statement, verdict = EXCLUDED.verdict`,
			id, l, in.Statement[l], in.Verdict[l]); err != nil {
			return uuid.Nil, fmt.Errorf("write prediction text %s: %w", l, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit prediction: %w", err)
	}
	return id, nil
}

// Delete removes a forecast. Reserved for a duplicate or a typo: a forecast
// deleted because it turned out badly would empty the ledger of its meaning,
// which is why the admin marks the action as such.
func (s *PredictionStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `DELETE FROM predictions WHERE id = $1`, id)
	return err
}

// PredictionYears groups the ledger by the year each forecast was made, newest
// year first, for the public page's headings.
func PredictionYears(list []*Prediction) []int {
	seen := map[int]bool{}
	out := []int{}
	for _, p := range list {
		if y := p.MadeOn.Year(); !seen[y] {
			seen[y] = true
			out = append(out, y)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(out)))
	return out
}
