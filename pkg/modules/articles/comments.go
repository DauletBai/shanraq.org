package articles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"shanraq.org/pkg/modules/auth"
)

const maxCommentLen = 2000

// commentCollapseScore is the score at or below which a comment is folded away.
//
// Folded, not removed: one line stands in its place and one click opens it.
// Being voted down is a statement about how a comment was received, and hiding
// it outright would make that statement into a verdict — which is what the
// model used to do here, and the reason it no longer does.
//
// Five is roughly five readers with no reputation of their own, or two who have
// earned some. Below that a comment has to have genuinely lost the room.
const commentCollapseScore = -5

// Comment is one reader comment on an article.
type Comment struct {
	ID         string
	AuthorName string
	Body       string
	CreatedAt  time.Time
	// Mine marks a comment written by the reader currently viewing the page.
	Mine bool

	// Score is the weighted sum of readers' votes, and UserVote is what the
	// reader looking at the page voted (-1, 0 or +1).
	Score    int
	UserVote int
}

// Collapsed reports whether readers have voted this comment far enough down to
// be folded away behind a single line.
func (c Comment) Collapsed() bool { return c.Score <= commentCollapseScore }

// CommentStore persists reader comments.
type CommentStore struct {
	db *pgxpool.Pool
}

func NewCommentStore(db *pgxpool.Pool) *CommentStore { return &CommentStore{db: db} }

// Create stores a published comment. The body is trimmed and length-capped.
func (s *CommentStore) Create(ctx context.Context, articleID, userID uuid.UUID, body string) error {
	return s.CreateWithStatus(ctx, articleID, userID, body, "published")
}

// CreateWithStatus stores a comment with an explicit moderation status.
// 'hidden' is what a human moderator sets when a comment has to come down; it
// is no longer set on the way in. A model used to make that call for every
// comment before anyone saw it — readers now make it by voting, and a comment
// they bury is folded rather than hidden.
func (s *CommentStore) CreateWithStatus(ctx context.Context, articleID, userID uuid.UUID, body, status string) error {
	body = strings.TrimSpace(body)
	if body == "" {
		return fmt.Errorf("empty comment")
	}
	if len(body) > maxCommentLen {
		body = body[:maxCommentLen]
	}
	if status != "hidden" {
		status = "published"
	}
	_, err := s.db.Exec(ctx, `INSERT INTO comments (article_id, user_id, body, status) VALUES ($1,$2,$3,$4)`,
		articleID, userID, body, status)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

// Delete removes a reader's own comment. Scoped to user_id in the statement, so
// someone else's id deletes nothing. A comment is published under the author's
// real name, which is precisely why they must be able to take it back without
// asking an administrator.
func (s *CommentStore) Delete(ctx context.Context, id, userID uuid.UUID) error {
	ct, err := s.db.Exec(ctx, `DELETE FROM comments WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ListForArticle returns published comments oldest first, with the author name.
// viewer marks which of them belong to the reader looking at the page, so the
// template can offer them a delete button; pass uuid.Nil for a guest.
func (s *CommentStore) ListForArticle(ctx context.Context, articleID, viewer uuid.UUID) ([]Comment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, u.email, u.first_name, u.last_name, u.middle_name, c.body, c.created_at,
		       c.user_id = $2 AS mine, c.score, COALESCE(v.value, 0)
		FROM comments c
		JOIN auth_users u ON u.id = c.user_id
		LEFT JOIN comment_votes v ON v.comment_id = c.id AND v.user_id = $2
		WHERE c.article_id = $1 AND c.status = 'published'
		ORDER BY c.created_at`, articleID, viewer)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()
	out := []Comment{}
	for rows.Next() {
		var c Comment
		var id uuid.UUID
		var email, first, last, middle string
		if err := rows.Scan(&id, &email, &first, &last, &middle, &c.Body, &c.CreatedAt,
			&c.Mine, &c.Score, &c.UserVote); err != nil {
			return nil, err
		}
		c.ID = id.String()
		// Comments are attributed as "Фамилия И.О." — formal and compact.
		c.AuthorName = auth.ShortName(first, last, middle, email)
		out = append(out, c)
	}
	return out, rows.Err()
}
