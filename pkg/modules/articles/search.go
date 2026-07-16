package articles

import (
	"context"
	"fmt"
	"strings"
)

// searchLimit caps how many results a query returns.
const searchLimit = 40

// Search returns published articles matching the full-text query across all
// their language translations (title/summary/body), newest first.
func (s *Store) Search(ctx context.Context, query string) ([]*Article, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.status = 'published'
		  AND a.id IN (
		    SELECT article_id FROM article_translations
		    WHERE search_vector @@ websearch_to_tsquery('simple', $1)
		  )
		ORDER BY a.published_at DESC NULLS LAST
		LIMIT $2`, query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("search articles: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// Search returns active published listings matching the full-text query over
// title/description/location, promoted first then newest.
func (s *ListingStore) Search(ctx context.Context, query string) ([]*Listing, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	q := fmt.Sprintf(`SELECT %s FROM listings l JOIN auth_users u ON u.id = l.author_id
		WHERE l.status = 'published' AND l.expires_at > NOW()
		  AND l.search_vector @@ websearch_to_tsquery('simple', $1)
		ORDER BY COALESCE(l.promoted_until > NOW(), false) DESC, l.created_at DESC
		LIMIT $2`, listingCols)
	rows, err := s.db.Query(ctx, q, query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("search listings: %w", err)
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
