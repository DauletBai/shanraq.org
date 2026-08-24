package articles

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CorrectionStore keeps reader-reported typos and their outcomes.
type CorrectionStore struct{ db *pgxpool.Pool }

// NewCorrectionStore builds the store over the shared pool.
func NewCorrectionStore(db *pgxpool.Pool) *CorrectionStore { return &CorrectionStore{db: db} }

// Insert records a claim as pending and returns its id. The row is written
// before anything is decided, so a claim is never lost to a failing checker or a
// dropped connection — an undecided correction can be looked at, a forgotten one
// cannot.
func (s *CorrectionStore) Insert(ctx context.Context, c Correction) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		INSERT INTO article_corrections
		    (article_id, lang, reporter_id, chapter, sentence, word)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id`,
		c.ArticleID, c.Lang, c.Reporter, c.Chapter, c.Sentence, c.Word).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert correction: %w", err)
	}
	return id, nil
}

// Decide records the outcome.
func (s *CorrectionStore) Decide(ctx context.Context, id uuid.UUID, status, fixed, reason string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE article_corrections
		SET status = $2, fixed = $3, reason = $4, decided_at = NOW()
		WHERE id = $1`, id, status, fixed, reason)
	if err != nil {
		return fmt.Errorf("decide correction: %w", err)
	}
	return nil
}

// CountSince counts one reporter's claims in a window, for the rate limit.
func (s *CorrectionStore) CountSince(ctx context.Context, reporter uuid.UUID, since time.Time) (int, error) {
	var n int
	err := s.db.QueryRow(ctx, `
		SELECT count(*) FROM article_corrections
		WHERE reporter_id = $1 AND created_at >= $2`, reporter, since).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count corrections: %w", err)
	}
	return n, nil
}

// Duplicate reports whether this reader already sent this exact claim for this
// article. Re-reading a piece and reporting the same word twice is ordinary, and
// it should not cost them a second slot in the day's allowance or the author a
// second notification.
func (s *CorrectionStore) Duplicate(ctx context.Context, articleID, reporter uuid.UUID, lang, word string) (string, bool, error) {
	var status string
	err := s.db.QueryRow(ctx, `
		SELECT status FROM article_corrections
		WHERE article_id = $1 AND reporter_id = $2 AND lang = $3 AND lower(word) = lower($4)
		ORDER BY created_at DESC LIMIT 1`,
		articleID, reporter, lang, word).Scan(&status)
	if err != nil {
		return "", false, nil
	}
	return status, true, nil
}

// ApplyToBody writes the corrected text back, and only if the text is still
// exactly what we read a moment ago.
//
// The comparison is the whole point: between reading the body and writing it back
// the author may have edited the piece, or a second correction may have landed.
// Writing blind would silently undo their work, so a stale write fails and the
// reader is told to try again rather than being told it worked.
func (s *CorrectionStore) ApplyToBody(ctx context.Context, articleID uuid.UUID, lang, was, now string) (bool, error) {
	tag, err := s.db.Exec(ctx, `
		UPDATE article_translations
		SET body_md = $4
		WHERE article_id = $1 AND lang = $2 AND body_md = $3`,
		articleID, lang, was, now)
	if err != nil {
		return false, fmt.Errorf("apply correction: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return false, nil
	}
	// The article's own timestamp moves too: the piece changed, and the sitemap
	// and the dateModified in its structured data are built from this.
	if _, err := s.db.Exec(ctx,
		`UPDATE articles SET updated_at = NOW() WHERE id = $1`, articleID); err != nil {
		return true, fmt.Errorf("touch article: %w", err)
	}
	return true, nil
}

// ForArticle lists the decided corrections on an article, newest first. This is
// what makes an applied edit answerable: the author can see every change a reader
// caused in their text, and what was refused.
func (s *CorrectionStore) ForArticle(ctx context.Context, articleID uuid.UUID, limit int) ([]Correction, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, lang, chapter, sentence, word, status, fixed, reason, created_at
		FROM article_corrections
		WHERE article_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, articleID, limit)
	if err != nil {
		return nil, fmt.Errorf("list corrections: %w", err)
	}
	defer rows.Close()
	out := []Correction{}
	for rows.Next() {
		var c Correction
		if err := rows.Scan(&c.ID, &c.Lang, &c.Chapter, &c.Sentence, &c.Word,
			&c.Status, &c.Fixed, &c.Reason, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ArticleID = articleID
		out = append(out, c)
	}
	return out, rows.Err()
}
