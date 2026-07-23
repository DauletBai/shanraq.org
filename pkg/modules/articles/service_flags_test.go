package articles

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestSiteMaintenanceTimer checks the global switch's countdown + auto-restore:
// a past window makes the site expired, RestoreExpired brings it back and clears
// the timer, and turning a service on drops any timer. Needs SHANRAQ_TEST_DB.
func TestSiteMaintenanceTimer(t *testing.T) {
	dsn := requireTestDB(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	flags := NewServiceFlags(pool)
	if err := flags.Load(ctx); err != nil {
		t.Fatalf("load: %v", err)
	}

	// Take the site down with a window that already elapsed.
	past := time.Now().Add(-time.Minute)
	if err := flags.Set(ctx, SvcSite, svcMaintenance, "", "Скоро вернёмся", "", &past, nil); err != nil {
		t.Fatalf("set site down: %v", err)
	}
	if flags.SiteUp() {
		t.Fatal("site should be down right after Set maintenance")
	}
	if !flags.SiteExpired() {
		t.Fatal("a past window should read as expired")
	}
	changed, err := flags.RestoreExpired(ctx)
	if err != nil || !changed {
		t.Fatalf("RestoreExpired: changed=%v err=%v", changed, err)
	}
	if !flags.SiteUp() {
		t.Fatal("site should be up after auto-restore")
	}
	if flags.SiteFlag().HasTimer() {
		t.Fatal("timer must be cleared after restore")
	}

	// A future window stays down and is not expired.
	future := time.Now().Add(time.Hour)
	if err := flags.Set(ctx, SvcSite, svcMaintenance, "", "x", "", &future, nil); err != nil {
		t.Fatalf("set future: %v", err)
	}
	if flags.SiteExpired() {
		t.Fatal("a future window must not read as expired")
	}
	if got := flags.SiteFlag().UntilUnixMillis(); got == 0 {
		t.Fatal("expected a non-zero until for the countdown")
	}

	// Turning the site back on clears the timer even if a future one was set.
	if err := flags.Set(ctx, SvcSite, svcOn, "", "x", "", &future, nil); err != nil {
		t.Fatalf("set on: %v", err)
	}
	if flags.SiteFlag().HasTimer() {
		t.Fatal("turning on must drop the timer")
	}
}

// TestServiceFlagsIntegration checks that a service flag round-trips through the
// DB and the in-memory cache: the beta seed puts payments in maintenance, a
// write flips a service on and updates its notice, and Available reflects it
// immediately. Skipped unless SHANRAQ_TEST_DB names a test database.
func TestServiceFlagsIntegration(t *testing.T) {
	dsn := requireTestDB(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	flags := NewServiceFlags(pool)
	if err := flags.Load(ctx); err != nil {
		t.Fatalf("load: %v", err)
	}

	// The migration seeds ad_orders in maintenance for the beta.
	if flags.Available(SvcAdOrders) {
		t.Fatalf("ad_orders should start in maintenance (seed), got available")
	}
	if got := flags.Flag(SvcAdOrders); got.Message(LangRU) == "" {
		t.Fatalf("seeded maintenance message should be non-empty")
	}

	// Turn it on with a fresh notice; the cache must reflect it at once.
	if err := flags.Set(ctx, SvcAdOrders, svcOn, "kz", "ru", "en", nil, nil); err != nil {
		t.Fatalf("set on: %v", err)
	}
	if !flags.Available(SvcAdOrders) {
		t.Fatalf("ad_orders should be available after Set on")
	}

	// Back into maintenance with a custom message; Available flips off again.
	if err := flags.Set(ctx, SvcAdOrders, svcMaintenance, "kz2", "ru2", "en2", nil, nil); err != nil {
		t.Fatalf("set maintenance: %v", err)
	}
	if flags.Available(SvcAdOrders) {
		t.Fatalf("ad_orders should be unavailable in maintenance")
	}
	if got := flags.Flag(SvcAdOrders).Message(LangEN); got != "en2" {
		t.Fatalf("message not updated: got %q", got)
	}

	// An unknown service or bad status must be rejected.
	if err := flags.Set(ctx, "nope", svcOn, "", "", "", nil, nil); err == nil {
		t.Fatalf("expected error for unknown service")
	}
	if err := flags.Set(ctx, SvcAdOrders, "banana", "", "", "", nil, nil); err == nil {
		t.Fatalf("expected error for bad status")
	}

	// invite_only is allowed on a free function but rejected on a paid one.
	if err := flags.Set(ctx, SvcRegistration, svcInviteOnly, "", "", "", nil, nil); err != nil {
		t.Fatalf("registration invite_only should be allowed: %v", err)
	}
	if flags.Flag(SvcRegistration).Status != svcInviteOnly {
		t.Fatalf("registration should be invite_only")
	}
	if !flags.Flag(SvcRegistration).Invitable() {
		t.Fatalf("registration must be invitable")
	}
	if err := flags.Set(ctx, SvcAdOrders, svcInviteOnly, "", "", "", nil, nil); err == nil {
		t.Fatalf("invite_only must be rejected on a paid service")
	}
	// Reset registration to on so re-runs start clean.
	if err := flags.Set(ctx, SvcRegistration, svcOn, "", "", "", nil, nil); err != nil {
		t.Fatalf("reset registration: %v", err)
	}

	// Restore the beta seed state so re-runs start clean.
	if err := flags.Set(ctx, SvcAdOrders, svcMaintenance,
		"Жарнама төлемі уақытша қолжетімсіз.",
		"Оплата рекламы временно недоступна.",
		"Paid advertising is temporarily unavailable.", nil, nil); err != nil {
		t.Fatalf("restore seed: %v", err)
	}
}
