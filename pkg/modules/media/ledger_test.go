package media

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
)

func ledgerPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run media ledger tests")
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

func seedOwner(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	id := uuid.New()
	if _, err := pool.Exec(ctx,
		`INSERT INTO auth_users (id, email, password_hash, role) VALUES ($1,$2,'x','user')`,
		id, "media-"+uuid.NewString()[:8]+"@t.test"); err != nil {
		t.Fatalf("seed owner: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, id) })
	return id
}

// storeKey files an object as though it had been uploaded ageHours ago.
func storeKey(t *testing.T, pool *pgxpool.Pool, l *ledger, owner uuid.UUID, key string, size int64, ageHours int) {
	t.Helper()
	ctx := context.Background()
	if err := l.record(ctx, key, size, "image/jpeg", owner); err != nil {
		t.Fatalf("record %s: %v", key, err)
	}
	if _, err := pool.Exec(ctx,
		`UPDATE media_objects SET created_at = now() - make_interval(hours => $2) WHERE key = $1`,
		key, ageHours); err != nil {
		t.Fatalf("backdate %s: %v", key, err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM media_objects WHERE key=$1`, key) })
}

// The 25 MB request cap only ever bounded one request. Nothing bounded how many
// requests, so an account with a valid session could fill the volume the site
// runs on — and on a single VPS that is the site going down, not a warning.
func TestQuotaBoundsWhatOneAccountCanHold(t *testing.T) {
	pool := ledgerPool(t)
	l := &ledger{db: pool}
	ctx := context.Background()
	owner := seedOwner(t, pool)

	m := &Module{ledger: l, logger: zap.NewNop(), cfg: config.MediaConfig{QuotaBytes: 1000}}

	first := "aa/" + uuid.NewString() + ".jpg"
	if err := m.reserve(ctx, owner, first, 600); err != nil {
		t.Fatalf("first upload inside the quota was refused: %v", err)
	}
	storeKey(t, pool, l, owner, first, 600, 0)

	second := "bb/" + uuid.NewString() + ".jpg"
	if err := m.reserve(ctx, owner, second, 600); !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("the upload that would exceed the quota returned %v, want ErrQuotaExceeded", err)
	}

	// Storing what you already hold costs nothing: a form resubmitted after a
	// validation error must not be charged twice for the same photo.
	if err := m.reserve(ctx, owner, first, 600); err != nil {
		t.Fatalf("re-storing an already-held object was refused: %v", err)
	}

	// Another account has its own quota, and shares the object rather than a
	// second copy on disk.
	other := seedOwner(t, pool)
	if err := m.reserve(ctx, other, first, 600); err != nil {
		t.Fatalf("a second account was charged for someone else's file: %v", err)
	}
}

// The store cap is the floor under the service: every account can be inside its
// own quota and the volume still full.
func TestStoreCapRefusesNewBytesButNotSharedOnes(t *testing.T) {
	pool := ledgerPool(t)
	l := &ledger{db: pool}
	ctx := context.Background()
	owner := seedOwner(t, pool)

	existing := "cc/" + uuid.NewString() + ".jpg"
	storeKey(t, pool, l, owner, existing, 500, 0)

	total, err := l.total(ctx)
	if err != nil {
		t.Fatalf("total: %v", err)
	}
	m := &Module{ledger: l, logger: zap.NewNop(), cfg: config.MediaConfig{MaxTotalBytes: total}}

	fresh := "dd/" + uuid.NewString() + ".jpg"
	if err := m.reserve(ctx, seedOwner(t, pool), fresh, 1); !errors.Is(err, ErrStoreFull) {
		t.Fatalf("a new object on a full store returned %v, want ErrStoreFull", err)
	}
	// An object already on disk adds no bytes, so a second claim on it is fine
	// even at the cap.
	if err := m.reserve(ctx, seedOwner(t, pool), existing, 500); err != nil {
		t.Fatalf("claiming an object already on disk was refused at the cap: %v", err)
	}
}

// An upload becomes reachable only when the form it belongs to is saved. The
// ones that never were are invisible on every screen and permanent on disk.
func TestSweepCollectsOnlyWhatNothingRefers(t *testing.T) {
	pool := ledgerPool(t)
	l := &ledger{db: pool}
	ctx := context.Background()
	owner := seedOwner(t, pool)

	dir := t.TempDir()
	store, err := NewFSStore(dir, "media")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	m := &Module{
		store:  store,
		ledger: l,
		logger: zap.NewNop(),
		cfg:    config.MediaConfig{PublicPrefix: "/media", OrphanGraceHours: 1},
	}

	avatar := "ee/" + uuid.NewString() + ".jpg"
	inBody := "ff/" + uuid.NewString() + ".jpg"
	orphan := "07/" + uuid.NewString() + ".jpg"
	for _, k := range []string{avatar, inBody, orphan} {
		if err := store.Put(ctx, k, []byte("x"), "image/jpeg"); err != nil {
			t.Fatalf("put %s: %v", k, err)
		}
		storeKey(t, pool, l, owner, k, 1, 48)
	}

	// Referenced as a whole-value URL column…
	if _, err := pool.Exec(ctx, `UPDATE auth_users SET avatar_url = $2 WHERE id = $1`,
		owner, "/media/"+avatar); err != nil {
		t.Fatalf("set avatar: %v", err)
	}
	// …and referenced only from inside free text, which the URL columns cannot
	// see and a sweep that ignored prose would happily delete.
	pageKey := "sweep-" + uuid.NewString()[:8]
	if _, err := pool.Exec(ctx,
		`INSERT INTO content_pages (page_key, lang, title, body_md) VALUES ($1,'ru','t',$2)`,
		pageKey, "текст ![фото](/media/"+inBody+") дальше"); err != nil {
		t.Fatalf("seed page: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM content_pages WHERE page_key=$1`, pageKey) })

	m.sweepOrphans(ctx)

	gone := func(key string) bool {
		_, err := os.Stat(dir + "/" + key)
		return os.IsNotExist(err)
	}
	if !gone(orphan) {
		t.Error("the unreferenced upload survived the sweep")
	}
	if gone(avatar) {
		t.Error("the sweep deleted a file an avatar points at")
	}
	if gone(inBody) {
		t.Error("the sweep deleted a file referenced only from page text")
	}

	var left int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM media_objects WHERE key = $1`, orphan).Scan(&left); err != nil {
		t.Fatalf("recount: %v", err)
	}
	if left != 0 {
		t.Error("the collected file is still in the ledger")
	}
}

// Deleting the account releases its claim, which is what lets the sweep reach
// files nobody is left to own.
func TestDeletingAnAccountReleasesItsFiles(t *testing.T) {
	pool := ledgerPool(t)
	l := &ledger{db: pool}
	ctx := context.Background()
	owner := seedOwner(t, pool)

	key := "1a/" + uuid.NewString() + ".jpg"
	storeKey(t, pool, l, owner, key, 42, 0)

	if used, err := l.usage(ctx, owner); err != nil || used != 42 {
		t.Fatalf("usage = %d, %v; want 42", used, err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM auth_users WHERE id = $1`, owner); err != nil {
		t.Fatalf("delete owner: %v", err)
	}
	var claims int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM media_owners WHERE key = $1`, key).Scan(&claims); err != nil {
		t.Fatalf("count claims: %v", err)
	}
	if claims != 0 {
		t.Fatalf("%d claims survived the account, want 0", claims)
	}
}
