package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// SessionCookieName is the HttpOnly cookie holding a signed access token for
// browser (non-API) flows such as the author cabinet.
const SessionCookieName = "shanraq_session"

// ErrInvalidCredentials is returned when email/password verification fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrInvalidEmail is returned when a registration email is empty or malformed.
var ErrInvalidEmail = errors.New("invalid email address")

// ErrWeakPassword is returned when a registration password is too short.
var ErrWeakPassword = errors.New("password too short")

// minPasswordLen / maxPasswordLen bound the password. The max mirrors bcrypt's
// 72-byte input limit so an over-long password gets a clean validation error
// instead of a bcrypt error surfacing as a 500.
const (
	minPasswordLen = 8
	maxPasswordLen = 72
)

// ConsentDocument and ConsentVersion identify the legal documents a user
// accepts at registration. Bump ConsentVersion when the Terms or Privacy
// Policy change materially so re-consent can be required later.
const (
	ConsentDocument = "terms_privacy"
	ConsentVersion  = "2026-07-18"
)

// AuthorConsentDocument/Version identify the author's one-time acknowledgment of
// the documents AND tariffs, required before publishing. Bump the version when
// tariffs or author terms change materially to require re-consent.
const (
	AuthorConsentDocument = "author_terms_tariffs"
	AuthorConsentVersion  = "2026-07-19"
)

// HasAuthorConsent reports whether the user has acknowledged the author
// documents + tariffs at least once. A nil store (tests) treats it as granted.
func (m *Module) HasAuthorConsent(ctx context.Context, userID uuid.UUID) (bool, error) {
	if m.store == nil {
		return true, nil
	}
	return m.store.HasConsent(ctx, userID, AuthorConsentDocument, AuthorConsentVersion)
}

// RecordAuthorConsent appends the author's acknowledgment of the documents and
// tariffs (append-only proof), required once before publishing.
func (m *Module) RecordAuthorConsent(ctx context.Context, r *http.Request, userID uuid.UUID, source string) error {
	if m.store == nil {
		return nil
	}
	return m.store.InsertConsent(ctx, userID, AuthorConsentDocument, AuthorConsentVersion, source, clientIdentifier(r))
}

// RecordConsent appends a consent record for a newly registered user. source is
// "web" or "api"; the client IP is taken from the request. Failure is returned
// so the caller can decide, but registration should not hard-fail on it.
func (m *Module) RecordConsent(ctx context.Context, r *http.Request, userID uuid.UUID, source string) error {
	if m.store == nil {
		return nil
	}
	return m.store.InsertConsent(ctx, userID, ConsentDocument, ConsentVersion, source, clientIdentifier(r))
}

// emailVerificationTTL is how long a verification link stays valid.
const emailVerificationTTL = 48 * time.Hour

// IssueEmailVerification mints a verification token for the user and emails the
// confirmation link. Best-effort: a delivery/store failure is returned but the
// caller should not fail registration on it (the user can request a resend).
func (m *Module) IssueEmailVerification(ctx context.Context, userID uuid.UUID, email string) error {
	token, err := generateSecureToken(refreshTokenSize)
	if err != nil {
		return err
	}
	if err := m.store.CreateEmailVerification(ctx, userID, hashToken(token), time.Now().Add(emailVerificationTTL)); err != nil {
		return err
	}
	link := fmt.Sprintf("%s/auth/verify?token=%s", m.rt.Config.PublicBase(), token)
	return m.sendVerificationEmail(ctx, email, link)
}

// VerifyEmail consumes a verification token and marks the email verified.
func (m *Module) VerifyEmail(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrEmailVerificationInvalid
	}
	_, err := m.store.ConsumeEmailVerification(ctx, hashToken(token))
	return err
}

// IsEmailVerified reports whether the user has confirmed their email. On error
// it returns false so gated actions fail closed.
func (m *Module) IsEmailVerified(ctx context.Context, userID uuid.UUID) bool {
	if m.store == nil {
		return false
	}
	ok, err := m.store.IsEmailVerified(ctx, userID)
	return err == nil && ok
}

// IsPhoneVerified reports whether the user has confirmed their phone. On error
// it returns false so gated actions fail closed.
func (m *Module) IsPhoneVerified(ctx context.Context, userID uuid.UUID) bool {
	if m.store == nil {
		return false
	}
	ok, err := m.store.IsPhoneVerified(ctx, userID)
	return err == nil && ok
}

func (m *Module) sendVerificationEmail(ctx context.Context, to, link string) error {
	subject := "Confirm your email"
	body := fmt.Sprintf("Welcome to Shanraq. Please confirm your email by opening the link below:\n\n%s\n\nIf you did not create an account, you can ignore this message.", link)
	return m.deliverOrDevLink(ctx, to, subject, body, link, "email verification")
}

// deliverOrDevLink sends a transactional email. If no mailer is wired or the
// mailer can't deliver (e.g. SMTP unconfigured), it fails loudly in production
// but, outside production, logs the link so the flow can be tested locally. It
// never logs the token/link in production.
func (m *Module) deliverOrDevLink(ctx context.Context, to, subject, body, link, purpose string) error {
	prod := strings.EqualFold(m.rt.Config.Environment, "production")
	if m.mailer != nil {
		if err := m.mailer.Send(ctx, to, subject, body); err == nil {
			return nil
		} else if prod {
			return err
		}
	} else if prod {
		return errors.New("email delivery is not configured")
	}
	// Non-production: surface the link so developers can complete the flow.
	m.rt.Logger.Info(purpose+": dev-only link (mailer unavailable)", zap.String("link", link))
	return nil
}

// NormalizeEmail lowercases and trims an email; ok is false when it is empty or
// not a syntactically valid address.
func NormalizeEmail(email string) (normalized string, ok bool) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(email) > 254 {
		return "", false
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", false
	}
	// Require a dotted domain (reject "user@localhost"-style typos on a public site).
	at := strings.LastIndexByte(email, '@')
	domain := email[at+1:]
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", false
	}
	return email, true
}

// isSecureRequest reports whether the request arrived over TLS (directly or via
// a trusted proxy), so the Secure cookie flag can be set appropriately.
func isSecureRequest(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

// SetSessionCookie writes the session cookie for a browser client.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
	if ttl <= 0 {
		ttl = time.Hour
	}
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
	})
}

// ClearSessionCookie removes the session cookie (sign-out).
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (m *Module) claimsFromCookie(r *http.Request) (*Claims, bool) {
	if m.tokens == nil {
		return nil, false
	}
	c, err := r.Cookie(SessionCookieName)
	if err != nil || c.Value == "" {
		return nil, false
	}
	claims, err := m.tokens.Parse(c.Value)
	if err != nil {
		return nil, false
	}
	return claims, true
}

// LoadSession injects claims from the session cookie when the request has no
// bearer-token claims yet. It never blocks the request.
func (m *Module) LoadSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := ClaimsFromContext(r.Context()); ok {
			next.ServeHTTP(w, r)
			return
		}
		if claims, ok := m.claimsFromCookie(r); ok {
			r = r.WithContext(ContextWithClaims(r.Context(), claims))
		}
		next.ServeHTTP(w, r)
	})
}

// RequireSession guards browser pages: it resolves claims from context or the
// session cookie and redirects unauthenticated visitors to loginPath.
func (m *Module) RequireSession(loginPath string, roles ...string) func(http.Handler) http.Handler {
	normalized := sanitizeRoles(roles)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				claims, ok = m.claimsFromCookie(r)
			}
			if !ok {
				http.Redirect(w, r, loginPath, http.StatusSeeOther)
				return
			}
			if len(normalized) > 0 && !claims.HasAnyRole(normalized...) {
				http.Redirect(w, r, loginPath, http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithClaims(r.Context(), claims)))
		})
	}
}

// LoginPassword verifies credentials and returns the user with a freshly signed
// access token suitable for a session cookie.
func (m *Module) LoginPassword(ctx context.Context, email, password string) (User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return User{}, "", ErrInvalidCredentials
	}
	user, err := m.store.FindByEmail(ctx, email)
	if err != nil {
		return User{}, "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return User{}, "", ErrInvalidCredentials
	}
	token, err := m.tokens.Generate(user)
	if err != nil {
		return User{}, "", err
	}
	return user, token, nil
}

// RegisterPassword creates a new user account and returns a signed access token.
func (m *Module) RegisterPassword(ctx context.Context, email, password string) (User, string, error) {
	email, ok := NormalizeEmail(email)
	if !ok {
		return User{}, "", ErrInvalidEmail
	}
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return User{}, "", ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, "", err
	}
	user, err := m.store.CreateUser(ctx, email, string(hash), defaultRoleName)
	if err != nil {
		return User{}, "", err
	}
	token, err := m.tokens.Generate(user)
	if err != nil {
		return User{}, "", err
	}
	return user, token, nil
}

// SessionTTL exposes the access-token lifetime for cookie expiry.
func (m *Module) SessionTTL() time.Duration {
	if m.tokens == nil {
		return time.Hour
	}
	return m.tokens.TTL()
}
