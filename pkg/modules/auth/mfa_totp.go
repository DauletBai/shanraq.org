package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

var (
	// ErrInvalidMFACode is returned when the provided MFA code cannot be validated.
	ErrInvalidMFACode = errors.New("invalid verification code")
)

// TOTPProvider implements MFAProvider using time-based one-time passwords.
type TOTPProvider struct {
	store  *Store
	issuer string
	digits otp.Digits
	period uint
	skew   uint
	logger *zap.Logger
}

// NewTOTPProvider builds a provider that stores secrets in the auth store and verifies TOTP codes.
func NewTOTPProvider(store *Store, issuer string, logger *zap.Logger) *TOTPProvider {
	if issuer == "" {
		issuer = "Shanraq"
	}
	return &TOTPProvider{
		store:  store,
		issuer: issuer,
		digits: otp.DigitsSix,
		period: 30,
		skew:   1,
		logger: logger,
	}
}

func (p *TOTPProvider) Challenge(ctx context.Context, user User) (MFAChallenge, error) {
	record, err := p.store.GetTOTPByUser(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, ErrMFATOTPNotFound) {
			return MFAChallenge{}, err
		}

		secret, uri, genErr := p.generateSecret(user)
		if genErr != nil {
			return MFAChallenge{}, genErr
		}

		record, err = p.store.CreateTOTP(ctx, user.ID, secret)
		if err != nil {
			return MFAChallenge{}, err
		}
		record.Secret = secret
		record.URI = uri
	}

	challenge := MFAChallenge{
		ID:        record.ID.String(),
		Channel:   "totp",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Data: map[string]any{
			"totp_confirmed": record.Confirmed,
		},
	}

	if !record.Confirmed {
		uri := record.URI
		if uri == "" {
			uri = p.formatURI(user, record.Secret)
		}
		challenge.URI = uri
		challenge.Secret = record.Secret
	}

	return challenge, nil
}

func (p *TOTPProvider) Verify(ctx context.Context, challengeID, code string) (MFAResult, error) {
	id, err := uuid.Parse(strings.TrimSpace(challengeID))
	if err != nil {
		return MFAResult{}, ErrInvalidMFACode
	}

	record, err := p.store.GetTOTPByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMFATOTPNotFound) {
			return MFAResult{}, ErrInvalidMFACode
		}
		return MFAResult{}, err
	}

	valid, validateErr := totp.ValidateCustom(
		strings.TrimSpace(code),
		record.Secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    p.period,
			Skew:      p.skew,
			Digits:    p.digits,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	if validateErr != nil {
		return MFAResult{}, fmt.Errorf("validate totp: %w", validateErr)
	}
	if !valid {
		return MFAResult{}, ErrInvalidMFACode
	}

	if !record.Confirmed {
		if err := p.store.MarkTOTPConfirmed(ctx, record.ID); err != nil {
			if p.logger != nil {
				p.logger.Warn("mark totp confirmed", zap.Error(err))
			}
			return MFAResult{}, fmt.Errorf("mark totp confirmed: %w", err)
		}
	}

	user, err := p.store.GetByID(ctx, record.UserID.String())
	if err != nil {
		return MFAResult{}, fmt.Errorf("load user for mfa: %w", err)
	}

	return MFAResult{User: user}, nil
}

func (p *TOTPProvider) generateSecret(user User) (secret string, uri string, err error) {
	account := user.Email
	if account == "" {
		account = user.ID.String()
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      p.issuer,
		AccountName: account,
		Period:      p.period,
		Digits:      p.digits,
	})
	if err != nil {
		return "", "", fmt.Errorf("generate totp secret: %w", err)
	}
	return key.Secret(), key.URL(), nil
}

func (p *TOTPProvider) formatURI(user User, secret string) string {
	account := user.Email
	if account == "" {
		account = user.ID.String()
	}
	query := url.Values{}
	query.Set("secret", secret)
	query.Set("issuer", p.issuer)
	query.Set("period", fmt.Sprintf("%d", p.period))
	query.Set("algorithm", "SHA1")
	query.Set("digits", fmt.Sprintf("%d", p.digits.Length()))

	return fmt.Sprintf("otpauth://totp/%s:%s?%s",
		url.QueryEscape(p.issuer),
		url.QueryEscape(account),
		query.Encode(),
	)
}
