package articles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrPaymentsDisabled is returned when no payment provider is configured — the
// default state, exactly like the AI and SMTP modules ship disabled.
var ErrPaymentsDisabled = errors.New("payments are not configured")

// Charge is what a provider hands back for an amount: a QR to display and/or a
// checkout URL to open, plus the provider's own reference.
type Charge struct {
	QR  string // QR payload or data: image, for display
	URL string // hosted checkout / redirect URL
	Ref string // provider's transaction reference
}

// WebhookResult is the outcome a provider extracts from an authenticated
// callback: which of OUR payments it concerns and whether it succeeded.
type WebhookResult struct {
	PaymentID uuid.UUID
	Paid      bool
	Failed    bool
}

// PaymentProvider is the seam every acquirer/aggregator plugs into. A concrete
// adapter (ioka, Freedom Pay, Halyk ePay, Kaspi) implements it once the owner
// picks a provider and supplies sandbox keys; the order flow never changes.
type PaymentProvider interface {
	Name() string
	// CreateCharge asks the provider to open a charge for a payment we created,
	// returning a QR/URL to present to the buyer.
	CreateCharge(ctx context.Context, p Payment) (Charge, error)
	// HandleWebhook authenticates an incoming provider callback (signature,
	// source) and reports which payment it settles. It MUST reject an unsigned
	// or forged request — a payment is only ever marked paid from a verified
	// callback, never from the client.
	HandleWebhook(r *http.Request) (WebhookResult, error)
}

// disabledProvider is the default: no money moves, nothing to misconfigure, no
// credentials in the repo. Every method fails closed.
type disabledProvider struct{}

func (disabledProvider) Name() string { return "" }
func (disabledProvider) CreateCharge(context.Context, Payment) (Charge, error) {
	return Charge{}, ErrPaymentsDisabled
}
func (disabledProvider) HandleWebhook(*http.Request) (WebhookResult, error) {
	return WebhookResult{}, ErrPaymentsDisabled
}

// Payment is one intent to charge, in our own records.
type Payment struct {
	ID        uuid.UUID
	Kind      string // ad_order | listing_promo
	TargetID  uuid.UUID
	Amount    int64
	Currency  string
	Provider  string
	Status    string // pending | paid | failed | expired
	ExpiresAt time.Time
}

// PaymentStore persists payments and settles them against what they pay for.
type PaymentStore struct{ db *pgxpool.Pool }

func NewPaymentStore(db *pgxpool.Pool) *PaymentStore { return &PaymentStore{db: db} }

// Create records a pending payment for a target, held for holdMinutes. The hold
// is why an unpaid checkout cannot occupy an ad slot forever.
func (s *PaymentStore) Create(ctx context.Context, kind string, target uuid.UUID, amount int64, holdMinutes int) (Payment, error) {
	if holdMinutes <= 0 {
		holdMinutes = 30
	}
	p := Payment{Kind: kind, TargetID: target, Amount: amount, Currency: "KZT", Status: "pending"}
	err := s.db.QueryRow(ctx, `
		INSERT INTO payments (kind, target_id, amount, expires_at)
		VALUES ($1, $2, $3, NOW() + (($4::int) || ' minutes')::interval)
		RETURNING id, expires_at`, kind, target, amount, holdMinutes).Scan(&p.ID, &p.ExpiresAt)
	if err != nil {
		return Payment{}, fmt.Errorf("create payment: %w", err)
	}
	return p, nil
}

// MarkPaid settles a payment and activates what it pays for, in one
// transaction, idempotently: a webhook that arrives twice flips the status only
// once, and the second call is a no-op. Returns whether this call was the one
// that settled it.
func (s *PaymentStore) MarkPaid(ctx context.Context, id uuid.UUID) (settled bool, err error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var kind string
	var target uuid.UUID
	err = tx.QueryRow(ctx, `
		UPDATE payments SET status = 'paid', paid_at = NOW()
		 WHERE id = $1 AND status = 'pending'
		RETURNING kind, target_id`, id).Scan(&kind, &target)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil // unknown, or already settled — idempotent
	}
	if err != nil {
		return false, fmt.Errorf("mark paid: %w", err)
	}

	switch kind {
	case "ad_order":
		if _, err := tx.Exec(ctx, `
			UPDATE ad_orders SET status = 'active'
			 WHERE id = $1 AND status = 'pending_payment'`, target); err != nil {
			return false, fmt.Errorf("activate ad order: %w", err)
		}
	case "listing_promo":
		// Listing promotions are not routed through payments yet; the hook is
		// here so the same settlement path serves them when they are.
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

// ExpirePending releases holds that were never paid: the payment becomes
// 'expired' and its ad order is cancelled, so the slot frees up. Meant to run
// on a timer; returns how many were expired.
func (s *PaymentStore) ExpirePending(ctx context.Context) (int, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, `
		UPDATE payments SET status = 'expired'
		 WHERE status = 'pending' AND expires_at < NOW()
		RETURNING kind, target_id`)
	if err != nil {
		return 0, fmt.Errorf("expire pending: %w", err)
	}
	type ref struct {
		kind   string
		target uuid.UUID
	}
	var refs []ref
	for rows.Next() {
		var rf ref
		if err := rows.Scan(&rf.kind, &rf.target); err != nil {
			rows.Close()
			return 0, err
		}
		refs = append(refs, rf)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, rf := range refs {
		if rf.kind == "ad_order" {
			if _, err := tx.Exec(ctx, `
				UPDATE ad_orders SET status = 'cancelled'
				 WHERE id = $1 AND status = 'pending_payment'`, rf.target); err != nil {
				return 0, fmt.Errorf("cancel expired ad order: %w", err)
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(refs), nil
}

// Get loads a payment by id.
func (s *PaymentStore) Get(ctx context.Context, id uuid.UUID) (Payment, error) {
	var p Payment
	err := s.db.QueryRow(ctx, `
		SELECT id, kind, target_id, amount, currency, COALESCE(provider,''), status, COALESCE(expires_at, NOW())
		FROM payments WHERE id = $1`, id).Scan(
		&p.ID, &p.Kind, &p.TargetID, &p.Amount, &p.Currency, &p.Provider, &p.Status, &p.ExpiresAt)
	if err != nil {
		return Payment{}, fmt.Errorf("get payment: %w", err)
	}
	return p, nil
}
