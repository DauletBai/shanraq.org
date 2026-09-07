package blog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// ErrNotFound is returned when nothing is filed under that address. The caller
// compares against it with errors.Is rather than reading the message.
var ErrNotFound = errors.New("мақала табылмады")

// Store keeps the articles. Until this step it kept them in a map, and the
// comment there said the map could become a database without anything outside
// noticing. This is that day: the field changed, the methods did not — except
// that every one of them now returns an error, because a file can fail in ways
// a map never could.
type Store struct {
	db *sql.DB
}

// Open opens the database file and brings its schema up to date. The pragma
// goes in the connection string, not in a query: the pool holds several
// connections and a query would reach only one of them.
//
// The schema itself is no longer written here. It lives in migrations/, one
// numbered file per change, and migrate applies whatever this database has not
// seen yet — that is what lets the blog change shape without losing what is
// already written.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("дерекқорды ашу: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("дерекқор жауап бермейді: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("көші-қон: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Add files an article. A duplicate address is refused by the database itself,
// because slug is unique in the schema.
func (s *Store) Add(ctx context.Context, a Article) error {
	_, err := s.db.ExecContext(ctx,
		`insert into articles (slug, title, words, lang, body) values (?, ?, ?, ?, ?)`,
		a.Slug, a.Title, a.Words, a.Lang, a.Body)
	if err != nil {
		return fmt.Errorf("%q қосу: %w", a.Slug, err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, slug string) (Article, error) {
	var a Article
	var updated string
	err := s.db.QueryRowContext(ctx,
		`select slug, title, words, lang, body, updated_at, cover from articles where slug = ?`, slug).
		Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body, &updated, &a.Cover)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return Article{}, fmt.Errorf("%q оқу: %w", slug, err)
	}
	a.UpdatedAt = parseTime(updated)
	return a, nil
}

// Sorts are the orders the list can be asked for. The key comes from outside
// -- out of a link like /?sort=title -- and a column name cannot travel as a
// parameter: only a value can. So it is not cleaned or escaped but looked up
// here, in a list we wrote. Anything not on it falls back to the default, and
// that is the whole defence: it never has to guess what an attacker may send.
var sorts = map[string]string{
	"":        "id",
	"title":   "title",
	"updated": "updated_at desc, id",
}

// SortKeys are the keys a page may offer, in a fixed order.
var SortKeys = []string{"", "title", "updated"}

// All hands back the articles in the order asked for, or in the order they
// were added when the key is not one of ours.
func (s *Store) All(ctx context.Context, sort string) ([]Article, error) {
	order, ok := sorts[sort]
	if !ok {
		order = sorts[""]
	}
	rows, err := s.db.QueryContext(ctx,
		`select slug, title, words, lang, body, updated_at, cover from articles order by `+order)
	if err != nil {
		return nil, fmt.Errorf("тізім: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		var updated string
		if err := rows.Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body, &updated, &a.Cover); err != nil {
			return nil, fmt.Errorf("жолды талдау: %w", err)
		}
		a.UpdatedAt = parseTime(updated)
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) Len(ctx context.Context) (int, error) {
	var n int
	if err := s.db.QueryRowContext(ctx, `select count(*) from articles`).Scan(&n); err != nil {
		return 0, fmt.Errorf("санау: %w", err)
	}
	return n, nil
}

// Update changes the title of an article that is already filed. A row that is
// not there is not an error to the database: it matches nothing and reports
// success, so RowsAffected is the only place the truth shows up.
func (s *Store) Update(ctx context.Context, slug, title string) error {
	res, err := s.db.ExecContext(ctx,
		`update articles set title = ?, updated_at = datetime('now') where slug = ?`,
		title, slug)
	if err != nil {
		return fmt.Errorf("%q өңдеу: %w", slug, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%q өңдеу: %w", slug, err)
	}
	if n == 0 {
		return fmt.Errorf("%q: %w", slug, ErrNotFound)
	}
	return nil
}

// sqliteTime is the shape datetime('now') writes: no zone, no letter T. An
// empty string means the article has never been edited, and that is not an
// error but the zero time.
const sqliteTime = "2006-01-02 15:04:05"

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(sqliteTime, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// Backup writes a consistent copy of the database into path while the blog
// keeps answering. Copying the file with cp is not the same thing: a write in
// the middle of the copy leaves the halves disagreeing, and the copy that
// results may not open at all. VACUUM INTO is SQLite's own answer -- it walks
// the pages under a read transaction and writes a whole database out.
func (s *Store) Backup(ctx context.Context, path string) error {
	if _, err := s.db.ExecContext(ctx, "vacuum into ?", path); err != nil {
		return fmt.Errorf("көшірме жасау: %w", err)
	}
	return nil
}
