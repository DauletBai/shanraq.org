package articles

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// favTypes is the closed set of things a user can bookmark.
var favTypes = map[string]bool{"article": true, "listing": true}

// FavoriteStore persists a user's bookmarks over articles and listings.
type FavoriteStore struct{ db *pgxpool.Pool }

// NewFavoriteStore builds the store.
func NewFavoriteStore(db *pgxpool.Pool) *FavoriteStore { return &FavoriteStore{db: db} }

// Toggle flips the favorite for (user, type, id) and reports whether it is now
// favorited. It deletes an existing row, or inserts one if none was present.
func (s *FavoriteStore) Toggle(ctx context.Context, userID uuid.UUID, itemType string, itemID uuid.UUID) (bool, error) {
	if !favTypes[itemType] {
		return false, fmt.Errorf("unknown favorite type %q", itemType)
	}
	tag, err := s.db.Exec(ctx,
		`DELETE FROM favorites WHERE user_id=$1 AND item_type=$2 AND item_id=$3`,
		userID, itemType, itemID)
	if err != nil {
		return false, fmt.Errorf("favorite remove: %w", err)
	}
	if tag.RowsAffected() > 0 {
		return false, nil
	}
	if _, err := s.db.Exec(ctx,
		`INSERT INTO favorites (user_id, item_type, item_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		userID, itemType, itemID); err != nil {
		return false, fmt.Errorf("favorite add: %w", err)
	}
	return true, nil
}

// IsFavorite reports whether the user has bookmarked the item. Errors are
// swallowed into false so a display path never fails on a bookmark check.
func (s *FavoriteStore) IsFavorite(ctx context.Context, userID uuid.UUID, itemType string, itemID uuid.UUID) bool {
	var exists bool
	if err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id=$1 AND item_type=$2 AND item_id=$3)`,
		userID, itemType, itemID).Scan(&exists); err != nil {
		return false
	}
	return exists
}

// ListFavorited returns the user's bookmarked published articles, newest saved
// first, hydrated with translations for rendering.
func (s *Store) ListFavorited(ctx context.Context, userID uuid.UUID) ([]*Article, error) {
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM favorites f
		JOIN articles a  ON a.id = f.item_id
		JOIN auth_users u ON u.id = a.author_id
		WHERE f.user_id = $1 AND f.item_type = 'article' AND a.status = 'published'
		ORDER BY f.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list favorited articles: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// ListFavorited returns the user's bookmarked published listings, newest saved first.
func (s *ListingStore) ListFavorited(ctx context.Context, userID uuid.UUID) ([]*Listing, error) {
	q := fmt.Sprintf(`SELECT %s FROM favorites f
		JOIN listings l   ON l.id = f.item_id
		JOIN auth_users u ON u.id = l.author_id
		WHERE f.user_id = $1 AND f.item_type = 'listing' AND l.status = 'published'
		ORDER BY f.created_at DESC`, listingCols)
	rows, err := s.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list favorited listings: %w", err)
	}
	defer rows.Close()
	var out []*Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
