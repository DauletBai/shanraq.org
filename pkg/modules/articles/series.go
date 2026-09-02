package articles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Article series -- an ordered course assembled from ordinary articles.
//
// Every lesson stays a normal article on its own URL: indexable, translatable,
// commentable. The series only supplies order and a hub page, because the course
// earns its traffic lesson by lesson from search, not from the hub. Gating the
// text would remove the only reason to build it.

// Series statuses. A draft series is reachable by direct link but stays out of
// listings, so a course can be assembled in the open before it is announced.
const (
	SeriesDraft     = "draft"
	SeriesPublished = "published"
)

// Series is a course: metadata in every language plus its ordered lessons.
type Series struct {
	ID       uuid.UUID
	Slug     string
	CoverURL string
	Status   string

	// Title and Summary keyed by language code.
	Title   map[string]string
	Summary map[string]string

	// Items are the lessons in reading order; empty unless loaded.
	Items []*SeriesItem

	CreatedAt time.Time
	UpdatedAt time.Time
}

// SeriesItem is one lesson's place in a course.
type SeriesItem struct {
	ArticleID uuid.UUID
	Position  int
	Slug      string
	Title     string
	Summary   string
	Minutes   int
	Published bool
}

// TitleIn returns the course title in the reader's language, falling back to any
// language that has one -- a course shown with a blank name would be worse than
// one shown in the wrong language.
func (s *Series) TitleIn(lang string) string { return pickLang(s.Title, lang) }

// SummaryIn returns the course blurb in the reader's language.
func (s *Series) SummaryIn(lang string) string { return pickLang(s.Summary, lang) }

// Published reports whether the course is listed publicly.
func (s *Series) IsPublished() bool { return s.Status == SeriesPublished }

// Lessons counts the published lessons -- the number a reader is shown, which is
// not the number of rows when part of the course is still being written.
func (s *Series) Lessons() int {
	n := 0
	for _, it := range s.Items {
		if it.Published {
			n++
		}
	}
	return n
}

// Minutes totals the reading time of the published lessons.
func (s *Series) Minutes() int {
	n := 0
	for _, it := range s.Items {
		if it.Published {
			n += it.Minutes
		}
	}
	return n
}

// SeriesPlace locates one article inside a course: which lesson it is and what
// sits either side of it. This is what the strip on a lesson page needs.
type SeriesPlace struct {
	Series *Series
	Number int // 1-based position among published lessons
	Total  int
	Prev   *SeriesItem
	Next   *SeriesItem
}

// SeriesStore reads and writes courses.
type SeriesStore struct{ db *pgxpool.Pool }

// NewSeriesStore wires the store to the pool.
func NewSeriesStore(db *pgxpool.Pool) *SeriesStore { return &SeriesStore{db: db} }

// ErrSeriesNotFound is returned when no course has the requested slug or id.
var ErrSeriesNotFound = errors.New("series not found")

const seriesCols = `s.id, s.slug, s.cover_url, s.status, s.created_at, s.updated_at`

func scanSeries(row pgx.Row) (*Series, error) {
	s := &Series{Title: map[string]string{}, Summary: map[string]string{}}
	if err := row.Scan(&s.ID, &s.Slug, &s.CoverURL, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return s, nil
}

// loadI18n fills titles and blurbs for the given courses in one query.
func (st *SeriesStore) loadI18n(ctx context.Context, list []*Series) error {
	if len(list) == 0 {
		return nil
	}
	byID := make(map[uuid.UUID]*Series, len(list))
	ids := make([]uuid.UUID, 0, len(list))
	for _, s := range list {
		byID[s.ID] = s
		ids = append(ids, s.ID)
	}
	rows, err := st.db.Query(ctx,
		`SELECT series_id, lang, title, summary FROM article_series_i18n WHERE series_id = ANY($1)`, ids)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var lang, title, summary string
		if err := rows.Scan(&id, &lang, &title, &summary); err != nil {
			return err
		}
		if s := byID[id]; s != nil {
			s.Title[lang] = title
			s.Summary[lang] = summary
		}
	}
	return rows.Err()
}

// List returns the published courses, newest first. Drafts are excluded: they
// are reachable by link but must not appear in a menu.
func (st *SeriesStore) List(ctx context.Context) ([]*Series, error) {
	rows, err := st.db.Query(ctx,
		`SELECT `+seriesCols+` FROM article_series s WHERE s.status = $1 ORDER BY s.created_at DESC`,
		SeriesPublished)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Series
	for rows.Next() {
		s, err := scanSeries(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, st.loadI18n(ctx, list)
}

// ListAll returns every course including drafts, for the admin screen.
func (st *SeriesStore) ListAll(ctx context.Context) ([]*Series, error) {
	rows, err := st.db.Query(ctx, `SELECT `+seriesCols+` FROM article_series s ORDER BY s.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Series
	for rows.Next() {
		s, err := scanSeries(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, st.loadI18n(ctx, list)
}

// BySlug loads one course with its lessons in reading order, titled in lang.
func (st *SeriesStore) BySlug(ctx context.Context, slug, lang string) (*Series, error) {
	s, err := scanSeries(st.db.QueryRow(ctx,
		`SELECT `+seriesCols+` FROM article_series s WHERE s.slug = $1`, slug))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	if err := st.loadI18n(ctx, []*Series{s}); err != nil {
		return nil, err
	}
	if s.Items, err = st.items(ctx, s.ID, lang); err != nil {
		return nil, err
	}
	return s, nil
}

// items loads a course's lessons. The title comes from the reader's language
// where it exists and from the article's original language otherwise, so a
// lesson that has not been translated yet still shows a name rather than a gap.
func (st *SeriesStore) items(ctx context.Context, seriesID uuid.UUID, lang string) ([]*SeriesItem, error) {
	rows, err := st.db.Query(ctx, `
		SELECT i.article_id, i.position, a.slug, a.status,
		       COALESCE(tr.title, orig.title, ''),
		       COALESCE(tr.summary, orig.summary, ''),
		       COALESCE(tr.body_md, orig.body_md, '')
		FROM article_series_items i
		JOIN articles a ON a.id = i.article_id
		LEFT JOIN article_translations tr ON tr.article_id = a.id AND tr.lang = $2
		LEFT JOIN article_translations orig ON orig.article_id = a.id AND orig.lang = a.original_lang
		WHERE i.series_id = $1
		ORDER BY i.position, a.created_at`, seriesID, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*SeriesItem
	for rows.Next() {
		it := &SeriesItem{}
		var status, body string
		if err := rows.Scan(&it.ArticleID, &it.Position, &it.Slug, &status, &it.Title, &it.Summary, &body); err != nil {
			return nil, err
		}
		it.Published = status == "published"
		it.Minutes = readingMinutes(body)
		list = append(list, it)
	}
	return list, rows.Err()
}

// ForArticle returns the article's place in each course it belongs to. Prev and
// Next skip unpublished lessons: a reader following the chain must never be sent
// to a draft, and a gap in the middle of a course under construction would
// otherwise dead-end them.
func (st *SeriesStore) ForArticle(ctx context.Context, articleID uuid.UUID, lang string) ([]*SeriesPlace, error) {
	// Published courses only. A draft course is meant to be reachable by its own
	// link and nowhere else, but a lesson of it that went live carried the strip
	// and the link out onto a public page — so the course announced itself
	// through the very article that was supposed to be independent of it.
	rows, err := st.db.Query(ctx,
		`SELECT `+seriesCols+` FROM article_series s
		 JOIN article_series_items i ON i.series_id = s.id
		 WHERE i.article_id = $1 AND s.status = $2 ORDER BY s.created_at`,
		articleID, SeriesPublished)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var found []*Series
	for rows.Next() {
		s, err := scanSeries(rows)
		if err != nil {
			return nil, err
		}
		found = append(found, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := st.loadI18n(ctx, found); err != nil {
		return nil, err
	}

	places := make([]*SeriesPlace, 0, len(found))
	for _, s := range found {
		items, err := st.items(ctx, s.ID, lang)
		if err != nil {
			return nil, err
		}
		s.Items = items
		place := &SeriesPlace{Series: s}
		var pub []*SeriesItem
		for _, it := range items {
			if it.Published {
				pub = append(pub, it)
			}
		}
		place.Total = len(pub)
		for i, it := range pub {
			if it.ArticleID == articleID {
				place.Number = i + 1
				if i > 0 {
					place.Prev = pub[i-1]
				}
				if i+1 < len(pub) {
					place.Next = pub[i+1]
				}
				break
			}
		}
		places = append(places, place)
	}
	return places, nil
}

// Save creates or updates a course and its per-language text. Returns the id.
func (st *SeriesStore) Save(ctx context.Context, id *uuid.UUID, slug, coverURL, status string, title, summary map[string]string) (uuid.UUID, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return uuid.Nil, fmt.Errorf("series slug is required")
	}
	if status != SeriesPublished {
		status = SeriesDraft
	}

	tx, err := st.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var sid uuid.UUID
	if id != nil && *id != uuid.Nil {
		sid = *id
		if _, err := tx.Exec(ctx,
			`UPDATE article_series SET slug = $2, cover_url = $3, status = $4, updated_at = now() WHERE id = $1`,
			sid, slug, coverURL, status); err != nil {
			return uuid.Nil, err
		}
	} else if err := tx.QueryRow(ctx,
		`INSERT INTO article_series (slug, cover_url, status) VALUES ($1, $2, $3) RETURNING id`,
		slug, coverURL, status).Scan(&sid); err != nil {
		return uuid.Nil, err
	}

	for _, lang := range Langs {
		t, sm := strings.TrimSpace(title[lang]), strings.TrimSpace(summary[lang])
		if t == "" && sm == "" {
			if _, err := tx.Exec(ctx,
				`DELETE FROM article_series_i18n WHERE series_id = $1 AND lang = $2`, sid, lang); err != nil {
				return uuid.Nil, err
			}
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO article_series_i18n (series_id, lang, title, summary)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (series_id, lang) DO UPDATE SET title = EXCLUDED.title, summary = EXCLUDED.summary`,
			sid, lang, t, sm); err != nil {
			return uuid.Nil, err
		}
	}
	return sid, tx.Commit(ctx)
}

// Attach puts an article into a course at the given position, or moves it there
// if it is already a member.
func (st *SeriesStore) Attach(ctx context.Context, seriesID, articleID uuid.UUID, position int) error {
	_, err := st.db.Exec(ctx, `
		INSERT INTO article_series_items (series_id, article_id, position)
		VALUES ($1, $2, $3)
		ON CONFLICT (series_id, article_id) DO UPDATE SET position = EXCLUDED.position`,
		seriesID, articleID, position)
	return err
}

// Detach removes an article from a course. The article itself is untouched.
func (st *SeriesStore) Detach(ctx context.Context, seriesID, articleID uuid.UUID) error {
	_, err := st.db.Exec(ctx,
		`DELETE FROM article_series_items WHERE series_id = $1 AND article_id = $2`, seriesID, articleID)
	return err
}

// Delete removes a course. Its lessons stay published as ordinary articles --
// deleting a course must never take reader-facing content down with it.
func (st *SeriesStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := st.db.Exec(ctx, `DELETE FROM article_series WHERE id = $1`, id)
	return err
}

// ByID loads one course with its lessons, for the admin editor.
func (st *SeriesStore) ByID(ctx context.Context, id uuid.UUID, lang string) (*Series, error) {
	s, err := scanSeries(st.db.QueryRow(ctx,
		`SELECT `+seriesCols+` FROM article_series s WHERE s.id = $1`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSeriesNotFound
		}
		return nil, err
	}
	if err := st.loadI18n(ctx, []*Series{s}); err != nil {
		return nil, err
	}
	if s.Items, err = st.items(ctx, s.ID, lang); err != nil {
		return nil, err
	}
	return s, nil
}

// Candidates lists articles that can be added to a course, newest first.
// Drafts are included: a lesson should take its place in the running order
// before it goes live, not after.
func (st *SeriesStore) Candidates(ctx context.Context, lang string, limit int) ([]predArticleOption, error) {
	if limit <= 0 {
		limit = 200
	}
	rows, err := st.db.Query(ctx, `
		SELECT a.slug, COALESCE(tr.title, orig.title, a.slug), a.status
		FROM articles a
		LEFT JOIN article_translations tr ON tr.article_id = a.id AND tr.lang = $1
		LEFT JOIN article_translations orig ON orig.article_id = a.id AND orig.lang = a.original_lang
		WHERE a.status <> 'archived'
		ORDER BY a.created_at DESC
		LIMIT $2`, lang, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []predArticleOption
	for rows.Next() {
		var slug, title, status string
		if err := rows.Scan(&slug, &title, &status); err != nil {
			return nil, err
		}
		if status != "published" {
			title += " (" + status + ")"
		}
		out = append(out, predArticleOption{ID: slug, Title: title})
	}
	return out, rows.Err()
}

// AttachBySlug places an article in a course by its slug, which is what the
// admin form has to hand.
func (st *SeriesStore) AttachBySlug(ctx context.Context, seriesID uuid.UUID, slug string, position int) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("article slug is required")
	}
	var aid uuid.UUID
	if err := st.db.QueryRow(ctx, `SELECT id FROM articles WHERE slug = $1`, slug).Scan(&aid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("no article with slug %q", slug)
		}
		return err
	}
	return st.Attach(ctx, seriesID, aid, position)
}
