package articles

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// editablePageKeys lists the info/legal pages an operator may edit from the
// admin panel, in display order. Every key must exist in staticPages, which
// supplies both the seed content and the fallback if a row is ever missing.
var editablePageKeys = []string{"about", "guide", "pricing", "support", "formatting", "privacy", "terms"}

// ContentStore persists the editable pages (title + Markdown body per page and
// language). It backs both the public reader and the admin editor.
type ContentStore struct{ db *pgxpool.Pool }

// NewContentStore builds a ContentStore over the shared pgx pool.
func NewContentStore(db *pgxpool.Pool) *ContentStore { return &ContentStore{db: db} }

// Get returns the stored page for (key, lang). ok is false when no row exists,
// so callers fall back to the built-in default.
func (s *ContentStore) Get(ctx context.Context, key, lang string) (title, body string, ok bool, err error) {
	err = s.db.QueryRow(ctx,
		`SELECT title, body_md FROM content_pages WHERE page_key=$1 AND lang=$2`,
		key, lang).Scan(&title, &body)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	return title, body, true, nil
}

// Upsert stores an edited page for (key, lang), stamping who changed it.
func (s *ContentStore) Upsert(ctx context.Context, key, lang, title, body string, editor uuid.UUID) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO content_pages (page_key, lang, title, body_md, updated_at, updated_by)
		 VALUES ($1,$2,$3,$4,NOW(),$5)
		 ON CONFLICT (page_key, lang) DO UPDATE
		   SET title=EXCLUDED.title, body_md=EXCLUDED.body_md,
		       updated_at=NOW(), updated_by=EXCLUDED.updated_by`,
		key, lang, title, body, editor)
	return err
}

// seed inserts a default row only when absent, so first boot fills the table
// from the built-in pages without ever clobbering a later edit.
func (s *ContentStore) seed(ctx context.Context, key, lang, title, body string) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO content_pages (page_key, lang, title, body_md)
		 VALUES ($1,$2,$3,$4) ON CONFLICT (page_key, lang) DO NOTHING`,
		key, lang, title, body)
	return err
}

// LastEdited reports when and by whom a page was most recently changed (across
// its languages), for the admin editor's audit line.
func (s *ContentStore) LastEdited(ctx context.Context, key string) (when time.Time, byEmail string, ok bool, err error) {
	err = s.db.QueryRow(ctx, `
		SELECT cp.updated_at, COALESCE(u.email, '')
		  FROM content_pages cp
		  LEFT JOIN auth_users u ON u.id = cp.updated_by
		 WHERE cp.page_key = $1
		 ORDER BY cp.updated_at DESC
		 LIMIT 1`, key).Scan(&when, &byEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, "", false, nil
	}
	if err != nil {
		return time.Time{}, "", false, err
	}
	return when, byEmail, true, nil
}

// seedContentPages fills content_pages from the built-in staticPages the first
// time the app boots against a fresh table. Idempotent (ON CONFLICT DO NOTHING)
// and best-effort — a failure just leaves the built-in text serving as before.
func (m *Module) seedContentPages(ctx context.Context) {
	for _, key := range editablePageKeys {
		for _, lang := range []string{"kz", "ru", "en"} {
			c := staticContent(key, lang)
			if c.Title == "" {
				continue
			}
			if err := m.content.seed(ctx, key, lang, c.Title, c.Body); err != nil {
				// Likely the table isn't there yet — stop; the fallback below keeps
				// serving the built-in text, so pages never break.
				m.rt.Logger.Warn("seed content page", zap.String("key", key), zap.String("lang", lang), zap.Error(err))
				return
			}
		}
	}
}

// pageContent returns the effective page (DB override or built-in default) for
// (key, lang), so a missing or empty row never blanks a page.
func (m *Module) pageContent(ctx context.Context, key, lang string) (title, body string) {
	if m.content != nil {
		if t, b, ok, err := m.content.Get(ctx, key, lang); err == nil && ok && t != "" {
			return t, b
		}
	}
	c := staticContent(key, lang)
	return c.Title, c.Body
}
