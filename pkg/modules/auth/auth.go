package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/respond"
	"shanraq.org/pkg/transport/validate"
)

const (
	refreshTokenTTL        = 30 * 24 * time.Hour
	passwordResetTTL       = time.Hour
	refreshTokenSize       = 48
	maxActiveRefreshTokens = 5
)

type contextKey string

const claimsContextKey contextKey = "shanraq/auth.claims"

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type Option func(*Module)

func WithMailer(mailer Mailer) Option {
	return func(m *Module) {
		m.mailer = mailer
	}
}

func WithTOTP(issuer string) Option {
	return func(m *Module) {
		m.requireTOTP = true
		m.totpIssuer = issuer
	}
}

func WithRateLimiter(limiter RateLimiter) Option {
	return func(m *Module) {
		m.rateLimiter = limiter
	}
}

func WithMFAProvider(provider MFAProvider) Option {
	return func(m *Module) {
		m.mfaProvider = provider
	}
}

var errTooManyRequests = errors.New("too many attempts, please try again later")

// Module wires auth routes plus helper services (store + tokens).
type Module struct {
	rt          *shanraq.Runtime
	store       *Store
	tokens      *TokenService
	views       *template.Template
	validator   *validate.Validator
	mailer      Mailer
	rateLimiter RateLimiter
	mfaProvider MFAProvider
	requireTOTP bool
	totpIssuer  string
}

//go:embed templates/*.html
var viewFiles embed.FS

type templateData struct {
	Token                string
	Message              string
	Error                string
	FrameworkName        string
	FrameworkDescription string
	BrandLogoPath        string
	Year                 int
}

// New returns a ready-to-register authentication module.
func New(opts ...Option) *Module {
	m := &Module{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Module) Name() string {
	return "auth"
}

// Init prepares dependencies shared across handlers.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	if strings.EqualFold(rt.Config.Environment, "production") && isWeakSecret(rt.Config.Auth.TokenSecret) {
		return fmt.Errorf("auth token secret must be overridden in production")
	}
	m.tokens = NewTokenService(rt.Config.Auth.TokenSecret, rt.Config.Auth.TokenTTL)
	tmpl, err := template.ParseFS(viewFiles, "templates/*.html")
	if err != nil {
		return fmt.Errorf("parse auth templates: %w", err)
	}
	m.views = tmpl
	m.validator = validate.New()
	if m.rateLimiter == nil {
		m.rateLimiter = newMemoryRateLimiter(defaultRateLimitRules())
	}
	if m.mfaProvider == nil && m.requireTOTP {
		issuer := m.totpIssuer
		if issuer == "" {
			issuer = "Shanraq"
		}
		m.mfaProvider = NewTOTPProvider(m.store, issuer, m.rt.Logger)
	}
	return nil
}

// Routes registers auth endpoints into the shared router.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", m.handleSignup)
		r.Post("/signin", m.handleSignin)
		r.Post("/refresh", m.handleRefresh)
		r.Post("/signout", m.handleSignout)
		r.Get("/profile", m.handleProfile)
		r.Get("/password/reset", m.renderPasswordResetPage)
		r.Post("/password/reset", m.handlePasswordResetRequest)
		r.Get("/password/confirm", m.renderPasswordConfirmPage)
		r.Post("/password/confirm", m.handlePasswordResetConfirm)
		r.Post("/mfa/verify", m.handleMFAVerify)
	})
}

func (m *Module) handleSignup(w http.ResponseWriter, r *http.Request) {
	if !m.enforceRateLimit(r, "signup", true) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	var req signupRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if !m.enforceRateLimit(r, "signup", false, req.Email) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		respond.Validation(w, fields)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	user, err := m.store.CreateUser(ctx, req.Email, string(hash), "user")
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			respond.Error(w, http.StatusConflict, err)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, refreshToken, err := m.issueTokenPair(ctx, user.ID, user)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	m.rt.Logger.Info("user signed up", zap.String("email", user.Email))
	m.writeTokenResponse(w, http.StatusCreated, user, accessToken, refreshToken)
}

func (m *Module) handleSignin(w http.ResponseWriter, r *http.Request) {
	if !m.enforceRateLimit(r, "signin", true) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	var req signinRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if !m.enforceRateLimit(r, "signin", false, req.Email) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		respond.Validation(w, fields)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := m.store.FindByEmail(ctx, req.Email)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		respond.Error(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	accessToken, refreshToken, err := m.issueTokenPair(ctx, user.ID, user)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	if m.mfaProvider != nil {
		challenge, challengeErr := m.mfaProvider.Challenge(ctx, user)
		if challengeErr != nil {
			if m.rt != nil {
				m.rt.Logger.Error("mfa challenge", zap.Error(challengeErr))
			}
			respond.Error(w, http.StatusInternalServerError, errors.New("multi-factor challenge failed"))
			return
		}
		resp := map[string]any{
			"mfa_required": true,
			"challenge_id": challenge.ID,
			"expires_at":   challenge.ExpiresAt,
			"channel":      challenge.Channel,
		}
		if challenge.Secret != "" {
			resp["totp_secret"] = challenge.Secret
		}
		if challenge.URI != "" {
			resp["totp_uri"] = challenge.URI
		}
		if len(challenge.Data) > 0 {
			for k, v := range challenge.Data {
				resp[k] = v
			}
		}
		respond.JSON(w, http.StatusAccepted, resp)
		return
	}

	m.writeTokenResponse(w, http.StatusOK, user, accessToken, refreshToken)
}

func (m *Module) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	token := strings.TrimSpace(req.RefreshToken)
	req.RefreshToken = token
	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		respond.Validation(w, fields)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	stored, err := m.store.GetRefreshToken(ctx, hashToken(token))
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, errors.New("invalid refresh token"))
		return
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		if stored.RevokedAt != nil {
			m.rt.Logger.Warn("refresh token reuse attempt", zap.String("token_id", stored.ID.String()), zap.String("user_id", stored.UserID.String()))
		} else {
			m.rt.Logger.Info("refresh token expired", zap.String("token_id", stored.ID.String()), zap.String("user_id", stored.UserID.String()))
		}
		respond.Error(w, http.StatusUnauthorized, errors.New("refresh token expired"))
		return
	}

	user, err := m.store.GetByID(ctx, stored.UserID.String())
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, errors.New("invalid refresh token"))
		return
	}

	if err := m.store.RevokeRefreshToken(ctx, stored.ID); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, refreshToken, err := m.issueTokenPair(ctx, stored.UserID, user)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	m.writeTokenResponse(w, http.StatusOK, user, accessToken, refreshToken)
}

func (m *Module) handleSignout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	token := strings.TrimSpace(req.RefreshToken)
	req.RefreshToken = token
	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		respond.Validation(w, fields)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	stored, err := m.store.GetRefreshToken(ctx, hashToken(token))
	if err == nil {
		_ = m.store.RevokeRefreshToken(ctx, stored.ID)
	}

	respond.JSON(w, http.StatusOK, map[string]string{"status": "signed_out"})
}

func (m *Module) renderPasswordResetPage(w http.ResponseWriter, _ *http.Request) {
	m.renderResetRequestTemplate(w, http.StatusOK, templateData{})
}

func (m *Module) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	if !m.enforceRateLimit(r, "password_reset", true) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	isJSON := isJSONRequest(r)
	var email string
	if isJSON {
		var req passwordResetRequest
		if err := respond.Decode(r, &req); err != nil {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}
		email = strings.TrimSpace(strings.ToLower(req.Email))
	} else {
		if err := r.ParseForm(); err != nil {
			m.renderResetRequestTemplate(w, http.StatusBadRequest, templateData{Error: "Invalid form submission."})
			return
		}
		email = strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	}

	if !m.enforceRateLimit(r, "password_reset", false, email) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	req := passwordResetRequest{Email: email}
	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		if isJSON {
			respond.Validation(w, fields)
		} else {
			var msg string
			for _, mMsg := range fields {
				msg = fmt.Sprintf("Email %s.", strings.TrimPrefix(mMsg, "is "))
				break
			}
			if msg == "" {
				msg = "Invalid email address."
			}
			m.renderResetRequestTemplate(w, http.StatusBadRequest, templateData{Error: msg})
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := m.store.FindByEmail(ctx, email)
	if err != nil {
		m.respondPasswordResetAck(w, isJSON)
		return
	}

	token, err := generateSecureToken(refreshTokenSize)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	if _, err := m.store.CreatePasswordReset(ctx, user.ID, hashToken(token), time.Now().Add(passwordResetTTL)); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	resetLink := fmt.Sprintf("http://localhost:8080/auth/password/confirm?token=%s", token)
	m.rt.Logger.Info("password reset issued", zap.String("email", user.Email), zap.String("link", resetLink))
	if err := m.sendPasswordResetEmail(ctx, user.Email, resetLink); err != nil {
		m.rt.Logger.Error("send password reset email", zap.Error(err))
		if isJSON {
			respond.Error(w, http.StatusInternalServerError, errors.New("failed to send password reset email"))
		} else {
			m.renderResetRequestTemplate(w, http.StatusInternalServerError, templateData{Error: "Unable to send reset email at this time."})
		}
		return
	}

	m.respondPasswordResetAck(w, isJSON)
}

func (m *Module) renderPasswordConfirmPage(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	m.renderResetConfirmTemplate(w, http.StatusOK, templateData{Token: token})
}

func (m *Module) handlePasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	if !m.enforceRateLimit(r, "password_reset_confirm", true) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	isJSON := isJSONRequest(r)
	var req passwordResetConfirmRequest
	if isJSON {
		if err := respond.Decode(r, &req); err != nil {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			m.renderResetConfirmTemplate(w, http.StatusBadRequest, templateData{Error: "Invalid form submission."})
			return
		}
		req.Token = r.FormValue("token")
		req.Password = r.FormValue("password")
	}

	req.Token = strings.TrimSpace(req.Token)
	req.Password = strings.TrimSpace(req.Password)

	if !m.enforceRateLimit(r, "password_reset_confirm", false, req.Token) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	fields, err := m.validatePayload(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if len(fields) > 0 {
		if isJSON {
			respond.Validation(w, fields)
		} else {
			var msg string
			for field, fieldMsg := range fields {
				label := field
				if len(label) > 0 {
					label = strings.ToUpper(label[:1]) + label[1:]
				}
				msg = fmt.Sprintf("%s %s.", label, fieldMsg)
				break
			}
			if msg == "" {
				msg = "Validation failed."
			}
			m.renderResetConfirmTemplate(w, http.StatusBadRequest, templateData{Token: req.Token, Error: msg})
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	reset, err := m.store.GetPasswordReset(ctx, hashToken(req.Token))
	if err != nil {
		m.respondPasswordResetInvalid(w, isJSON, req.Token)
		return
	}
	if reset.UsedAt != nil || time.Now().After(reset.ExpiresAt) {
		m.respondPasswordResetInvalid(w, isJSON, req.Token)
		return
	}

	userID := reset.UserID
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err := m.store.UpdatePassword(ctx, userID, string(hash)); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err := m.store.MarkPasswordResetUsed(ctx, reset.ID); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	if err := m.store.RevokeUserTokens(ctx, userID); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	m.writePasswordResetSuccess(w, isJSON)
}

func (m *Module) handleProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := m.authenticateRequest(r)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	user, err := m.store.GetByID(ctx, claims.UserID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]any{
		"id":                      user.ID.String(),
		"email":                   user.Email,
		"roles":                   user.Roles,
		"role":                    user.Role,
		"password_reset_required": user.PasswordResetRequired,
	})
}

func (m *Module) authenticateRequest(r *http.Request) (*Claims, error) {
	token, err := tokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	claims, err := m.tokens.Parse(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *Module) issueTokenPair(ctx context.Context, userID uuid.UUID, user User) (string, string, error) {
	if len(user.Roles) == 0 {
		role := strings.TrimSpace(strings.ToLower(user.Role))
		if role == "" {
			role = defaultRoleName
		}
		user.Role = role
		user.Roles = []string{role}
	}
	accessToken, err := m.tokens.Generate(user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateSecureToken(refreshTokenSize)
	if err != nil {
		return "", "", err
	}
	if _, err := m.store.InsertRefreshToken(ctx, userID, hashToken(refreshToken), time.Now().Add(refreshTokenTTL)); err != nil {
		return "", "", err
	}
	if err := m.store.DeleteExpiredRefreshTokens(ctx, userID); err != nil {
		m.rt.Logger.Warn("delete expired refresh tokens", zap.Error(err))
	}
	if err := m.store.TrimActiveRefreshTokens(ctx, userID, maxActiveRefreshTokens); err != nil {
		m.rt.Logger.Warn("trim refresh tokens", zap.Error(err))
	}
	return accessToken, refreshToken, nil
}

func (m *Module) writeTokenResponse(w http.ResponseWriter, status int, user User, accessToken, refreshToken string) {
	respond.JSON(w, status, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(m.tokens.TTL().Seconds()),
		"user": map[string]any{
			"id":                      user.ID.String(),
			"email":                   user.Email,
			"roles":                   user.Roles,
			"role":                    user.Role,
			"password_reset_required": user.PasswordResetRequired,
		},
	})
}

func generateSecureToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("token entropy: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sum[:])
}

func isJSONRequest(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json")
}

func (m *Module) respondPasswordResetAck(w http.ResponseWriter, asJSON bool) {
	if asJSON {
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	m.renderResetRequestTemplate(w, http.StatusOK, templateData{Message: "If the account exists, a reset link has been issued. Please check the application logs for details."})
}

func (m *Module) respondPasswordResetInvalid(w http.ResponseWriter, asJSON bool, token string) {
	if asJSON {
		respond.Error(w, http.StatusBadRequest, errors.New("reset token invalid or expired"))
		return
	}
	m.renderResetConfirmTemplate(w, http.StatusBadRequest, templateData{Token: token, Error: "Reset token is invalid or has expired."})
}

func (m *Module) writePasswordResetSuccess(w http.ResponseWriter, asJSON bool) {
	if asJSON {
		respond.JSON(w, http.StatusOK, map[string]string{"status": "password_updated"})
		return
	}
	m.renderResetConfirmTemplate(w, http.StatusOK, templateData{Message: "Password updated. You may now sign in with your new password."})
}

func (m *Module) handleMFAVerify(w http.ResponseWriter, r *http.Request) {
	if m.mfaProvider == nil {
		respond.Error(w, http.StatusNotFound, errors.New("multi-factor authentication not configured"))
		return
	}

	if !m.enforceRateLimit(r, "mfa_verify", true) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	var req struct {
		ChallengeID string `json:"challenge_id"`
		Code        string `json:"code"`
	}
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	req.ChallengeID = strings.TrimSpace(req.ChallengeID)
	req.Code = strings.TrimSpace(req.Code)
	if req.ChallengeID == "" || req.Code == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("challenge_id and code are required"))
		return
	}

	if !m.enforceRateLimit(r, "mfa_verify", false, req.ChallengeID) {
		respond.Error(w, http.StatusTooManyRequests, errTooManyRequests)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := m.mfaProvider.Verify(ctx, req.ChallengeID, req.Code)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err)
		return
	}

	user := result.User
	accessToken, refreshToken, err := m.issueTokenPair(ctx, user.ID, user)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	if m.rt != nil {
		m.rt.Logger.Info("user passed MFA", zap.String("email", user.Email))
	}
	m.writeTokenResponse(w, http.StatusOK, user, accessToken, refreshToken)
}

func (m *Module) enforceRateLimit(r *http.Request, action string, includeIP bool, extraKeys ...string) bool {
	if m.rateLimiter == nil {
		return true
	}

	keys := make(map[string]struct{})
	if includeIP {
		if ip := clientIdentifier(r); ip != "" {
			keys[ip] = struct{}{}
		}
	}

	for _, key := range extraKeys {
		key = strings.TrimSpace(strings.ToLower(key))
		if key == "" {
			continue
		}
		keys[key] = struct{}{}
	}

	if len(keys) == 0 {
		return true
	}

	for key := range keys {
		if !m.rateLimiter.Allow(action, key) {
			if m.rt != nil {
				m.rt.Logger.Warn("auth rate limit exceeded",
					zap.String("action", action),
					zap.String("key_hash", hashKey(key)))
			}
			return false
		}
	}
	return true
}

func clientIdentifier(r *http.Request) string {
	xForwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func hashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", sum[:8])
}

func (m *Module) validatePayload(payload any) (map[string]string, error) {
	if m.validator == nil {
		return nil, nil
	}
	if err := m.validator.Struct(payload); err != nil {
		var verr validate.ValidationError
		if errors.As(err, &verr) {
			return verr.Fields, nil
		}
		return nil, err
	}
	return nil, nil
}

func (m *Module) sendPasswordResetEmail(ctx context.Context, to, link string) error {
	if m.mailer == nil {
		m.rt.Logger.Info("password reset email not sent (mailer disabled)", zap.String("email", to), zap.String("link", link))
		return nil
	}
	subject := "Password reset instructions"
	body := fmt.Sprintf("You requested a password reset. Use the link below to set a new password:\n\n%s\n\nIf you did not request this change, you can ignore this message.", link)
	if err := m.mailer.Send(ctx, to, subject, body); err != nil {
		return err
	}
	m.rt.Logger.Info("password reset email sent", zap.String("email", to))
	return nil
}

func (m *Module) renderResetRequestTemplate(w http.ResponseWriter, status int, data templateData) {
	m.renderTemplate(w, status, "reset_request", data)
}

func (m *Module) renderResetConfirmTemplate(w http.ResponseWriter, status int, data templateData) {
	m.renderTemplate(w, status, "reset_confirm", data)
}

func (m *Module) renderTemplate(w http.ResponseWriter, status int, name string, data templateData) {
	data = m.decorateTemplateData(data)
	if m.views == nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, "%s\n%s", name, data.Message)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := m.views.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template rendering error", http.StatusInternalServerError)
	}
}

func (m *Module) decorateTemplateData(data templateData) templateData {
	if data.FrameworkName == "" {
		data.FrameworkName = "Shanraq Framework"
	}
	if data.FrameworkDescription == "" {
		data.FrameworkDescription = "Batteries-included Go platform for resilient services."
	}
	if data.BrandLogoPath == "" {
		data.BrandLogoPath = "/static/brand/logo.svg"
	}
	if data.Year == 0 {
		data.Year = time.Now().Year()
	}
	return data
}

func (m *Module) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	normalized := sanitizeRoles(roles)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				var err error
				claims, err = m.authenticateRequest(r)
				if err != nil {
					respond.Error(w, http.StatusUnauthorized, err)
					return
				}
			}
			if len(normalized) > 0 && !claims.HasAnyRole(normalized...) {
				respond.Error(w, http.StatusForbidden, errors.New("insufficient role to access this resource"))
				return
			}
			nextCtx := ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(nextCtx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	if claims == nil {
		return ctx
	}
	return context.WithValue(ctx, claimsContextKey, claims)
}

func sanitizeRoles(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		role = strings.TrimSpace(strings.ToLower(role))
		if role == "" {
			continue
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		result = append(result, role)
	}
	return result
}

func isWeakSecret(secret string) bool {
	normalized := strings.TrimSpace(strings.ToLower(secret))
	if len(normalized) < 24 {
		switch normalized {
		case "", "replace-me-now", "super-secret-key", "secret", "changeme":
			return true
		}
	}
	return false
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
