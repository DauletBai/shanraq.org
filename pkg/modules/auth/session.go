package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SessionCookieName is the HttpOnly cookie holding a signed access token for
// browser (non-API) flows such as the author cabinet.
const SessionCookieName = "shanraq_session"

// ErrInvalidCredentials is returned when email/password verification fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

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
	email = strings.TrimSpace(strings.ToLower(email))
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
