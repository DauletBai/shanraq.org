package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"shanraq.org/pkg/shanraq"
)

// Module exposes canonical health and readiness endpoints.
type Module struct {
	rt *shanraq.Runtime
}

// New returns a ready-to-use health module.
func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "health"
}

// Init captures runtime dependencies for later use in handlers.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	return nil
}

// Routes wires /healthz and /readyz.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}

	r.Get("/healthz", m.handleLiveness)
	r.Get("/readyz", m.handleReadiness)
}

func (m *Module) handleLiveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":      "ok",
		"environment": m.rt.Config.Environment,
	})
}

func (m *Module) handleReadiness(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := m.rt.DB.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "degraded",
			"error":  err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
