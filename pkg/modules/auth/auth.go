package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/respond"
)

// Module wires auth routes plus helper services (store + tokens).
type Module struct {
	rt     *shanraq.Runtime
	store  *Store
	tokens *TokenService
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
		r.Get("/profile", m.handleProfile)
	})
}

func (m *Module) handleSignup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
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

	user, err := m.store.CreateUser(ctx, req.Email, string(hash))
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			respond.Error(w, http.StatusConflict, err)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	m.rt.Logger.Info("user signed up", zap.String("email", user.Email))
	respond.JSON(w, http.StatusCreated, map[string]any{
		"id":    user.ID.String(),
		"email": user.Email,
	})
}

func (m *Module) handleSignin(w http.ResponseWriter, r *http.Request) {
	var req signinRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
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

	token, err := m.tokens.Generate(user)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (m *Module) handleProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := m.authenticateRequest(r)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	user, err := m.store.GetByID(ctx, userID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{
		"id":    user.ID.String(),
		"email": user.Email,
	})
}

func (m *Module) authenticateRequest(r *http.Request) (string, error) {
	token, err := tokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		return "", err
	}
	claims, err := m.tokens.Parse(token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
