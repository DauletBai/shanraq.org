package articles

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// referralRewardDays is the promotion credit a referrer earns when someone they
// invited posts a real listing. Three days of "top" placement — inventory the
// platform already sells, so the marginal cost is zero.
const referralRewardDays = 3

// ReferralStore holds invite codes, the referral graph, and the promotion-day
// credit ledger.
type ReferralStore struct{ db *pgxpool.Pool }

func NewReferralStore(db *pgxpool.Pool) *ReferralStore { return &ReferralStore{db: db} }

// codeAlphabet excludes look-alike characters (0/O, 1/I/L) so a code read from
// a screen and typed by hand does not turn into the wrong one.
const codeAlphabet = "abcdefghijkmnpqrstuvwxyz23456789"

func newCode() (string, error) {
	b := make([]byte, 7)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, len(b))
	for i, x := range b {
		out[i] = codeAlphabet[int(x)%len(codeAlphabet)]
	}
	return string(out), nil
}

// EnsureCode returns the user's invite code, generating it once on first use.
func (s *ReferralStore) EnsureCode(ctx context.Context, userID uuid.UUID) (string, error) {
	var code string
	err := s.db.QueryRow(ctx, `SELECT code FROM referral_codes WHERE user_id = $1`, userID).Scan(&code)
	if err == nil {
		return code, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("referral code lookup: %w", err)
	}
	// Generate, retrying on the rare collision.
	for attempt := 0; attempt < 5; attempt++ {
		c, gerr := newCode()
		if gerr != nil {
			return "", gerr
		}
		_, ierr := s.db.Exec(ctx, `INSERT INTO referral_codes (user_id, code) VALUES ($1, $2)
			ON CONFLICT (user_id) DO NOTHING`, userID, c)
		if ierr != nil {
			// Unique violation on the code column — try another.
			continue
		}
		// Re-read: the row may have been inserted by a concurrent request.
		if err := s.db.QueryRow(ctx, `SELECT code FROM referral_codes WHERE user_id = $1`, userID).Scan(&code); err == nil {
			return code, nil
		}
	}
	return "", fmt.Errorf("could not allocate a referral code")
}

// ReferrerByCode resolves an invite code to the user who owns it.
func (s *ReferralStore) ReferrerByCode(ctx context.Context, code string) (uuid.UUID, bool) {
	var id uuid.UUID
	if err := s.db.QueryRow(ctx, `SELECT user_id FROM referral_codes WHERE code = $1`, code).Scan(&id); err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// RecordReferral links an invited user to their referrer at registration. It is
// a no-op on self-referral or if the invited user is already referred, so it is
// safe to call unconditionally.
func (s *ReferralStore) RecordReferral(ctx context.Context, referrer, referred uuid.UUID) error {
	if referrer == referred {
		return nil
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO referrals (referrer_id, referred_id) VALUES ($1, $2)
		ON CONFLICT (referred_id) DO NOTHING`, referrer, referred)
	if err != nil {
		return fmt.Errorf("record referral: %w", err)
	}
	return nil
}

// IsReferred reports whether a user joined through someone's invite link. It is
// the "invited" signal for invite_only mode: a beta tester who registered via an
// admin's invite is referred, so they may create content while the public cannot.
func (s *ReferralStore) IsReferred(ctx context.Context, userID uuid.UUID) bool {
	var exists bool
	if err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM referrals WHERE referred_id = $1)`, userID).Scan(&exists); err != nil {
		return false
	}
	return exists
}

// Qualify is called when a user does the rewardable action (posts a real
// listing). If that user was referred and the referral is still pending, it
// flips to qualified and grants the referrer their credit — both in one
// transaction, so a grant can never exist without its qualified referral.
// Returns whether a reward was granted.
func (s *ReferralStore) Qualify(ctx context.Context, referred uuid.UUID) (bool, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var refID, referrer uuid.UUID
	err = tx.QueryRow(ctx, `
		UPDATE referrals SET status = 'qualified', qualified_at = NOW()
		 WHERE referred_id = $1 AND status = 'pending'
		RETURNING id, referrer_id`, referred).Scan(&refID, &referrer)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil // not referred, or already qualified
	}
	if err != nil {
		return false, fmt.Errorf("qualify referral: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO promo_credit_ledger (user_id, delta_days, reason, ref_id)
		VALUES ($1, $2, 'referral', $3)`, referrer, referralRewardDays, refID); err != nil {
		return false, fmt.Errorf("grant credit: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

// Balance returns a user's unused promotion-day credit.
func (s *ReferralStore) Balance(ctx context.Context, userID uuid.UUID) (int, error) {
	var days int
	if err := s.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(delta_days), 0) FROM promo_credit_ledger WHERE user_id = $1`,
		userID).Scan(&days); err != nil {
		return 0, fmt.Errorf("credit balance: %w", err)
	}
	return days, nil
}

// ReferralStats is the referrer's own view: how many they invited and how many
// turned into a reward.
type ReferralStats struct {
	Code      string
	Invited   int
	Qualified int
	Credit    int
}

func (s *ReferralStore) Stats(ctx context.Context, userID uuid.UUID) (ReferralStats, error) {
	var st ReferralStats
	code, err := s.EnsureCode(ctx, userID)
	if err != nil {
		return st, err
	}
	st.Code = code
	if err := s.db.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'qualified')
		  FROM referrals WHERE referrer_id = $1`, userID).Scan(&st.Invited, &st.Qualified); err != nil {
		return st, fmt.Errorf("referral stats: %w", err)
	}
	if st.Credit, err = s.Balance(ctx, userID); err != nil {
		return st, err
	}
	return st, nil
}

// SpendCredit debits promotion days and applies them to a listing the caller
// owns, atomically: the balance is re-checked inside the transaction, so two
// concurrent requests cannot spend the same credit twice. The listing's
// promoted_until is extended by the spent days.
func (s *ReferralStore) SpendCredit(ctx context.Context, userID, listingID uuid.UUID, days int) error {
	if days <= 0 {
		return fmt.Errorf("days must be positive")
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lock this user's ledger rows so the balance cannot change under us.
	// FOR UPDATE is not allowed with an aggregate, so lock the rows first, then
	// sum them within the same transaction.
	if _, err := tx.Exec(ctx,
		`SELECT 1 FROM promo_credit_ledger WHERE user_id = $1 FOR UPDATE`, userID); err != nil {
		return fmt.Errorf("spend: lock: %w", err)
	}
	var balance int
	if err := tx.QueryRow(ctx,
		`SELECT COALESCE(SUM(delta_days),0) FROM promo_credit_ledger WHERE user_id = $1`,
		userID).Scan(&balance); err != nil {
		return fmt.Errorf("spend: balance: %w", err)
	}
	if balance < days {
		return ErrInsufficientCredit
	}
	// The listing must belong to the caller; touch its promotion window.
	tag, err := tx.Exec(ctx, `
		UPDATE listings
		   SET promoted_until = GREATEST(COALESCE(promoted_until, NOW()), NOW()) + (($3::int) || ' days')::interval
		 WHERE id = $1 AND author_id = $2`, listingID, userID, days)
	if err != nil {
		return fmt.Errorf("spend: promote: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO promo_credit_ledger (user_id, delta_days, reason, ref_id)
		VALUES ($1, $2, 'spend_promote', $3)`, userID, -days, listingID); err != nil {
		return fmt.Errorf("spend: debit: %w", err)
	}
	return tx.Commit(ctx)
}

// ErrInsufficientCredit is returned when a user tries to spend credit they do
// not have.
var ErrInsufficientCredit = errors.New("insufficient promotion credit")
