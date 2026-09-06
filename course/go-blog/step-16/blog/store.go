package blog

import (
	"database/sql"
	"errors"
	"fmt"

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
func (s *Store) Add(a Article) error {
	_, err := s.db.Exec(
		`insert into articles (slug, title, words, lang, body) values (?, ?, ?, ?, ?)`,
		a.Slug, a.Title, a.Words, a.Lang, a.Body)
	if err != nil {
		return fmt.Errorf("%q қосу: %w", a.Slug, err)
	}
	return nil
}

func (s *Store) Get(slug string) (Article, error) {
	var a Article
	err := s.db.QueryRow(
		`select slug, title, words, lang, body, updated_at from articles where slug = ?`, slug).
		Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body, &a.UpdatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return Article{}, fmt.Errorf("%q оқу: %w", slug, err)
	}
	return a, nil
}

// All hands back the articles in the order they were added.
func (s *Store) All() ([]Article, error) {
	rows, err := s.db.Query(`select slug, title, words, lang, body, updated_at from articles order by id`)
	if err != nil {
		return nil, fmt.Errorf("тізім: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("жолды талдау: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) Len() (int, error) {
	var n int
	if err := s.db.QueryRow(`select count(*) from articles`).Scan(&n); err != nil {
		return 0, fmt.Errorf("санау: %w", err)
	}
	return n, nil
}

// Update changes the title of an article that is already filed. A row that is
// not there is not an error to the database: it matches nothing and reports
// success, so RowsAffected is the only place the truth shows up.
func (s *Store) Update(slug, title string) error {
	res, err := s.db.Exec(
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
