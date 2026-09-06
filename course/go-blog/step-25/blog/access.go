package blog

import (
	"database/sql"
	"errors"
	"fmt"
)

// The three answers the blog can give to "may I change this". They are errors
// rather than a bool because a handler answers each of them differently: a
// guest is sent to sign in, a stranger is refused, and what does not exist
// stays not existing.
var (
	ErrNeedLogin  = errors.New("алдымен кіріңіз")
	ErrNotAllowed = errors.New("бұл сіздің мақалаңыз емес")
)

// MayEdit is the only place where it is decided who may change what. Two such
// places would drift apart: one would learn about editors and the other would
// not.
func (s *Store) MayEdit(user User, slug string) error {
	if user.ID == 0 {
		return ErrNeedLogin
	}

	var author sql.NullInt64
	err := s.db.QueryRow(`select author_id from articles where slug = ?`, slug).Scan(&author)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return fmt.Errorf("%q: %w", slug, err)
	}

	// Checked after the article is found, so that a stranger cannot learn
	// which addresses exist by the difference between 403 and 404.
	if !author.Valid || author.Int64 != user.ID {
		return fmt.Errorf("%q: %w", slug, ErrNotAllowed)
	}
	return nil
}
