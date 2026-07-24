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

// Create stores a published comment. The body is trimmed and length-capped.
func (s *CommentStore) Create(ctx context.Context, articleID, userID uuid.UUID, body string) error {
	return s.CreateWithStatus(ctx, articleID, userID, body, "published")
}

// CreateWithStatus stores a comment with an explicit moderation status. The AI
// moderator uses this to file a suspect comment straight into 'hidden', the same
// queue human reports feed, so a human can confirm or restore it.
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

// ListForArticle returns published comments oldest first, with the author name.
func (s *CommentStore) ListForArticle(ctx context.Context, articleID uuid.UUID) ([]Comment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, u.email, u.first_name, u.last_name, u.middle_name, c.body, c.created_at
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
		var email, first, last, middle string
		if err := rows.Scan(&id, &email, &first, &last, &middle, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = id.String()
		// Comments are attributed as "Фамилия И.О." — formal and compact.
		c.AuthorName = auth.ShortName(first, last, middle, email)
		out = append(out, c)
	}
	return out, rows.Err()
}
