package articles

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestPaymentsIntegration covers the provider-agnostic core: creating a hold,
// settling it once (idempotently) and activating the ad order, and expiring an
// abandoned hold so the slot frees. Skipped unless SHANRAQ_TEST_DB names a
// test database.
func TestPaymentsIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the payments integration test")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatalf("SHANRAQ_TEST_DB must name a test database; refusing %q", dsn)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	ps := NewPaymentStore(pool)

	// A minimal advertiser + ad order in pending_payment.
	owner := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO auth_users (id,email,password_hash,role) VALUES ($1,$2,'x','user')`,
		owner, "pay-"+owner.String()+"@t.test"); err != nil {
		t.Fatalf("user: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id=$1`, owner) })
	var advID uuid.UUID
	if err := pool.QueryRow(ctx, `INSERT INTO advertisers (owner_id,company_name,contact_name,contact_phone)
		VALUES ($1,'C','N','+7') RETURNING id`, owner).Scan(&advID); err != nil {
		t.Fatalf("advertiser: %v", err)
	}
	var orderID uuid.UUID
	if err := pool.QueryRow(ctx, `INSERT INTO ad_orders (advertiser_id,title,target_url,surfaces,format,price,status)
		VALUES ($1,'T','https://x',ARRAY['home'],'horizontal',180000,'pending_payment') RETURNING id`, advID).Scan(&orderID); err != nil {
		t.Fatalf("order: %v", err)
	}

	// Create the hold.
	p, err := ps.Create(ctx, "ad_order", orderID, 180000, 30)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Settle once → order active; settle again → no-op (idempotent).
	settled, err := ps.MarkPaid(ctx, p.ID)
	if err != nil || !settled {
		t.Fatalf("mark paid: %v settled=%v", err, settled)
	}
	again, _ := ps.MarkPaid(ctx, p.ID)
	if again {
		t.Fatal("second MarkPaid must be a no-op")
	}
	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM ad_orders WHERE id=$1`, orderID).Scan(&status); err != nil {
		t.Fatalf("read order: %v", err)
	}
	if status != "active" {
		t.Fatalf("order should be active, got %q", status)
	}

	// A second order whose hold has already lapsed must expire and cancel.
	var order2 uuid.UUID
	if err := pool.QueryRow(ctx, `INSERT INTO ad_orders (advertiser_id,title,target_url,surfaces,format,price,status)
		VALUES ($1,'T2','https://x',ARRAY['home'],'horizontal',180000,'pending_payment') RETURNING id`, advID).Scan(&order2); err != nil {
		t.Fatalf("order2: %v", err)
	}
	var pay2 uuid.UUID
	if err := pool.QueryRow(ctx, `INSERT INTO payments (kind,target_id,amount,expires_at)
		VALUES ('ad_order',$1,180000, NOW() - interval '1 minute') RETURNING id`, order2).Scan(&pay2); err != nil {
		t.Fatalf("pay2: %v", err)
	}
	n, err := ps.ExpirePending(ctx)
	if err != nil || n < 1 {
		t.Fatalf("expire: %v n=%d", err, n)
	}
	if err := pool.QueryRow(ctx, `SELECT status FROM ad_orders WHERE id=$1`, order2).Scan(&status); err != nil {
		t.Fatalf("read order2: %v", err)
	}
	if status != "cancelled" {
		t.Fatalf("expired order should be cancelled, got %q", status)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM payments WHERE target_id IN ($1,$2)`, orderID, order2)
		_, _ = pool.Exec(ctx, `DELETE FROM ad_orders WHERE advertiser_id=$1`, advID)
		_, _ = pool.Exec(ctx, `DELETE FROM advertisers WHERE id=$1`, advID)
	})
}
