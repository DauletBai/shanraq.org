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

// Open opens the database file and makes sure the schema is there. The pragma
// goes in the connection string, not in a query: the pool holds several
// connections and a query would reach only one of them.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("дерекқорды ашу: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("дерекқор жауап бермейді: %w", err)
	}
	if _, err := db.Exec(`create table if not exists articles (
		id    integer primary key,
		slug  text    not null unique,
		title text    not null,
		words integer not null default 0,
		lang  text    not null default 'kz',
		body  text    not null default ''
	) strict`); err != nil {
		db.Close()
		return nil, fmt.Errorf("схема: %w", err)
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
		`select slug, title, words, lang, body from articles where slug = ?`, slug).
		Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body)
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
	rows, err := s.db.Query(`select slug, title, words, lang, body from articles order by id`)
	if err != nil {
		return nil, fmt.Errorf("тізім: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body); err != nil {
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
