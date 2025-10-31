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
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/respond"
)

const (
	refreshTokenTTL  = 30 * 24 * time.Hour
	passwordResetTTL = time.Hour
	refreshTokenSize = 48
)

// Module wires auth routes plus helper services (store + tokens).
type Module struct {
	rt     *shanraq.Runtime
	store  *Store
	tokens *TokenService
	views  *template.Template
}

//go:embed templates/*.html
var viewFiles embed.FS

type templateData struct {
	Token   string
	Message string
	Error   string
}

// New returns a ready-to-register authentication module.
func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "auth"
}

// Init prepares dependencies shared across handlers.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	m.tokens = NewTokenService(rt.Config.Auth.TokenSecret, rt.Config.Auth.TokenTTL)
	tmpl, err := template.ParseFS(viewFiles, "templates/*.html")
	if err != nil {
		return fmt.Errorf("parse auth templates: %w", err)
	}
	m.views = tmpl
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
	})
}

func (m *Module) handleSignup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("email and password required"))
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
	var req signinRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("email and password required"))
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

	m.writeTokenResponse(w, http.StatusOK, user, accessToken, refreshToken)
}

func (m *Module) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	token := strings.TrimSpace(req.RefreshToken)
	if token == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("refresh token required"))
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
	if token == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("refresh token required"))
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

	if email == "" {
		if isJSON {
			respond.Error(w, http.StatusBadRequest, errors.New("email required"))
		} else {
			m.renderResetRequestTemplate(w, http.StatusBadRequest, templateData{Error: "Email is required."})
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

	m.respondPasswordResetAck(w, isJSON)
}

func (m *Module) renderPasswordConfirmPage(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	m.renderResetConfirmTemplate(w, http.StatusOK, templateData{Token: token})
}

func (m *Module) handlePasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
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
	if req.Token == "" || req.Password == "" {
		if isJSON {
			respond.Error(w, http.StatusBadRequest, errors.New("token and password required"))
		} else {
			m.renderResetConfirmTemplate(w, http.StatusBadRequest, templateData{Token: req.Token, Error: "Token and password are required."})
		}
		return
	}
	if len(req.Password) < 8 {
		if isJSON {
			respond.Error(w, http.StatusBadRequest, errors.New("password must be at least 8 characters"))
		} else {
			m.renderResetConfirmTemplate(w, http.StatusBadRequest, templateData{Token: req.Token, Error: "Password must be at least 8 characters."})
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

func (m *Module) renderResetRequestTemplate(w http.ResponseWriter, status int, data templateData) {
	m.renderTemplate(w, status, "reset_request", data)
}

func (m *Module) renderResetConfirmTemplate(w http.ResponseWriter, status int, data templateData) {
	m.renderTemplate(w, status, "reset_confirm", data)
}

func (m *Module) renderTemplate(w http.ResponseWriter, status int, name string, data templateData) {
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

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
