package auth

import (
	"context"
	"time"
)

// MFAProvider allows the auth module to hand off multi-factor challenges.
type MFAProvider interface {
	// Challenge should initiate a second-factor flow for the given user (e.g. TOTP, SMS, WebAuthn)
	// and return metadata required by the client to continue the flow.
	Challenge(ctx context.Context, user User) (MFAChallenge, error)
	// Verify validates the submitted code or assertion for the referenced challenge and, on success,
	// returns the user that completed the multi-factor step.
	Verify(ctx context.Context, challengeID, code string) (MFAResult, error)
}

// MFAChallenge describes an outstanding multi-factor verification prompt.
type MFAChallenge struct {
	ID        string
	Channel   string
	ExpiresAt time.Time
	Secret    string
	URI       string
	Data      map[string]any
}

// MFAResult captures the user record issued by the MFA provider upon a successful verification.
type MFAResult struct {
	User User
}
