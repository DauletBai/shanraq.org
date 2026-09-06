package blog

import (
	"errors"
	"fmt"
	"strings"

	sqlite "modernc.org/sqlite"
)

// uniqueViolation is the code SQLite answers a broken unique with. Codes are
// listed at sqlite.org/rescode.html and do not change between versions, which
// the text of the message does.
const uniqueViolation = 2067

// ErrTaken says the address is already in use. The handler needs this, not a
// number out of the database.
var ErrTaken = errors.New("мекенжай бос емес")

// ParseTags turns what the form sent into a list without repeats. The trimming
// and the lower case are what keep "Go" and "go" one tag rather than two.
func ParseTags(s string) []string {
	var out []string
	seen := map[string]bool{}
	for _, part := range strings.Split(s, ",") {
		t := strings.ToLower(strings.TrimSpace(part))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// AddWithTags files an article together with its tags. Two tables, one action,
// so one transaction: a tag made for an article that was not filed is rubbish
// that stays for ever.
func (s *Store) AddWithTags(a Article, tags []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`insert into articles (slug, title, words, lang, body, author_id) values (?, ?, ?, ?, ?, ?)`,
		a.Slug, a.Title, a.Words, a.Lang, a.Body, nullID(a.AuthorID))
	if err != nil {
		var serr *sqlite.Error
		if errors.As(err, &serr) && serr.Code() == uniqueViolation {
			return fmt.Errorf("%q: %w", a.Slug, ErrTaken)
		}
		return fmt.Errorf("%q қосу: %w", a.Slug, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, name := range tags {
		var tagID int64
		// do update rather than do nothing: a row that conflicts is only
		// counted as touched by an update, and returning needs it touched to
		// hand the id back.
		if err := tx.QueryRow(`insert into tags (name) values (?)
			on conflict (name) do update set name = excluded.name
			returning id`, name).Scan(&tagID); err != nil {
			return fmt.Errorf("тег %q: %w", name, err)
		}
		if _, err := tx.Exec(
			`insert into article_tags (article_id, tag_id) values (?, ?)`,
			id, tagID); err != nil {
			return fmt.Errorf("тег %q: %w", name, err)
		}
	}
	return tx.Commit()
}

// TagsOf hands back the tags of one article, in alphabetical order so the page
// does not reshuffle them between visits.
func (s *Store) TagsOf(slug string) ([]string, error) {
	rows, err := s.db.Query(`select t.name
		from tags t
		join article_tags at on at.tag_id = t.id
		join articles a on a.id = at.article_id
		where a.slug = ?
		order by t.name`, slug)
	if err != nil {
		return nil, fmt.Errorf("%q тегтері: %w", slug, err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

// ByTag lists the articles carrying one tag.
func (s *Store) ByTag(tag string) ([]Article, error) {
	rows, err := s.db.Query(`select a.slug, a.title, a.words, a.lang, a.body, a.updated_at
		from articles a
		join article_tags at on at.article_id = a.id
		join tags t on t.id = at.tag_id
		where t.name = ?
		order by a.id`, tag)
	if err != nil {
		return nil, fmt.Errorf("тег %q: %w", tag, err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		var updated string
		if err := rows.Scan(&a.Slug, &a.Title, &a.Words, &a.Lang, &a.Body, &updated); err != nil {
			return nil, err
		}
		a.UpdatedAt = parseTime(updated)
		out = append(out, a)
	}
	return out, rows.Err()
}

// nullID turns a zero into a NULL: an article filed by nobody has no owner,
// and a zero in the column would point at a user who does not exist.
func nullID(id int64) any {
	if id == 0 {
		return nil
	}
	return id
}
