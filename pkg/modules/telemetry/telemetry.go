package telemetry

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"
	internaltelemetry "shanraq.org/internal/telemetry"
	"shanraq.org/pkg/shanraq"
)

// Module exposes Prometheus metrics handler and optional tracing middleware when enabled.
type Module struct {
	handler          http.Handler
	metricsPath      string
	metricsEnabled   bool
	metricsToken     string
	environment      string
	log              *zap.Logger
	tracingEnabled   bool
	tracingShutdown  func(context.Context) error
	serviceName      string
	tracerMiddleware bool
}

// New creates a telemetry module bound to application configuration.
func New() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "telemetry"
}

func (m *Module) Init(ctx context.Context, rt *shanraq.Runtime) error {
	m.log = rt.Logger
	m.environment = rt.Config.Environment
	if rt.Config.Telemetry.EnableMetrics {
		m.handler = promhttp.Handler()
		path := rt.Config.Telemetry.MetricsPath
		if path == "" {
			path = "/metrics"
		}
		m.metricsPath = path
		m.metricsToken = strings.TrimSpace(rt.Config.Telemetry.MetricsToken)
		m.metricsEnabled = true
	}

	shutdown, err := internaltelemetry.SetupTracing(ctx, rt.Config.Telemetry.Tracing, rt.Config.Environment, rt.Logger)
	if err != nil {
		return err
	}
	m.tracingShutdown = shutdown
	m.tracingEnabled = rt.Config.Telemetry.Tracing.Enabled
	m.serviceName = rt.Config.Telemetry.Tracing.ServiceName
	if m.serviceName == "" {
		m.serviceName = "shanraq-app"
	}
	return nil
}

func (m *Module) Routes(r chi.Router) {
	if m.tracingEnabled && !m.tracerMiddleware {
		r.Use(otelchi.Middleware(m.serviceName))
		m.tracerMiddleware = true
	}
	if m.metricsEnabled {
		r.Handle(m.metricsPath, m.guardMetrics(m.handler))
	}
}

// guardMetrics protects the Prometheus endpoint. With a token configured it
// requires "Authorization: Bearer <token>" (constant-time compare). Without a
// token it serves only outside production — production config validation already
// requires a token when metrics are enabled, so this is a defense-in-depth net.
func (m *Module) guardMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.metricsToken == "" {
			if strings.EqualFold(m.environment, "production") {
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		const prefix = "Bearer "
		got := r.Header.Get("Authorization")
		if len(got) > len(prefix) && strings.EqualFold(got[:len(prefix)], prefix) &&
			subtle.ConstantTimeCompare([]byte(got[len(prefix):]), []byte(m.metricsToken)) == 1 {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

func (m *Module) Start(ctx context.Context, rt *shanraq.Runtime) error {
	<-ctx.Done()

	if m.tracingShutdown != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.tracingShutdown(shutdownCtx); err != nil && rt.Logger != nil {
			rt.Logger.Warn("telemetry tracing shutdown", zap.Error(err))
		}
	}
	return ctx.Err()
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
	shanraq.StarterModule
} = (*Module)(nil)
