package articles

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
	if err := flags.Set(ctx, SvcAdOrders, svcOn, "kz", "ru", "en", nil); err != nil {
		t.Fatalf("set on: %v", err)
	}
	if !flags.Available(SvcAdOrders) {
		t.Fatalf("ad_orders should be available after Set on")
	}

	// Back into maintenance with a custom message; Available flips off again.
	if err := flags.Set(ctx, SvcAdOrders, svcMaintenance, "kz2", "ru2", "en2", nil); err != nil {
		t.Fatalf("set maintenance: %v", err)
	}
	if flags.Available(SvcAdOrders) {
		t.Fatalf("ad_orders should be unavailable in maintenance")
	}
	if got := flags.Flag(SvcAdOrders).Message(LangEN); got != "en2" {
		t.Fatalf("message not updated: got %q", got)
	}

	// An unknown service or bad status must be rejected.
	if err := flags.Set(ctx, "nope", svcOn, "", "", "", nil); err == nil {
		t.Fatalf("expected error for unknown service")
	}
	if err := flags.Set(ctx, SvcAdOrders, "banana", "", "", "", nil); err == nil {
		t.Fatalf("expected error for bad status")
	}

	// Restore the beta seed state so re-runs start clean.
	if err := flags.Set(ctx, SvcAdOrders, svcMaintenance,
		"Жарнама төлемі уақытша қолжетімсіз.",
		"Оплата рекламы временно недоступна.",
		"Paid advertising is temporarily unavailable.", nil); err != nil {
		t.Fatalf("restore seed: %v", err)
	}
}
