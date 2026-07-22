package articles

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestReferralLoopIntegration exercises invite → capture → qualify → credit →
// spend against a real database. Skipped unless SHANRAQ_TEST_DB names a test DB.
func TestReferralLoopIntegration(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to run the referral integration test")
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
	rs := NewReferralStore(pool)

	referrer, referred := uuid.New(), uuid.New()
	for _, id := range []uuid.UUID{referrer, referred} {
		if _, err := pool.Exec(ctx, `INSERT INTO auth_users (id,email,password_hash,role) VALUES ($1,$2,'x','user')`,
			id, "ref-"+id.String()+"@t.test"); err != nil {
			t.Fatalf("user: %v", err)
		}
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM auth_users WHERE id IN ($1,$2)`, referrer, referred) })

	code, err := rs.EnsureCode(ctx, referrer)
	if err != nil || code == "" {
		t.Fatalf("code: %v %q", err, code)
	}
	// Idempotent: same code on second call.
	if again, _ := rs.EnsureCode(ctx, referrer); again != code {
		t.Fatalf("code not stable: %q vs %q", code, again)
	}
	// Resolve.
	if got, ok := rs.ReferrerByCode(ctx, code); !ok || got != referrer {
		t.Fatalf("resolve: %v %v", got, ok)
	}
	// Self-referral is a no-op.
	if err := rs.RecordReferral(ctx, referrer, referrer); err != nil {
		t.Fatalf("self referral errored: %v", err)
	}
	// Capture.
	if err := rs.RecordReferral(ctx, referrer, referred); err != nil {
		t.Fatalf("record: %v", err)
	}
	// No credit before qualification.
	if bal, _ := rs.Balance(ctx, referrer); bal != 0 {
		t.Fatalf("credit before qualify: %d", bal)
	}
	// Qualify grants the reward.
	granted, err := rs.Qualify(ctx, referred)
	if err != nil || !granted {
		t.Fatalf("qualify: %v granted=%v", err, granted)
	}
	if bal, _ := rs.Balance(ctx, referrer); bal != referralRewardDays {
		t.Fatalf("credit after qualify: want %d got %d", referralRewardDays, bal)
	}
	// Qualifying again grants nothing (one reward per referred user).
	if granted, _ := rs.Qualify(ctx, referred); granted {
		t.Fatal("second qualify must not grant again")
	}
	if bal, _ := rs.Balance(ctx, referrer); bal != referralRewardDays {
		t.Fatalf("credit double-granted: %d", bal)
	}

	// Spend on a listing the referrer owns.
	listing := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO listings (id,author_id,title,status) VALUES ($1,$2,$3,'published')`,
		listing, referrer, "ref-itest-"+listing.String()[:8]); err != nil {
		t.Fatalf("listing: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM listings WHERE id=$1`, listing) })
	if err := rs.SpendCredit(ctx, referrer, listing, referralRewardDays); err != nil {
		t.Fatalf("spend: %v", err)
	}
	if bal, _ := rs.Balance(ctx, referrer); bal != 0 {
		t.Fatalf("credit after spend: %d", bal)
	}
	// Overspend is refused.
	if err := rs.SpendCredit(ctx, referrer, listing, referralRewardDays); err == nil {
		t.Fatal("overspend must be refused")
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM promo_credit_ledger WHERE user_id=$1`, referrer) })
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM referrals WHERE referrer_id=$1`, referrer) })
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM referral_codes WHERE user_id=$1`, referrer) })
}
