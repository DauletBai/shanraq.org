package ratings

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrSelfVote is returned when an author tries to vote on their own article.
var ErrSelfVote = errors.New("cannot vote on your own article")

type pgxPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

// Store persists votes and reputation.
type Store struct {
	db pgxPool
}

// NewStore builds a Store over the shared pgx pool.
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// Vote records, updates, or (when value==VoteNone) retracts a reader's vote on
// an article, then recomputes the article's cached score and the author's
// karma inside one transaction. It returns the article's new score.
func (s *Store) Vote(ctx context.Context, articleID, voterID, authorID uuid.UUID, value int) (int, error) {
	if voterID == authorID {
		return 0, ErrSelfVote
	}
	if value != VoteUp && value != VoteDown && value != VoteNone {
		return 0, fmt.Errorf("invalid vote value %d", value)
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if value == VoteNone {
		if _, err := tx.Exec(ctx, `DELETE FROM article_votes WHERE article_id = $1 AND user_id = $2`, articleID, voterID); err != nil {
			return 0, fmt.Errorf("delete vote: %w", err)
		}
	} else {
		var voterKarma int
		err := tx.QueryRow(ctx, `SELECT COALESCE(karma, 0) FROM author_reputation WHERE user_id = $1`, voterID).Scan(&voterKarma)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("voter karma: %w", err)
		}
		weight := Weight(voterKarma)
		if _, err := tx.Exec(ctx, `
			INSERT INTO article_votes (article_id, user_id, value, weight)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (article_id, user_id) DO UPDATE SET
				value = EXCLUDED.value,
				weight = EXCLUDED.weight,
				updated_at = NOW()
		`, articleID, voterID, value, weight); err != nil {
			return 0, fmt.Errorf("upsert vote: %w", err)
		}
	}

	// Recompute the article's cached score.
	var score int
	if err := tx.QueryRow(ctx, `
		UPDATE articles
		SET score = COALESCE((SELECT SUM(value * weight) FROM article_votes WHERE article_id = $1), 0),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING score
	`, articleID).Scan(&score); err != nil {
		return 0, fmt.Errorf("recompute score: %w", err)
	}

	// Recompute the author's karma across all of their articles.
	if err := recomputeKarma(ctx, tx, authorID); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return score, nil
}

func recomputeKarma(ctx context.Context, tx pgx.Tx, authorID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO author_reputation (user_id, karma, updated_at)
		VALUES ($1, COALESCE((
			SELECT SUM(av.value * av.weight)
			FROM article_votes av
			JOIN articles a ON a.id = av.article_id
			WHERE a.author_id = $1
		), 0), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			karma = EXCLUDED.karma,
			updated_at = NOW()
	`, authorID)
	if err != nil {
		return fmt.Errorf("recompute karma: %w", err)
	}
	return nil
}

// ForArticle returns an article's score plus the viewer's own vote (0 if the
// viewer is anonymous or has not voted).
func (s *Store) ForArticle(ctx context.Context, articleID, viewerID uuid.UUID) (Rating, error) {
	var r Rating
	err := s.db.QueryRow(ctx, `
		SELECT a.score, COALESCE(v.value, 0)
		FROM articles a
		LEFT JOIN article_votes v ON v.article_id = a.id AND v.user_id = $2
		WHERE a.id = $1
	`, articleID, viewerID).Scan(&r.Score, &r.UserVote)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Rating{}, nil
		}
		return Rating{}, fmt.Errorf("article rating: %w", err)
	}
	return r, nil
}

// AuthorKarma returns an author's accumulated reputation (0 if none yet).
func (s *Store) AuthorKarma(ctx context.Context, authorID uuid.UUID) (int, error) {
	var karma int
	err := s.db.QueryRow(ctx, `SELECT COALESCE(karma, 0) FROM author_reputation WHERE user_id = $1`, authorID).Scan(&karma)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("author karma: %w", err)
	}
	return karma, nil
}
