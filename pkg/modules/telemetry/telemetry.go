package telemetry

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"shanraq.org/pkg/shanraq"
)

// Module exposes Prometheus metrics handler when enabled.
type Module struct {
	handler     http.Handler
	metricsPath string
	enabled     bool
}

// New creates a telemetry module bound to application configuration.
func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "telemetry"
}

func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	if !rt.Config.Telemetry.EnableMetrics {
		return nil
	}
	m.handler = promhttp.Handler()
	path := rt.Config.Telemetry.MetricsPath
	if path == "" {
		path = "/metrics"
	}
	m.metricsPath = path
	m.enabled = true
	return nil
}

func (m *Module) Routes(r chi.Router) {
	if !m.enabled {
		return
	}
	r.Handle(m.metricsPath, m.handler)
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
