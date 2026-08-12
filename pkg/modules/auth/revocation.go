package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrTokenRevoked is returned when a token is well-formed and unexpired but no
// longer speaks for the account it names.
var ErrTokenRevoked = errors.New("credentials are no longer valid; sign in again")

// BumpAuthVersion retires every token an account currently holds. Call it from
// anything that changes what the account is allowed to do, or who controls it:
// a role change, a deletion, a password reset.
func (s *Store) BumpAuthVersion(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx,
		`UPDATE auth_users SET auth_version = auth_version + 1, updated_at = now() WHERE id = $1`, id)
	return err
}

// authVersionOf reads the account's current revocation counter. A missing row
// yields 0, which matches no issued token — a deleted account is refused by the
// same comparison that refuses a demoted one, with no special case.
func (s *Store) authVersionOf(ctx context.Context, id uuid.UUID) (int, bool) {
	if s == nil || s.db == nil {
		return 0, false
	}
	var v int
	if err := s.db.QueryRow(ctx, `SELECT auth_version FROM auth_users WHERE id = $1`, id).Scan(&v); err != nil {
		// Absent row, or a database we cannot reach: either way this request
		// must not proceed on the strength of a role we cannot confirm.
		return 0, false
	}
	return v, true
}

// tokenStillValid reports whether claims still speak for their account.
//
// Only privileged requests pay for this. Reading the site is authorised by the
// signature alone, as before; the lookup happens where the answer can actually
// matter — when a role is being relied upon.
func (m *Module) tokenStillValid(ctx context.Context, claims *Claims) bool {
	if claims == nil || m.store == nil {
		return false
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return false
	}
	current, ok := m.store.authVersionOf(ctx, id)
	if !ok {
		return false
	}
	// Tokens minted before this column existed carry no version. Treating them
	// as valid against version 1 keeps everyone signed in across the deploy;
	// the first bump on an account ends that grace for it.
	if claims.AuthVersion == 0 {
		return current == 1
	}
	return claims.AuthVersion == current
}
