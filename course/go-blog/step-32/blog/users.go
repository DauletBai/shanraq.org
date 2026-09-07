package blog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	sqlite "modernc.org/sqlite"
)

// What the rest of the program sees. It knows nothing about bcrypt or about
// the numbers the database answers with.
var (
	ErrEmailTaken = errors.New("бұл пошта тіркелген")
	ErrBadLogin   = errors.New("пошта немесе құпиясөз келмейді")
	ErrShort      = errors.New("құпиясөз сегіз таңбадан қысқа")
	ErrLong       = errors.New("құпиясөз 72 байттан ұзын")
)

// MinPassword is counted in characters, for a person. MaxPasswordBytes is
// counted in bytes, for bcrypt, whose limit is exactly 72 -- which is only 36
// Kazakh letters, because each of them takes two bytes.
const (
	MinPassword      = 8
	MaxPasswordBytes = 72
)

// dummyHash is a real hash of a string nobody knows. It is what makes an
// address that does not exist cost the same time as one that does: otherwise
// a stopwatch would tell an outsider who is registered here.
const dummyHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

// NormalizeEmail is what makes Aigul@Example.KZ and aigul@example.kz one
// person, so that unique in the schema means what it should.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Register files a person. The password never reaches the database; its hash
// does, and the salt travels inside that hash.
func (s *Store) Register(ctx context.Context, email, password string) error {
	email = NormalizeEmail(email)
	if utf8.RuneCountInString(password) < MinPassword {
		return ErrShort
	}
	if len(password) > MaxPasswordBytes {
		return ErrLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("хеш: %w", err)
	}

	_, err = s.db.ExecContext(ctx,
		`insert into users (email, password_hash) values (?, ?)`, email, string(hash))
	if err != nil {
		var serr *sqlite.Error
		if errors.As(err, &serr) && serr.Code() == uniqueViolation {
			return ErrEmailTaken
		}
		return fmt.Errorf("тіркеу: %w", err)
	}
	return nil
}

// Authenticate answers the same way to a wrong password and to an address that
// was never registered -- in words and in time.
func (s *Store) Authenticate(ctx context.Context, email, password string) (int64, error) {
	email = NormalizeEmail(email)

	var id int64
	var hash string
	err := s.db.QueryRowContext(ctx,
		`select id, password_hash from users where email = ?`, email).Scan(&id, &hash)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		// The result is not needed; the time spent is.
		bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))
		return 0, ErrBadLogin
	case err != nil:
		return 0, fmt.Errorf("кіру: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, ErrBadLogin
	}
	return id, nil
}
