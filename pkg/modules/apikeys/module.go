package apikeys

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/respond"
)

// Option configures the API keys module.
type Option func(*Module)

// WithHTTPMiddleware registers middleware that should guard management endpoints (e.g. auth.RequireRoles).
func WithHTTPMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(m *Module) {
		m.middlewares = append(m.middlewares, mw...)
	}
}

// WithLogger overrides the logger used for diagnostics.
func WithLogger(logger *zap.Logger) Option {
	return func(m *Module) {
		m.logger = logger
	}
}

// Module manages API key issuance and validation.
type Module struct {
	rt          *shanraq.Runtime
	store       *Store
	logger      *zap.Logger
	middlewares []func(http.Handler) http.Handler
}

func New(opts ...Option) *Module {
	m := &Module{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Module) Name() string {
	return "apikeys"
}

func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	if m.logger == nil {
		m.logger = rt.Logger
	}
	return nil
}

func (m *Module) Routes(r chi.Router) {
	if m.store == nil {
		return
	}
	r.Route("/auth/apikeys", func(r chi.Router) {
		for _, mw := range m.middlewares {
			r.Use(mw)
		}
		r.Get("/", m.handleList)
		r.Post("/", m.handleCreate)
		r.Delete("/{id}", m.handleRevoke)
	})
}

// RequireAPIKey authenticates incoming requests using the X-API-Key or Authorization: ApiKey header.
func (m *Module) RequireAPIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := extractAPIKey(r)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			user, apiKey, err := m.store.Validate(r.Context(), key)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, errors.New("invalid api key"))
				return
			}
			claims := auth.ClaimsForUser(user)
			ctx := auth.ContextWithClaims(r.Context(), claims)
			ctx = context.WithValue(ctx, apiKeyContextKey{}, apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// APIKeyFromContext returns the validated API key metadata when present.
func APIKeyFromContext(ctx context.Context) (APIKey, bool) {
	key, ok := ctx.Value(apiKeyContextKey{}).(APIKey)
	return key, ok
}

type apiKeyContextKey struct{}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.userIDFromContext(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}

	keys, err := m.store.List(r.Context(), userID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, keys)
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.userIDFromContext(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}

	var req createRequest
	if err := respond.Decode(r, &req); err != nil && !errors.Is(err, io.EOF) {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	plain, key, err := m.store.Create(r.Context(), userID, strings.TrimSpace(req.Label))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]any{
		"api_key": plain,
		"key":     key,
	})
}

func (m *Module) handleRevoke(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.userIDFromContext(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, errors.New("missing auth claims"))
		return
	}

	keyID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, errors.New("invalid api key id"))
		return
	}
	if err := m.store.Revoke(r.Context(), userID, keyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, errors.New("api key not found"))
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

func (m *Module) userIDFromContext(r *http.Request) (uuid.UUID, bool) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, false
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}
	return userID, true
}

func extractAPIKey(r *http.Request) string {
	if key := strings.TrimSpace(r.Header.Get("X-API-Key")); key != "" {
		return key
	}
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(authHeader), "apikey ") {
		return strings.TrimSpace(authHeader[7:])
	}
	return ""
}
