package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Author-identity (level A): a verified phone number plus a real name. In
// Kazakhstan a SIM is registered to its owner's IIN, so phone verification ties
// an account to a real person without collecting identity documents.
const (
	phoneCodeTTL     = 10 * time.Minute
	phoneMaxAttempts = 5
	phoneCodeDigits  = 6
)

// ErrInvalidPhone is returned when a phone number is empty or malformed.
var ErrInvalidPhone = errors.New("invalid phone number")

// NormalizePhone trims formatting and validates a phone number to E.164-ish
// form (optional leading +, 10–15 digits). Returns ok=false when invalid.
func NormalizePhone(phone string) (normalized string, ok bool) {
	var b strings.Builder
	plus := false
	for i, r := range strings.TrimSpace(phone) {
		switch {
		case r == '+' && i == 0:
			plus = true
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '(' || r == ')':
			// ignore common separators
		default:
			return "", false
		}
	}
	digits := b.String()
	if len(digits) < 10 || len(digits) > 15 {
		return "", false
	}
	if plus {
		return "+" + digits, true
	}
	return digits, true
}

func generateNumericCode(n int) (string, error) {
	const digits = "0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = digits[int(buf[i])%len(digits)]
	}
	return string(buf), nil
}

// StartPhoneVerification mints a one-time code for the phone and sends it by SMS
// (dev-logged when no gateway is configured, outside production). Returns
// ErrInvalidPhone for a malformed number.
func (m *Module) StartPhoneVerification(ctx context.Context, userID uuid.UUID, phone string) error {
	norm, ok := NormalizePhone(phone)
	if !ok {
		return ErrInvalidPhone
	}
	code, err := generateNumericCode(phoneCodeDigits)
	if err != nil {
		return err
	}
	if err := m.store.CreatePhoneCode(ctx, userID, norm, hashToken(code), time.Now().Add(phoneCodeTTL)); err != nil {
		return err
	}
	text := "Shanraq: your verification code is " + code
	return m.deliverSMSOrDevLog(ctx, norm, text, code)
}

// ConfirmPhoneVerification checks the submitted code; on success the account's
// phone is stored and marked verified.
func (m *Module) ConfirmPhoneVerification(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return false, nil
	}
	return m.store.VerifyPhoneCode(ctx, userID, hashToken(code), phoneMaxAttempts)
}

// AuthorIdentity returns the user's real name and whether their phone is
// verified. CanPublish reports whether they may publish articles (both a name
// and a verified phone are required).
func (m *Module) AuthorIdentity(ctx context.Context, userID uuid.UUID) (first, last string, phoneVerified bool) {
	if m.store == nil {
		return "", "", false
	}
	f, l, v, err := m.store.AuthorIdentity(ctx, userID)
	if err != nil {
		return "", "", false
	}
	return f, l, v
}

// MiddleName returns the user's patronymic, or "" if unset or unavailable.
func (m *Module) MiddleName(ctx context.Context, userID uuid.UUID) string {
	if m.store == nil {
		return ""
	}
	middle, err := m.store.MiddleName(ctx, userID)
	if err != nil {
		return ""
	}
	return middle
}

// CanPublish reports whether the user is a verified author (real name + phone).
func (m *Module) CanPublish(ctx context.Context, userID uuid.UUID) bool {
	f, l, v := m.AuthorIdentity(ctx, userID)
	return v && strings.TrimSpace(f) != "" && strings.TrimSpace(l) != ""
}

// SetAuthorName stores the real first/last name shown as the article byline.
func (m *Module) SetAuthorName(ctx context.Context, userID uuid.UUID, first, last, middle string) error {
	first, last, middle = NormalizePersonName(first), NormalizePersonName(last), NormalizePersonName(middle)
	// The same rule as registration: letters only, capitalized.
	if err := ValidatePersonName(first); err != nil {
		return err
	}
	if err := ValidatePersonName(last); err != nil {
		return err
	}
	if err := ValidateOptionalPersonName(middle); err != nil {
		return err
	}
	return m.store.SetName(ctx, userID, first, last, middle)
}

func (m *Module) deliverSMSOrDevLog(ctx context.Context, phone, text, code string) error {
	prod := strings.EqualFold(m.rt.Config.Environment, "production")
	if m.sms != nil {
		if err := m.sms.SendSMS(ctx, phone, text); err == nil {
			return nil
		} else if prod {
			return err
		}
	} else if prod {
		return errors.New("sms delivery is not configured")
	}
	// Non-production: surface the code so the flow can be tested locally.
	m.rt.Logger.Info("phone verification: dev-only code (sms unavailable)", zap.String("code", code))
	return nil
}
