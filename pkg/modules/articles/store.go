package articles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indicates the requested article does not exist.
var ErrNotFound = errors.New("article not found")

type pgxPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

// Store manages persistence for articles and their translations.
type Store struct {
	db pgxPool
}

// NewStore builds a Store over the shared pgx pool.
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// TranslationInput carries one language version submitted from the editor.
type TranslationInput struct {
	Lang    string
	Title   string
	Summary string
	BodyMD  string
	Source  string
}

// Create inserts a new draft article plus any non-empty translations and
// returns the new article ID.
func (s *Store) Create(ctx context.Context, authorID uuid.UUID, slug, originalLang, category, subcategory, coverURL string, trs []TranslationInput) (uuid.UUID, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var id uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO articles (author_id, slug, original_lang, category, subcategory, cover_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'draft')
		RETURNING id
	`, authorID, slug, originalLang, category, subcategory, coverURL).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert article: %w", err)
	}

	for _, tr := range trs {
		if err := upsertTranslation(ctx, tx, id, tr); err != nil {
			return uuid.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit: %w", err)
	}
	return id, nil
}

// Update rewrites an article's translations (author-scoped).
func (s *Store) Update(ctx context.Context, id, authorID uuid.UUID, originalLang, category, subcategory, coverURL string, trs []TranslationInput) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	tag, err := tx.Exec(ctx, `
		UPDATE articles SET original_lang = $3, category = $4, subcategory = $5, cover_url = $6, updated_at = NOW()
		WHERE id = $1 AND author_id = $2
	`, id, authorID, originalLang, category, subcategory, coverURL)
	if err != nil {
		return fmt.Errorf("update article: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	for _, tr := range trs {
		if err := upsertTranslation(ctx, tx, id, tr); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func upsertTranslation(ctx context.Context, tx pgx.Tx, articleID uuid.UUID, tr TranslationInput) error {
	if tr.Title == "" && tr.BodyMD == "" && tr.Summary == "" {
		return nil
	}
	source := tr.Source
	if source == "" {
		source = "human"
	}
	status := "draft"
	if tr.Title != "" && tr.BodyMD != "" {
		status = "ready"
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO article_translations (article_id, lang, title, summary, body_md, source, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (article_id, lang) DO UPDATE SET
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			body_md = EXCLUDED.body_md,
			source = EXCLUDED.source,
			status = EXCLUDED.status,
			updated_at = NOW()
	`, articleID, tr.Lang, tr.Title, tr.Summary, tr.BodyMD, source, status)
	if err != nil {
		return fmt.Errorf("upsert translation %s: %w", tr.Lang, err)
	}
	return nil
}

// SetStatus transitions an article's lifecycle state (author-scoped).
func (s *Store) SetStatus(ctx context.Context, id, authorID uuid.UUID, status string) error {
	var publishedAt any
	if status == "published" {
		publishedAt = time.Now()
	}
	tag, err := s.db.Exec(ctx, `
		UPDATE articles
		SET status = $3,
		    published_at = COALESCE(articles.published_at, $4),
		    updated_at = NOW()
		WHERE id = $1 AND author_id = $2
	`, id, authorID, status, publishedAt)
	if err != nil {
		return fmt.Errorf("set status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SlugExists reports whether a slug is already taken.
func (s *Store) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM articles WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

// GetByID loads an article with all translations, scoped to an author.
func (s *Store) GetByID(ctx context.Context, id, authorID uuid.UUID) (*Article, error) {
	row := s.db.QueryRow(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.id = $1 AND a.author_id = $2
	`, id, authorID)
	art, err := scanArticle(row)
	if err != nil {
		return nil, err
	}
	if err := s.loadTranslations(ctx, art); err != nil {
		return nil, err
	}
	return art, nil
}

// GetPublishedBySlug loads a published article with all translations.
func (s *Store) GetPublishedBySlug(ctx context.Context, slug string) (*Article, error) {
	row := s.db.QueryRow(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.slug = $1 AND a.status = 'published'
	`, slug)
	art, err := scanArticle(row)
	if err != nil {
		return nil, err
	}
	if err := s.loadTranslations(ctx, art); err != nil {
		return nil, err
	}
	return art, nil
}

// ListPublished returns published articles for the feed. sort "top" orders by
// score (readers' choice); anything else by recency. A non-empty category
// filters to that rubric.
func (s *Store) ListPublished(ctx context.Context, sort, category, subcategory string, limit, offset int) ([]*Article, error) {
	if limit <= 0 || limit > 60 {
		limit = 24
	}
	orderBy := "a.published_at DESC NULLS LAST"
	if sort == "top" {
		orderBy = "a.score DESC, a.published_at DESC NULLS LAST"
	}

	where := "a.status = 'published'"
	args := []any{}
	if category != "" {
		args = append(args, category)
		where += fmt.Sprintf(" AND a.category = $%d", len(args))
	}
	if subcategory != "" {
		args = append(args, subcategory)
		where += fmt.Sprintf(" AND a.subcategory = $%d", len(args))
	}
	args = append(args, limit)
	limIdx := len(args)
	args = append(args, offset)
	offIdx := len(args)

	query := fmt.Sprintf(`
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, where, orderBy, limIdx, offIdx)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list published: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// ListPublishedByAuthor returns an author's published articles, newest first.
func (s *Store) ListPublishedByAuthor(ctx context.Context, authorID string, limit int) ([]*Article, error) {
	if limit <= 0 || limit > 100 {
		limit = 60
	}
	id, err := uuid.Parse(authorID)
	if err != nil {
		return nil, fmt.Errorf("author id: %w", err)
	}
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.status = 'published' AND a.author_id = $1
		ORDER BY a.published_at DESC NULLS LAST
		LIMIT $2
	`, id, limit)
	if err != nil {
		return nil, fmt.Errorf("list by author: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// ListByAuthor returns all of an author's articles, newest first.
func (s *Store) ListByAuthor(ctx context.Context, authorID uuid.UUID) ([]*Article, error) {
	rows, err := s.db.Query(ctx, `
		SELECT a.id, a.author_id, u.email, a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.author_id = $1
		ORDER BY a.updated_at DESC
	`, authorID)
	if err != nil {
		return nil, fmt.Errorf("list by author: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// RecordView increments the aggregate and per-day view counters.
func (s *Store) RecordView(ctx context.Context, articleID uuid.UUID, lang string) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `UPDATE articles SET views_count = views_count + 1 WHERE id = $1`, articleID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO article_views_daily (article_id, lang, day, views)
		VALUES ($1, $2, CURRENT_DATE, 1)
		ON CONFLICT (article_id, lang, day) DO UPDATE SET views = article_views_daily.views + 1
	`, articleID, lang); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AuthorStats aggregates dashboard metrics for one author.
func (s *Store) AuthorStats(ctx context.Context, authorID uuid.UUID) (AuthorStats, error) {
	var st AuthorStats
	st.ViewsByLang = map[string]int64{}

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'published'),
			COUNT(*) FILTER (WHERE status = 'draft'),
			COALESCE(SUM(views_count), 0)
		FROM articles WHERE author_id = $1
	`, authorID).Scan(&st.TotalArticles, &st.Published, &st.Drafts, &st.TotalViews)
	if err != nil {
		return st, fmt.Errorf("author stats: %w", err)
	}

	rows, err := s.db.Query(ctx, `
		SELECT v.lang, COALESCE(SUM(v.views), 0)
		FROM article_views_daily v
		JOIN articles a ON a.id = v.article_id
		WHERE a.author_id = $1
		GROUP BY v.lang
	`, authorID)
	if err != nil {
		return st, fmt.Errorf("views by lang: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var lang string
		var n int64
		if err := rows.Scan(&lang, &n); err != nil {
			return st, err
		}
		st.ViewsByLang[lang] = n
	}
	return st, rows.Err()
}

func (s *Store) loadTranslations(ctx context.Context, art *Article) error {
	rows, err := s.db.Query(ctx, `
		SELECT lang, title, summary, body_md, source, status
		FROM article_translations WHERE article_id = $1
	`, art.ID)
	if err != nil {
		return fmt.Errorf("load translations: %w", err)
	}
	defer rows.Close()
	art.Translations = map[string]*Translation{}
	for rows.Next() {
		var t Translation
		if err := rows.Scan(&t.Lang, &t.Title, &t.Summary, &t.BodyMD, &t.Source, &t.Status); err != nil {
			return err
		}
		tr := t
		art.Translations[t.Lang] = &tr
	}
	return rows.Err()
}

func (s *Store) attachTranslations(ctx context.Context, arts []*Article) ([]*Article, error) {
	for _, a := range arts {
		if err := s.loadTranslations(ctx, a); err != nil {
			return nil, err
		}
	}
	return arts, nil
}

func scanArticle(row pgx.Row) (*Article, error) {
	var a Article
	err := row.Scan(&a.ID, &a.AuthorID, &a.AuthorEmail, &a.Slug, &a.OriginalLang, &a.Status, &a.Category, &a.Subcategory,
		&a.CoverURL, &a.Score, &a.ViewsCount, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan article: %w", err)
	}
	a.Translations = map[string]*Translation{}
	return &a, nil
}

func scanArticles(rows pgx.Rows) ([]*Article, error) {
	defer rows.Close()
	var out []*Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.AuthorID, &a.AuthorEmail, &a.Slug, &a.OriginalLang, &a.Status, &a.Category, &a.Subcategory,
			&a.CoverURL, &a.Score, &a.ViewsCount, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan article row: %w", err)
		}
		a.Translations = map[string]*Translation{}
		art := a
		out = append(out, &art)
	}
	return out, rows.Err()
}
