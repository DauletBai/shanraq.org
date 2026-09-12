package payments

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Payment is one intent to charge, in our own records.
type Payment struct {
	ID        uuid.UUID
	Kind      string // what is being paid for, e.g. "ad_order"
	TargetID  uuid.UUID
	Amount    int64
	Currency  string
	Provider  string
	Status    string // pending | paid | failed | expired
	ExpiresAt time.Time
}

// Settlement is what the seller does with its own record when a payment of its
// kind settles or expires. Both hooks run inside the payment's transaction, so
// the order and the payment move together or not at all: there is no moment
// when money is taken and the thing paid for is not active.
//
// A missing hook is not an error. Some kinds are recorded before the thing they
// pay for has any state to flip.
type Settlement struct {
	Paid    func(ctx context.Context, tx pgx.Tx, target uuid.UUID) error
	Expired func(ctx context.Context, tx pgx.Tx, target uuid.UUID) error
}

// Store persists payments and settles them against what they pay for.
type Store struct {
	db *pgxpool.Pool

	mu    sync.RWMutex
	kinds map[string]Settlement
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db, kinds: map[string]Settlement{}}
}

// Handle registers the settlement hooks for one kind of payment. Modules call
// it during Init, before anything can be bought.
func (s *Store) Handle(kind string, h Settlement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kinds[kind] = h
}

func (s *Store) settlement(kind string) Settlement {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.kinds[kind]
}

// Create records a pending payment for a target, held for holdMinutes. The hold
// is why an unpaid checkout cannot occupy an ad slot forever.
func (s *Store) Create(ctx context.Context, kind string, target uuid.UUID, amount int64, holdMinutes int) (Payment, error) {
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
func (s *Store) MarkPaid(ctx context.Context, id uuid.UUID) (settled bool, err error) {
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
	if h := s.settlement(kind).Paid; h != nil {
		if err := h(ctx, tx, target); err != nil {
			return false, fmt.Errorf("settle %s: %w", kind, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

// ExpirePending releases holds that were never paid: the payment becomes
// 'expired' and whatever it was holding is released, so the slot frees up.
// Meant to run on a timer; returns how many were expired.
func (s *Store) ExpirePending(ctx context.Context) (int, error) {
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
		if h := s.settlement(rf.kind).Expired; h != nil {
			if err := h(ctx, tx, rf.target); err != nil {
				return 0, fmt.Errorf("release %s: %w", rf.kind, err)
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(refs), nil
}

// Get loads a payment by id.
func (s *Store) Get(ctx context.Context, id uuid.UUID) (Payment, error) {
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
