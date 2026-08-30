package auth

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A password-reset link is a temporary key to somebody's account, and every
// property below is a way that key can be turned into a permanent one. They run
// against a real database because the guarantees are the table's.

type authFixture struct {
	pool  *pgxpool.Pool
	store *Store
	ctx   context.Context
}

func newAuthFixture(t *testing.T) *authFixture {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the auth store tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return &authFixture{pool: pool, store: NewStore(pool), ctx: ctx}
}

func (f *authFixture) user(t *testing.T) uuid.UUID {
	t.Helper()
	id := uuid.New()
	if _, err := f.pool.Exec(f.ctx,
		`INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'old-hash','user')`,
		id, "pr-"+id.String()+"@t.test"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_password_resets WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_refresh_tokens WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_users WHERE id=$1`, id)
	})
	return id
}

// The link is looked up by a hash, so the token itself is not in the table: a
// database dump must not be a stack of working reset links.
func TestAResetLinkIsStoredOnlyAsItsHash(t *testing.T) {
	f := newAuthFixture(t)
	user := f.user(t)
	const hash = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	pr, err := f.store.CreatePasswordReset(f.ctx, user, hash, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if pr.TokenHash != hash {
		t.Errorf("stored %q, want the hash it was given", pr.TokenHash)
	}
	if pr.UsedAt != nil {
		t.Error("a fresh reset is already marked used")
	}

	got, err := f.store.GetPasswordReset(f.ctx, hash)
	if err != nil || got.ID != pr.ID {
		t.Fatalf("looking the reset up by its hash failed: %v", err)
	}
	// A hash nobody issued is not a reset, and saying so is a distinct error the
	// handler can answer without leaking whether an account exists.
	if _, err := f.store.GetPasswordReset(f.ctx, "no-such-hash"); !errors.Is(err, ErrPasswordResetNotFound) {
		t.Errorf("an unknown hash gave %v, want ErrPasswordResetNotFound", err)
	}
}

// Single use is the whole point. A link that still works after the password has
// been changed is a link an attacker can use later.
func TestAResetLinkIsSpentOnceItIsUsed(t *testing.T) {
	f := newAuthFixture(t)
	user := f.user(t)
	const hash = "aaaa789abcdef0123456789abcdef0123456789abcdef0123456789abcdefaaa"

	pr, err := f.store.CreatePasswordReset(f.ctx, user, hash, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := f.store.MarkPasswordResetUsed(f.ctx, pr.ID); err != nil {
		t.Fatalf("mark used: %v", err)
	}

	spent, err := f.store.GetPasswordReset(f.ctx, hash)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if spent.UsedAt == nil {
		t.Fatal("a used reset link is not marked used, so it would be accepted again")
	}
	firstUse := *spent.UsedAt

	// Marking it again must not move the timestamp: the row records when the
	// link was spent, and a second attempt is not a second spending.
	if err := f.store.MarkPasswordResetUsed(f.ctx, pr.ID); err != nil {
		t.Fatalf("second mark: %v", err)
	}
	again, _ := f.store.GetPasswordReset(f.ctx, hash)
	if !again.UsedAt.Equal(firstUse) {
		t.Errorf("the used-at moved from %v to %v on a second attempt", firstUse, *again.UsedAt)
	}
}

// An expired link is stored with its expiry so the handler can refuse it; the
// row does not disappear on its own.
func TestAnExpiredResetKeepsItsExpiry(t *testing.T) {
	f := newAuthFixture(t)
	user := f.user(t)
	const hash = "bbbb789abcdef0123456789abcdef0123456789abcdef0123456789abcdefbbb"
	past := time.Now().Add(-time.Hour)

	if _, err := f.store.CreatePasswordReset(f.ctx, user, hash, past); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := f.store.GetPasswordReset(f.ctx, hash)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !got.ExpiresAt.Before(time.Now()) {
		t.Errorf("expiry came back as %v, which is not in the past", got.ExpiresAt)
	}
}

// Changing a password signs out everyone else holding a token for the account.
// That is the point of a reset after a compromise, and it is done by moving
// auth_version, which every session is checked against.
func TestChangingAPasswordInvalidatesExistingSessions(t *testing.T) {
	f := newAuthFixture(t)
	user := f.user(t)

	var before int
	if err := f.pool.QueryRow(f.ctx, `SELECT auth_version FROM auth_users WHERE id=$1`, user).Scan(&before); err != nil {
		t.Fatalf("read version: %v", err)
	}
	if err := f.store.UpdatePassword(f.ctx, user, "new-hash"); err != nil {
		t.Fatalf("update password: %v", err)
	}

	var after int
	var hash string
	var mustReset bool
	if err := f.pool.QueryRow(f.ctx,
		`SELECT auth_version, password_hash, password_reset_required FROM auth_users WHERE id=$1`,
		user).Scan(&after, &hash, &mustReset); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if after <= before {
		t.Errorf("auth_version went %d to %d; sessions held elsewhere would survive the reset", before, after)
	}
	if hash != "new-hash" {
		t.Errorf("password_hash = %q, want the new one", hash)
	}
	if mustReset {
		t.Error("the account is still flagged as needing a reset after one was done")
	}
}
