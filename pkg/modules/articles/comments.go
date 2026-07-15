package articles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxCommentLen = 2000

// Comment is one reader comment on an article.
type Comment struct {
	ID         string
	AuthorName string
	Body       string
	CreatedAt  time.Time
}

// CommentStore persists reader comments.
type CommentStore struct {
	db *pgxpool.Pool
}

func NewCommentStore(db *pgxpool.Pool) *CommentStore { return &CommentStore{db: db} }

// Create stores a comment. The body is trimmed and length-capped.
func (s *CommentStore) Create(ctx context.Context, articleID, userID uuid.UUID, body string) error {
	body = strings.TrimSpace(body)
	if body == "" {
		return fmt.Errorf("empty comment")
	}
	if len(body) > maxCommentLen {
		body = body[:maxCommentLen]
	}
	_, err := s.db.Exec(ctx, `INSERT INTO comments (article_id, user_id, body) VALUES ($1,$2,$3)`,
		articleID, userID, body)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

// ListForArticle returns published comments oldest first, with the author name.
func (s *CommentStore) ListForArticle(ctx context.Context, articleID uuid.UUID) ([]Comment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, u.email, c.body, c.created_at
		FROM comments c JOIN auth_users u ON u.id = c.user_id
		WHERE c.article_id = $1 AND c.status = 'published'
		ORDER BY c.created_at`, articleID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()
	out := []Comment{}
	for rows.Next() {
		var c Comment
		var id uuid.UUID
		var email string
		if err := rows.Scan(&id, &email, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = id.String()
		c.AuthorName = displayName(email)
		out = append(out, c)
	}
	return out, rows.Err()
}
