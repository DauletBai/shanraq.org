package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func revocationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run revocation tests")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// seedUser inserts an account and returns it as the token issuer sees it.
func seedUser(t *testing.T, pool *pgxpool.Pool, role string) User {
	t.Helper()
	ctx := context.Background()
	u := User{ID: uuid.New(), Email: "rev-" + uuid.NewString()[:8] + "@t.test", Role: role, Roles: []string{role}}
	if _, err := pool.Exec(ctx,
		`INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x',$3)`,
		u.ID, u.Email, role); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, u.ID) })
	if err := pool.QueryRow(ctx, `SELECT auth_version FROM auth_users WHERE id=$1`, u.ID).
		Scan(&u.AuthVersion); err != nil {
		t.Fatalf("read auth_version: %v", err)
	}
	return u
}

// The regression the audit asked for: an administrator's token must stop
// working the moment the administrator stops being one. Roles live in the JWT,
// so before this the demoted account kept its powers until the token expired —
// two hours in production, with no session row anywhere to delete.
func TestPrivilegedTokenDiesWithTheRole(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	tokens := NewTokenService("test-token-secret-that-is-long-enough-1234567890", time.Hour)
	mod := &Module{tokens: tokens, store: store}

	admin := seedUser(t, pool, "admin")
	token, err := tokens.Generate(admin)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	reached := 0
	guarded := mod.RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached++
		w.WriteHeader(http.StatusTeapot)
	}))
	call := func() int {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		guarded.ServeHTTP(rec, req)
		return rec.Code
	}

	if got := call(); got != http.StatusTeapot {
		t.Fatalf("a sitting administrator was refused: %d", got)
	}

	// Demote through the real path, with the same token still in hand.
	if err := store.SetRoleByID(context.Background(), admin.ID, "user"); err != nil {
		t.Fatalf("demote: %v", err)
	}
	if got := call(); got != http.StatusUnauthorized {
		t.Fatalf("the demoted administrator's token still worked: %d", got)
	}
	if reached != 1 {
		t.Fatalf("the guarded handler ran %d times, want 1", reached)
	}
}

// Deleting the account has to end its tokens too — there is no row left to
// compare against, and "cannot confirm" must mean "refuse".
func TestPrivilegedTokenDiesWithTheAccount(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	tokens := NewTokenService("test-token-secret-that-is-long-enough-1234567890", time.Hour)
	mod := &Module{tokens: tokens, store: store}

	// A second administrator, so removing the first is not refused as the last.
	seedUser(t, pool, "admin")
	victim := seedUser(t, pool, "admin")
	token, err := tokens.Generate(victim)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	guarded := mod.RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	call := func() int {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		guarded.ServeHTTP(rec, req)
		return rec.Code
	}
	if got := call(); got != http.StatusTeapot {
		t.Fatalf("a sitting administrator was refused: %d", got)
	}
	if err := store.DeleteUserByAdmin(context.Background(), victim.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := call(); got != http.StatusUnauthorized {
		t.Fatalf("a deleted account's token still worked: %d", got)
	}
}

// Reading the site is still authorised by the signature alone: the database
// lookup happens only where a role is being relied upon.
func TestUnprivilegedRequestsSkipTheLookup(t *testing.T) {
	tokens := NewTokenService("test-token-secret-that-is-long-enough-1234567890", time.Hour)
	mod := &Module{tokens: tokens} // deliberately no store
	tok, err := tokens.Generate(User{ID: uuid.New(), Email: "a@b.c", Role: "user"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	guarded := mod.RequireRoles()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	guarded.ServeHTTP(rec, req)
	if rec.Code != http.StatusTeapot {
		t.Fatalf("a request needing no role was made to pay for a lookup: %d", rec.Code)
	}
}

// parkOtherAdmins demotes every administrator except the named ones, restoring
// them when the test ends.
func parkOtherAdmins(t *testing.T, pool *pgxpool.Pool, keep ...uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	rows, err := pool.Query(ctx,
		`SELECT id, role FROM auth_users WHERE role = ANY($1) AND id <> ALL($2)`, rootRoles, keep)
	if err != nil {
		t.Fatalf("list admins: %v", err)
	}
	type parked struct {
		id   uuid.UUID
		role string
	}
	var others []parked
	for rows.Next() {
		var pk parked
		if err := rows.Scan(&pk.id, &pk.role); err != nil {
			rows.Close()
			t.Fatalf("scan admin: %v", err)
		}
		others = append(others, pk)
	}
	rows.Close()
	for _, o := range others {
		if _, err := pool.Exec(ctx, `UPDATE auth_users SET role = 'user' WHERE id = $1`, o.id); err != nil {
			t.Fatalf("park admin: %v", err)
		}
	}
	t.Cleanup(func() {
		for _, o := range others {
			_, _ = pool.Exec(ctx, `UPDATE auth_users SET role = $2 WHERE id = $1`, o.id, o.role)
		}
	})
}

// Two administrators, two simultaneous demotions: exactly one must succeed.
//
// The count and the demotion used to be separate statements outside any
// transaction, so both callers read "two administrators", both concluded it was
// safe, and the site was left with none. Run repeatedly, because a race that
// only sometimes loses is still a race.
func TestLastAdminSurvivesConcurrentDemotion(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	for attempt := 0; attempt < 8; attempt++ {
		a := seedUser(t, pool, "admin")
		b := seedUser(t, pool, "admin")

		// Any other administrator — the bootstrap account, a leftover — would
		// make both demotions legitimately safe and prove nothing. Park them for
		// the duration with raw SQL, which deliberately bypasses the very guard
		// under test, and put them back afterwards.
		parkOtherAdmins(t, pool, a.ID, b.ID)

		start := make(chan struct{})
		errs := make(chan error, 2)
		for _, id := range []uuid.UUID{a.ID, b.ID} {
			go func(id uuid.UUID) {
				<-start
				errs <- store.SetRoleByID(ctx, id, "user")
			}(id)
		}
		close(start)
		first, second := <-errs, <-errs

		refused := 0
		for _, err := range []error{first, second} {
			switch {
			case err == nil:
			case errors.Is(err, ErrLastAdmin):
				refused++
			default:
				t.Fatalf("unexpected error: %v", err)
			}
		}
		if refused != 1 {
			t.Fatalf("attempt %d: %d of 2 demotions were refused, want exactly 1", attempt, refused)
		}

		var left int
		if err := pool.QueryRow(ctx,
			`SELECT count(*) FROM auth_users WHERE role = ANY($1)`, rootRoles).Scan(&left); err != nil {
			t.Fatalf("recount: %v", err)
		}
		if left != 1 {
			t.Fatalf("attempt %d left %d administrators, want 1", attempt, left)
		}

		_, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id = ANY($1)`, []uuid.UUID{a.ID, b.ID})
	}
}

// The hole the audit found: the console refuses to delete the last
// administrator, but the administrator could delete *themselves* from
// /studio/delete, which went straight to an unconditional DELETE. One
// administrator, one right-to-erasure click, and nobody left who can administer
// the site — a state no console action can undo.
func TestLastAdminCannotDeleteThemselves(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	a := seedUser(t, pool, "admin")
	parkOtherAdmins(t, pool, a.ID)

	err := store.DeleteAccount(ctx, a.ID)
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("self-delete of the last administrator returned %v, want ErrLastAdmin", err)
	}

	var left int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM auth_users WHERE role = ANY($1)`, rootRoles).Scan(&left); err != nil {
		t.Fatalf("recount: %v", err)
	}
	if left != 1 {
		t.Fatalf("left %d administrators, want 1", left)
	}

	// An ordinary account still deletes: the guard must not turn erasure off for
	// everyone to protect one case.
	u := seedUser(t, pool, "user")
	if err := store.DeleteAccount(ctx, u.ID); err != nil {
		t.Fatalf("ordinary account could not delete itself: %v", err)
	}
}

// Self-deletion and demotion are two doors onto the same invariant, so the race
// has to be tested across them, not only within each. Every pairing of the two
// must leave exactly one administrator standing.
func TestLastAdminSurvivesConcurrentRemoval(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	remove := map[string]func(uuid.UUID) error{
		"delete": func(id uuid.UUID) error { return store.DeleteAccount(ctx, id) },
		"demote": func(id uuid.UUID) error { return store.SetRoleByID(ctx, id, "user") },
	}
	pairs := [][2]string{{"delete", "delete"}, {"delete", "demote"}, {"demote", "delete"}}

	for _, pair := range pairs {
		for attempt := 0; attempt < 6; attempt++ {
			a := seedUser(t, pool, "admin")
			b := seedUser(t, pool, "admin")
			parkOtherAdmins(t, pool, a.ID, b.ID)

			start := make(chan struct{})
			errs := make(chan error, 2)
			for i, id := range []uuid.UUID{a.ID, b.ID} {
				act := remove[pair[i]]
				go func(id uuid.UUID) {
					<-start
					errs <- act(id)
				}(id)
			}
			close(start)

			refused := 0
			for i := 0; i < 2; i++ {
				switch err := <-errs; {
				case err == nil:
				case errors.Is(err, ErrLastAdmin):
					refused++
				default:
					t.Fatalf("%v attempt %d: unexpected error: %v", pair, attempt, err)
				}
			}
			if refused != 1 {
				t.Fatalf("%v attempt %d: %d of 2 removals refused, want exactly 1", pair, attempt, refused)
			}

			var left int
			if err := pool.QueryRow(ctx,
				`SELECT count(*) FROM auth_users WHERE role = ANY($1)`, rootRoles).Scan(&left); err != nil {
				t.Fatalf("recount: %v", err)
			}
			if left != 1 {
				t.Fatalf("%v attempt %d left %d administrators, want 1", pair, attempt, left)
			}

			_, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id = ANY($1)`, []uuid.UUID{a.ID, b.ID})
		}
	}
}

// adminctl is the third door onto the same change. A guard the web console
// enforces and the command line does not is a guard with a way around it.
func TestSetPrimaryRoleKeepsTheLastAdmin(t *testing.T) {
	pool := revocationPool(t)
	store := NewStore(pool)
	ctx := context.Background()

	a := seedUser(t, pool, "admin")
	parkOtherAdmins(t, pool, a.ID)

	found, err := store.SetPrimaryRole(ctx, a.Email, "user")
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("adminctl demotion of the last administrator returned (%v, %v), want ErrLastAdmin", found, err)
	}

	// Promotion is never the dangerous direction, and it must still bump
	// auth_version so a token issued under the old role stops being honoured.
	u := seedUser(t, pool, "user")
	if _, err := store.SetPrimaryRole(ctx, u.Email, "admin"); err != nil {
		t.Fatalf("promote: %v", err)
	}
	var role string
	var version int
	if err := pool.QueryRow(ctx,
		`SELECT role, auth_version FROM auth_users WHERE id = $1`, u.ID).Scan(&role, &version); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if role != "admin" {
		t.Fatalf("role is %q, want admin", role)
	}
	if version <= u.AuthVersion {
		t.Fatalf("auth_version stayed at %d: tokens issued under the old role still pass", version)
	}
}
