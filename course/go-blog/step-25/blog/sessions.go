package blog

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// CookieName is the name of the cookie the browser carries back and forth.
const CookieName = "session"

// SessionLife is how long a sign-in lasts: long enough not to sign in daily,
// short enough that a forgotten session does not live for years.
const SessionLife = 7 * 24 * time.Hour

// ErrNoSession says nobody is signed in with this token -- there is no such
// token, or it has run out.
var ErrNoSession = errors.New("сессия жоқ")

// User is who is reading, as far as a page is concerned.
type User struct {
	ID    int64
	Email string
}

// StartSession hands out a random token and keeps only its hash. The token
// goes into the cookie; a leaked database gives nobody a way in.
func (s *Store) StartSession(userID int64) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("кездейсоқ сан: %w", err)
	}
	// RawURLEncoding: 32 bytes become 43 characters with no +, / or = to
	// spoil a header or a link.
	token := base64.RawURLEncoding.EncodeToString(raw)

	if _, err := s.db.Exec(
		`insert into sessions (token_hash, user_id, expires_at) values (?, ?, ?)`,
		hashToken(token), userID, expiry()); err != nil {
		return "", fmt.Errorf("сессия: %w", err)
	}
	return token, nil
}

// UserByToken says who sent this request. The cookie is never believed on its
// own: everything but the random string lives in the database.
func (s *Store) UserByToken(token string) (User, error) {
	var u User
	err := s.db.QueryRow(`select u.id, u.email
		from sessions s
		join users u on u.id = s.user_id
		where s.token_hash = ? and s.expires_at > ?`,
		hashToken(token), time.Now().UTC().Format(time.RFC3339)).Scan(&u.ID, &u.Email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return User{}, ErrNoSession
	case err != nil:
		return User{}, fmt.Errorf("сессия: %w", err)
	}
	return u, nil
}

// EndSession is what signing out really is. Asking the browser to forget the
// cookie is politeness; deleting the row is what makes the token dead.
func (s *Store) EndSession(token string) error {
	if _, err := s.db.Exec(`delete from sessions where token_hash = ?`, hashToken(token)); err != nil {
		return fmt.Errorf("шығу: %w", err)
	}
	return nil
}

// hashToken keeps sha256 rather than bcrypt on purpose: the token is 32 random
// bytes of our own, so there is nothing to guess and nothing to pay slowness
// for. bcrypt is for secrets a person invents.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func expiry() string {
	return time.Now().Add(SessionLife).UTC().Format(time.RFC3339)
}
