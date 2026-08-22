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

// DeleteDraft removes an author's own article, but only while it is a draft.
//
// A published piece is deliberately out of reach: it has views, votes, comments,
// an RSS entry and a Telegram post behind it, and somewhere a reader may already
// have the link. Losing all of that to one mis-click is not a trade worth
// offering. The way out exists and is one step longer — unpublish, then delete —
// which is exactly the pause such a decision deserves.
//
// Returns ErrNotFound when the article does not exist, belongs to someone else,
// or is still published; the caller cannot tell these apart, which is also what
// stops the endpoint from confirming that a stranger's article ID is real.
func (s *Store) DeleteDraft(ctx context.Context, id, authorID uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("delete draft: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	tag, err := tx.Exec(ctx, `
		DELETE FROM articles
		WHERE id = $1 AND author_id = $2 AND status = 'draft'
	`, id, authorID)
	if err != nil {
		return fmt.Errorf("delete draft: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	// Translations, ratings and comments carry ON DELETE CASCADE and are gone
	// already. These two are keyed by article id without a foreign key — the
	// funnel because it is a plain counter table, favourites because they are
	// polymorphic over articles and listings — so they have to be swept by hand
	// or they linger as rows pointing at nothing.
	if _, err := tx.Exec(ctx, `DELETE FROM reading_depth WHERE article_id = $1`, id); err != nil {
		return fmt.Errorf("delete draft depth: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM favorites WHERE item_type = 'article' AND item_id = $1`, id); err != nil {
		return fmt.Errorf("delete draft favorites: %w", err)
	}
	// The moderation ledger is deliberately left intact: it is an append-only
	// record of decisions, and an author must not be able to erase it.
	return tx.Commit(ctx)
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
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at, a.indexable
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
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at, a.indexable
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

// placeClause narrows a feed to what its reader is addressed by: material with
// no place at all, plus material written for the places that contain them.
//
// A reader who never said where they live — every guest among them — is
// addressed by nothing, and sees only the material that was written for
// everyone. A power cut in one town is not news to somebody a thousand
// kilometres away; it is clutter, and it pushes down what they came for.
func placeClause(args *[]any, addressed []uuid.UUID) string {
	if len(addressed) == 0 {
		return " AND a.geo_node_id IS NULL"
	}
	*args = append(*args, addressed)
	return fmt.Sprintf(" AND (a.geo_node_id IS NULL OR a.geo_node_id = ANY($%d))", len(*args))
}

// ListPublished returns published articles for the feed. sort "top" orders by
// score (readers' choice); anything else by recency. A non-empty category
// filters to that rubric.
func (s *Store) ListPublished(ctx context.Context, sort, category, subcategory string, limit, offset int, addressed []uuid.UUID) ([]*Article, error) {
	if limit <= 0 || limit > 60 {
		limit = 24
	}
	// The id is a tiebreaker, not decoration. Articles published in the same
	// batch share a published_at to the microsecond, and rows tied under
	// ORDER BY have no defined order between one query and the next — so paging
	// by OFFSET could show an article on two pages and never show another at
	// all. Nobody noticed while the feed was a single page.
	orderBy := "a.published_at DESC NULLS LAST, a.id DESC"
	if sort == "top" {
		orderBy = "a.score DESC, a.published_at DESC NULLS LAST, a.id DESC"
	}

	args := []any{}
	where := "a.status = 'published'" + placeClause(&args, addressed)
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
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at, a.indexable
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
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at, a.indexable
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
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug, a.original_lang, a.status, a.category, a.subcategory,
		       a.cover_url, a.score, a.views_count, a.published_at, a.created_at, a.updated_at, a.indexable
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
	err := row.Scan(&a.ID, &a.AuthorID, &a.AuthorEmail, &a.AuthorFirst, &a.AuthorLast, &a.Slug, &a.OriginalLang, &a.Status, &a.Category, &a.Subcategory,
		&a.CoverURL, &a.Score, &a.ViewsCount, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &a.Indexable)
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
		err := rows.Scan(&a.ID, &a.AuthorID, &a.AuthorEmail, &a.AuthorFirst, &a.AuthorLast, &a.Slug, &a.OriginalLang, &a.Status, &a.Category, &a.Subcategory,
			&a.CoverURL, &a.Score, &a.ViewsCount, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &a.Indexable)
		if err != nil {
			return nil, fmt.Errorf("scan article row: %w", err)
		}
		a.Translations = map[string]*Translation{}
		art := a
		out = append(out, &art)
	}
	return out, rows.Err()
}

// RelatedPublished returns a handful of other articles to offer at the end of
// one, nearest first: same subcategory, then same category, then simply recent.
//
// An article page used to link to no other article at all. The site's whole
// link graph ran from the home page outwards and stopped: a reader who finished
// a piece had nowhere to go, and every article was a leaf with nothing pointing
// out of it — which is also how a crawler decides a page sits at the edge of a
// site rather than inside it.
//
// Only indexable articles are offered. The machine-written columns are open to
// readers who choose them and closed to search engines; putting them under a
// human piece would push a machine's opinion at somebody who did not ask for
// it, and spend a crawler's visit on a page that cannot be indexed anyway.
func (s *Store) RelatedPublished(ctx context.Context, exclude uuid.UUID, category, subcategory string, limit int, addressed []uuid.UUID) ([]*Article, error) {
	if limit <= 0 || limit > 12 {
		limit = 4
	}
	args := []any{exclude, subcategory, category, limit}
	rows, err := s.db.Query(ctx, fmt.Sprintf(`
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug,
		       a.original_lang, a.status, a.category, a.subcategory, a.cover_url, a.score, a.views_count,
		       a.published_at, a.created_at, a.updated_at, a.indexable
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.status = 'published' AND a.indexable AND a.id <> $1%s
		ORDER BY
			(a.subcategory = $2 AND $2 <> '') DESC,
			(a.category = $3) DESC,
			a.published_at DESC NULLS LAST,
			a.id DESC
		LIMIT $4
	`, placeClause(&args, addressed)), args...)
	if err != nil {
		return nil, fmt.Errorf("related published: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// CategoryFreshness returns, per category, when its newest article appeared.
//
// It feeds <lastmod> on the category feeds in the sitemap. Those pages are real
// hubs — they change whenever a piece lands in the section — and a crawler that
// knows the date has a reason to come back for it. The static pages get no date
// on purpose: /about and /terms do not change, and inventing a modification
// date for them would teach Google that our dates mean nothing, which costs
// more than the entries are worth.
func (s *Store) CategoryFreshness(ctx context.Context) (map[string]time.Time, error) {
	rows, err := s.db.Query(ctx, `
		SELECT category, MAX(published_at)
		FROM articles
		WHERE status = 'published' AND indexable AND published_at IS NOT NULL
		GROUP BY category`)
	if err != nil {
		return nil, fmt.Errorf("category freshness: %w", err)
	}
	defer rows.Close()
	out := map[string]time.Time{}
	for rows.Next() {
		var cat string
		var at *time.Time
		if err := rows.Scan(&cat, &at); err != nil {
			return nil, err
		}
		if at != nil {
			out[cat] = *at
		}
	}
	return out, rows.Err()
}

// ListForPlace returns the published articles addressed to a place — the ones
// tied to it, and the ones tied to any place that contains it — newest first.
//
// Upwards, not downwards. The place an author picks is the audience they wrote
// for: a notice for Kachar is meant for Kachar, and putting it on the oblast
// page would hand it to a hundred thousand people who were never its readers.
// An oblast-wide announcement, on the other hand, is addressed to everyone in
// the oblast, Kachar included, so it belongs on Kachar's page too.
//
// Non-indexable articles are excluded for the same reason they are excluded
// everywhere else: a place page is a page we want found, and it should carry
// what a person came for.
func (s *Store) ListForPlace(ctx context.Context, place uuid.UUID, limit, offset int) ([]*Article, error) {
	if limit <= 0 || limit > 60 {
		limit = 24
	}
	rows, err := s.db.Query(ctx, `
		WITH RECURSIVE up AS (
			SELECT id, parent_id FROM geo_nodes WHERE id = $1
			UNION ALL
			SELECT g.id, g.parent_id FROM geo_nodes g JOIN up ON g.id = up.parent_id
		)
		SELECT a.id, a.author_id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), a.slug,
		       a.original_lang, a.status, a.category, a.subcategory, a.cover_url, a.score, a.views_count,
		       a.published_at, a.created_at, a.updated_at, a.indexable
		FROM articles a
		JOIN auth_users u ON u.id = a.author_id
		WHERE a.status = 'published' AND a.indexable AND a.geo_node_id IN (SELECT id FROM up)
		ORDER BY a.published_at DESC NULLS LAST, a.id DESC
		LIMIT $2 OFFSET $3
	`, place, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list for place: %w", err)
	}
	arts, err := scanArticles(rows)
	if err != nil {
		return nil, err
	}
	return s.attachTranslations(ctx, arts)
}

// PlacesWithArticles lists the places something was published for, so the
// sitemap and the place index carry pages with content rather than eight
// hundred empty ones.
//
// Only the place the author chose. A region's article also shows on the page of
// every town inside that region, but offering thirty pages that each carry the
// same single article is thin content thirty times over, not reach.
func (s *Store) PlacesWithArticles(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `
		SELECT DISTINCT g.slug
		FROM articles a JOIN geo_nodes g ON g.id = a.geo_node_id
		WHERE a.status = 'published' AND a.indexable
		  AND g.slug IS NOT NULL AND g.slug <> ''
		ORDER BY g.slug`)
	if err != nil {
		return nil, fmt.Errorf("places with articles: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// SetArticlePlace ties an article to a place, or unties it when node is nil.
// Scoped to the author in the statement, so another account's id changes
// nothing.
func (s *Store) SetArticlePlace(ctx context.Context, id, author uuid.UUID, node *uuid.UUID) error {
	_, err := s.db.Exec(ctx,
		`UPDATE articles SET geo_node_id = $3, updated_at = NOW() WHERE id = $1 AND author_id = $2`,
		id, author, node)
	if err != nil {
		return fmt.Errorf("set article place: %w", err)
	}
	return nil
}

// ArticlePlace returns the place an article was published for, if any.
func (s *Store) ArticlePlace(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	var node *uuid.UUID
	if err := s.db.QueryRow(ctx, `SELECT geo_node_id FROM articles WHERE id = $1`, id).Scan(&node); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("article place: %w", err)
	}
	return node, nil
}
