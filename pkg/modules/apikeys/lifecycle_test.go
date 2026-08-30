package apikeys

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"shanraq.org/pkg/modules/auth"
)

// An API key is a password that is never typed, so the properties worth pinning
// down are the ones a mock cannot show: what is actually written to the table,
// and what the table refuses afterwards.

type keyFixture struct {
	pool  *pgxpool.Pool
	store *Store
	ctx   context.Context
}

func newKeyFixture(t *testing.T) *keyFixture {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the api-key lifecycle tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return &keyFixture{pool: pool, store: NewStore(pool), ctx: ctx}
}

func (f *keyFixture) user(t *testing.T) uuid.UUID {
	t.Helper()
	id := uuid.New()
	if _, err := f.pool.Exec(f.ctx,
		`INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`,
		id, "ak-"+id.String()+"@t.test"); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_api_keys WHERE user_id=$1`, id)
		_, _ = f.pool.Exec(f.ctx, `DELETE FROM auth_users WHERE id=$1`, id)
	})
	return id
}

// The secret is shown once and must not survive anywhere it could be read back.
// If this ever fails, a database dump becomes a set of working credentials.
func TestTheSecretIsNeverStored(t *testing.T) {
	f := newKeyFixture(t)
	owner := f.user(t)

	secret, key, err := f.store.Create(f.ctx, owner, "backup script")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	var stored, prefix string
	if err := f.pool.QueryRow(f.ctx,
		`SELECT key_hash, prefix FROM auth_api_keys WHERE id=$1`, key.ID).Scan(&stored, &prefix); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if stored == secret {
		t.Fatal("the key itself was written to the table")
	}
	if strings.Contains(stored, secret) || strings.Contains(secret, stored) {
		t.Fatal("the stored value and the secret share their text")
	}
	if stored != hashKey(secret) {
		t.Error("the stored value is not the hash of the secret, so no key could ever validate")
	}
	// The prefix is deliberate: enough to recognise a key in a list, never
	// enough to reconstruct it.
	if !strings.HasPrefix(secret, prefix) || len(prefix) >= len(secret) {
		t.Errorf("prefix %q is not a short opening of the secret", prefix)
	}
}

// Revocation is the only thing standing between a leaked key and an open door.
func TestARevokedKeyStopsWorking(t *testing.T) {
	f := newKeyFixture(t)
	owner := f.user(t)
	secret, key, err := f.store.Create(f.ctx, owner, "temporary")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if user, _, err := f.store.Validate(f.ctx, secret); err != nil || user.ID != owner {
		t.Fatalf("a fresh key did not validate: user %+v, err %v", user, err)
	}
	if err := f.store.Revoke(f.ctx, owner, key.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, _, err := f.store.Validate(f.ctx, secret); err == nil {
		t.Fatal("a revoked key still opens the door")
	}
	// Revoking again changes nothing and says so, rather than reporting success
	// for work it did not do.
	if err := f.store.Revoke(f.ctx, owner, key.ID); err == nil {
		t.Error("revoking an already revoked key reported success")
	}
}

// A key belongs to its owner: nobody else may retire it, and the refusal must
// not depend on knowing whether the id exists.
func TestOnlyTheOwnerCanRevoke(t *testing.T) {
	f := newKeyFixture(t)
	owner, stranger := f.user(t), f.user(t)
	secret, key, err := f.store.Create(f.ctx, owner, "mine")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := f.store.Revoke(f.ctx, stranger, key.ID); err == nil {
		t.Fatal("a stranger revoked someone else's key")
	}
	if _, _, err := f.store.Validate(f.ctx, secret); err != nil {
		t.Errorf("the owner's key stopped working after a stranger's attempt: %v", err)
	}
	if err := f.store.Revoke(f.ctx, stranger, uuid.New()); err == nil {
		t.Error("revoking a key that does not exist reported success")
	}
}

// Listing keys is how an owner audits them, and it must not hand back the
// secrets it is auditing.
func TestListingKeysReturnsNoSecrets(t *testing.T) {
	f := newKeyFixture(t)
	owner, other := f.user(t), f.user(t)
	secret, _, err := f.store.Create(f.ctx, owner, "one")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := f.store.Create(f.ctx, other, "someone else's"); err != nil {
		t.Fatalf("create: %v", err)
	}

	keys, err := f.store.List(f.ctx, owner)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("owner sees %d keys, want only their own", len(keys))
	}
	if keys[0].Label != "one" {
		t.Errorf("label = %q, want %q", keys[0].Label, "one")
	}
	if strings.Contains(keys[0].Prefix, secret[len(keys[0].Prefix):]) {
		t.Error("the listing carried more of the secret than its prefix")
	}
}

// Nothing is not a key, and neither is a plausible-looking string.
func TestValidateRefusesEmptyAndUnknownTokens(t *testing.T) {
	f := newKeyFixture(t)
	for _, token := range []string{"", "   ", "sk_not_a_real_key", "sk_"} {
		if _, _, err := f.store.Validate(f.ctx, token); err == nil {
			t.Errorf("token %q was accepted", token)
		}
	}
}

// Two keys issued in a row must not resemble each other: the entropy is the
// whole security of the scheme.
func TestEachKeyIsItsOwn(t *testing.T) {
	f := newKeyFixture(t)
	owner := f.user(t)
	seen := map[string]bool{}
	for i := 0; i < 8; i++ {
		secret, key, err := f.store.Create(f.ctx, owner, "n")
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if seen[secret] {
			t.Fatal("the same secret was issued twice")
		}
		seen[secret] = true
		if !strings.HasPrefix(secret, "sk_") {
			t.Errorf("secret %q does not carry the scheme prefix", secret[:6])
		}
		_ = key
	}
}

// The gate itself. A request carrying a key is either admitted as that user or
// refused; a request carrying none passes through unauthenticated, which is how
// the same routes stay open to a browser session.
func TestTheGateAdmitsRefusesAndStandsAside(t *testing.T) {
	f := newKeyFixture(t)
	owner := f.user(t)
	secret, key, err := f.store.Create(f.ctx, owner, "gate")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	mod := &Module{store: f.store}

	var sawUser string
	var sawKey APIKey
	var sawKeyOK bool
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := auth.ClaimsFromContext(r.Context()); ok {
			sawUser = c.UserID
		}
		sawKey, sawKeyOK = APIKeyFromContext(r.Context())
		w.WriteHeader(http.StatusTeapot)
	})
	gate := mod.RequireAPIKey()(inner)

	run := func(set func(*http.Request)) *httptest.ResponseRecorder {
		sawUser, sawKey, sawKeyOK = "", APIKey{}, false
		req := httptest.NewRequest(http.MethodGet, "/api/thing", nil)
		set(req)
		rec := httptest.NewRecorder()
		gate.ServeHTTP(rec, req)
		return rec
	}

	// Both spellings the extractor accepts.
	for _, set := range []func(*http.Request){
		func(r *http.Request) { r.Header.Set("X-API-Key", secret) },
		func(r *http.Request) { r.Header.Set("Authorization", "ApiKey "+secret) },
		func(r *http.Request) { r.Header.Set("Authorization", "apikey "+secret) },
	} {
		rec := run(set)
		if rec.Code != http.StatusTeapot {
			t.Errorf("a valid key was not let through: status %d", rec.Code)
		}
		if sawUser != owner.String() {
			t.Errorf("the request arrived as %q, want the key's owner", sawUser)
		}
		if !sawKeyOK || sawKey.ID != key.ID {
			t.Error("the handler could not tell which key admitted it")
		}
	}

	// No key at all is not a refusal: the browser session uses these routes too.
	if rec := run(func(*http.Request) {}); rec.Code != http.StatusTeapot {
		t.Errorf("a request with no key was blocked: status %d", rec.Code)
	} else if sawUser != "" {
		t.Error("a request with no key arrived carrying an identity")
	}

	// A key that is not one is refused outright rather than passed through
	// unauthenticated, which would silently downgrade a caller who meant to
	// authenticate.
	if rec := run(func(r *http.Request) { r.Header.Set("X-API-Key", "sk_nonsense") }); rec.Code != http.StatusUnauthorized {
		t.Errorf("a bad key gave status %d, want 401", rec.Code)
	}

	// And once revoked, the same request stops being admitted.
	if err := f.store.Revoke(f.ctx, owner, key.ID); err != nil {
		t.Fatal(err)
	}
	if rec := run(func(r *http.Request) { r.Header.Set("X-API-Key", secret) }); rec.Code != http.StatusUnauthorized {
		t.Errorf("a revoked key gave status %d, want 401", rec.Code)
	}
}
