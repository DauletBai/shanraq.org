package webui

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/web"
)

// Module renders a Bootstrap-powered dashboard for operators.
type Module struct {
	rt           *shanraq.Runtime
	renderer     *web.Renderer
	jobsStore    *jobs.Store
	workers      int
	pollInterval time.Duration
	reloadPeriod time.Duration
	aboutContent *AboutContent
}

// New constructs a module with the provided worker metadata (for UI display).
func New(workers int, pollInterval time.Duration) *Module {
	return &Module{
		workers:      workers,
		pollInterval: pollInterval,
		reloadPeriod: 15 * time.Second,
	}
}

func (m *Module) Name() string {
	return "webui"
}

// Init bootstraps the renderer and any required stores.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	renderer, err := web.NewRenderer()
	if err != nil {
		return err
	}
	m.rt = rt
	m.renderer = renderer
	m.jobsStore = jobs.NewStore(rt.DB)
	m.aboutContent = loadAboutContent(rt.DB)
	if m.aboutContent == nil {
		m.aboutContent = &AboutContent{
			Headline:     "Shanraq Console",
			Subheadline:  "A Go-first framework for resilient backends.",
			FeatureOne:   "PostgreSQL-native data layer with migrations built-in.",
			FeatureTwo:   "Composable module system for HTTP, workers, and observability.",
			FeatureThree: "Cloud-ready tooling: Docker, Prometheus telemetry, structured logging.",
		}
	}
	return nil
}

// Routes exposes the main dashboard.
func (m *Module) Routes(r chi.Router) {
	if m.renderer == nil {
		return
	}
	r.Handle("/static/*", http.StripPrefix("/static/", web.StaticHandler()))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/brand/favicon.svg", http.StatusMovedPermanently)
	})
	r.Get("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.StaticFS(), "brand/favicon.svg")
	})
	r.Get("/", m.handleHome)
	r.Get("/console", m.handleDashboard)
	r.Get("/docs", m.handleDocs)
	r.Get("/partials/dashboard", m.handleDashboardPartial)
}

func (m *Module) handleHome(w http.ResponseWriter, r *http.Request) {
	data := m.buildDashboardData(r.Context())
	if data == nil {
		http.Error(w, "failed to build home", http.StatusInternalServerError)
		return
	}
	data.SidebarLinks = nil
	data.PageID = "home"
	data.PageTitle = data.FrameworkName

	if err := m.renderer.Render(w, "home.html", *data); err != nil {
		m.rt.Logger.Error("render home", zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (m *Module) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := m.buildDashboardData(r.Context())
	if data == nil {
		http.Error(w, "failed to build dashboard", http.StatusInternalServerError)
		return
	}
	data.PageID = "dashboard"

	if err := m.renderer.Render(w, "dashboard.html", *data); err != nil {
		m.rt.Logger.Error("render dashboard", zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (m *Module) handleDashboardPartial(w http.ResponseWriter, r *http.Request) {
	data := m.buildDashboardData(r.Context())
	if data == nil {
		http.Error(w, "failed to build dashboard", http.StatusInternalServerError)
		return
	}
	data.PageID = "dashboard"

	buf := new(bytes.Buffer)
	if err := m.renderer.Template().ExecuteTemplate(buf, "dashboard-content", *data); err != nil {
		m.rt.Logger.Error("render dashboard partial", zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

func (m *Module) buildDashboardData(parent context.Context) *DashboardData {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	metrics, err := m.jobsStore.Metrics(ctx)
	if err != nil {
		m.rt.Logger.Warn("job metrics", zap.Error(err))
	}

	counts := map[string]int{
		"pending": metrics.Pending,
		"retry":   metrics.Retry,
		"running": metrics.Running,
		"failed":  metrics.Failed,
		"done":    metrics.Done,
	}
	if metrics.Total == 0 {
		if fallback, err := m.jobsStore.CountByStatus(ctx); err == nil {
			counts = fallback
			metrics.Pending = fallback["pending"]
			metrics.Retry = fallback["retry"]
			metrics.Running = fallback["running"]
			metrics.Failed = fallback["failed"]
			metrics.Done = fallback["done"]
			metrics.Total = totalJobs(fallback)
		}
	}

	recent, err := m.jobsStore.ListRecent(ctx, 8)
	if err != nil {
		m.rt.Logger.Warn("recent jobs", zap.Error(err))
		recent = []jobs.Job{}
	}

	frameworkName := "Shanraq"
	frameworkDescription := "Shanraq is a modular Go 1.25 framework for PostgreSQL-first services."
	if m.aboutContent != nil {
		frameworkDescription = m.aboutContent.Subheadline
	}

	total := totalJobs(counts)
	pending := counts["pending"]
	retrying := counts["retry"]
	running := counts["running"]
	failed := counts["failed"]
	completed := counts["done"]
	jobStats := JobStats{
		Pending:      pending,
		Retrying:     retrying,
		Running:      running,
		Failed:       failed,
		Completed:    completed,
		Total:        total,
		Workers:      m.workers,
		PollInterval: m.pollInterval.String(),
	}

	var queueOverview QueueOverview
	windowTotal := metrics.DoneLastHour + metrics.FailedLastHour
	var failureRate, successRate float64
	if windowTotal > 0 {
		failureRate = float64(metrics.FailedLastHour) / float64(windowTotal)
		successRate = float64(metrics.DoneLastHour) / float64(windowTotal)
	}
	queueOverview = QueueOverview{
		DoneLastHour:       metrics.DoneLastHour,
		FailedLastHour:     metrics.FailedLastHour,
		SuccessRate:        successRate,
		FailureRate:        failureRate,
		NextScheduled:      metrics.NextScheduled,
		NextScheduledValid: metrics.NextScheduledOk,
	}

	return &DashboardData{
		PageID:               "dashboard",
		FrameworkName:        frameworkName,
		FrameworkDescription: frameworkDescription,
		FrameworkBadge:       "Go 1.25 ready",
		PageTitle:            "Shanraq Console",
		CurrentYear:          time.Now().Year(),
		Environment:          m.rt.Config.Environment,
		LastUpdated:          time.Now(),
		ReloadAfter:          m.reloadPeriod,
		SidebarLinks: []SidebarLink{
			{Label: "Home", Href: "/"},
			{Label: "Console", Href: "/console"},
			{Label: "Documentation", Href: "/docs"},
			{Label: "Jobs API", Href: "/jobs"},
			{Label: "Health", Href: "/healthz"},
			{Label: "Readiness", Href: "/readyz"},
			{Label: "Metrics", Href: "/metrics"},
		},
		JobStats:   jobStats,
		RecentJobs: recent,
		Health:     m.healthSnapshot(ctx, counts),
		JobsAPIURL: "/jobs",
		MetricsURL: metricsPath(m.rt.Config.Telemetry),
		JobForm: JobFormDefaults{
			Name:        "send_welcome_email",
			MaxAttempts: 3,
			Payload:     "{\n  \"email\": \"demo@shanraq.org\"\n}",
		},
		Hero:          m.aboutContent,
		QueueOverview: queueOverview,
	}
}

// JobStats feeds the dashboard hero cards.
type JobStats struct {
	Pending      int
	Retrying     int
	Running      int
	Failed       int
	Completed    int
	Total        int
	Workers      int
	PollInterval string
}

// DashboardData is the template context for dashboard.html.
type DashboardData struct {
	PageID               string
	FrameworkName        string
	FrameworkDescription string
	FrameworkBadge       string
	PageTitle            string
	CurrentYear          int
	Environment          string
	LastUpdated          time.Time
	ReloadAfter          time.Duration
	SidebarLinks         []SidebarLink
	JobStats             JobStats
	RecentJobs           []jobs.Job
	Health               []HealthIndicator
	JobsAPIURL           string
	MetricsURL           string
	JobForm              JobFormDefaults
	Hero                 *AboutContent
	Docs                 []DocSection
	QueueOverview        QueueOverview
}

type HealthIndicator struct {
	Name        string
	Status      string
	Level       string
	Description string
	Link        string
}

type JobFormDefaults struct {
	Name        string
	Payload     string
	MaxAttempts int
}

type AboutContent struct {
	Headline     string
	Subheadline  string
	FeatureOne   string
	FeatureTwo   string
	FeatureThree string
}

// SidebarLink renders binary safe navigation items.
type SidebarLink struct {
	Label string
	Href  string
}

type DocSection struct {
	Title   string
	Summary string
	Items   []DocItem
}

type DocItem struct {
	Title       string
	Description string
	Link        string
	Command     string
}

type QueueOverview struct {
	DoneLastHour       int
	FailedLastHour     int
	SuccessRate        float64
	FailureRate        float64
	NextScheduled      time.Time
	NextScheduledValid bool
}

func totalJobs(counts map[string]int) int {
	var total int
	for _, v := range counts {
		total += v
	}
	return total
}

func metricsPath(t config.Telemetry) string {
	if t.MetricsPath != "" {
		return t.MetricsPath
	}
	return "/metrics"
}

func loadAboutContent(db *pgxpool.Pool) *AboutContent {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := db.QueryRow(ctx, `
		SELECT headline, subheadline, feature_one, feature_two, feature_three
		FROM framework_about
		ORDER BY created_at DESC
		LIMIT 1
	`)
	var content AboutContent
	if err := row.Scan(&content.Headline, &content.Subheadline, &content.FeatureOne, &content.FeatureTwo, &content.FeatureThree); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return nil
	}
	return &content
}

func (m *Module) healthSnapshot(ctx context.Context, counts map[string]int) []HealthIndicator {
	indicators := make([]HealthIndicator, 0, 4)
	status := "Operational"
	level := "success"
	stat := m.rt.DB.Stat()
	description := fmt.Sprintf("%d total · %d idle · %d in-use", stat.TotalConns(), stat.IdleConns(), stat.AcquiredConns())
	if err := m.rt.DB.Ping(ctx); err != nil {
		status = "Degraded"
		level = "danger"
		description = err.Error()
	}
	indicators = append(indicators, HealthIndicator{
		Name:        "PostgreSQL",
		Status:      status,
		Level:       level,
		Description: description,
	})

	workerStatus := "Operational"
	workerLevel := "info"
	if m.workers == 0 {
		workerStatus = "Paused"
		workerLevel = "warning"
	}
	indicators = append(indicators, HealthIndicator{
		Name:        "Workers",
		Status:      workerStatus,
		Level:       workerLevel,
		Description: fmt.Sprintf("%d workers · poll %s", m.workers, m.pollInterval),
	})

	failed := counts["failed"]
	if failed > 0 {
		indicators = append(indicators, HealthIndicator{
			Name:        "Failed Jobs",
			Status:      fmt.Sprintf("%d pending attention", failed),
			Level:       "warning",
			Description: "Investigate recent failures via Jobs API.",
			Link:        "/jobs",
		})
	} else {
		indicators = append(indicators, HealthIndicator{
			Name:        "Failed Jobs",
			Status:      "None",
			Level:       "success",
			Description: "Queue is healthy.",
		})
	}

	indicators = append(indicators, HealthIndicator{
		Name:        "Telemetry",
		Status:      "Scrape ready",
		Level:       "info",
		Description: "Prometheus endpoint exposed.",
		Link:        metricsPath(m.rt.Config.Telemetry),
	})

	return indicators
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)

func (m *Module) handleDocs(w http.ResponseWriter, r *http.Request) {
	data := m.buildDashboardData(r.Context())
	if data == nil {
		http.Error(w, "failed to build docs", http.StatusInternalServerError)
		return
	}
	data.PageID = "docs"
	data.PageTitle = "Framework Documentation"
	data.Docs = []DocSection{
		{
			Title:   "Quick Start",
			Summary: "Bootstrap the reference application and explore core capabilities.",
			Items: []DocItem{
				{
					Title:       "Run Locally",
					Description: "Launch the development server with hot-reloaded templates and the PostgreSQL-backed queue.",
					Command:     "go run ./cmd/app -config config.yaml",
				},
				{
					Title:       "Docker Compose",
					Description: "Build the distroless image and start the application together with PostgreSQL.",
					Command:     "docker compose up --build",
				},
			},
		},
		{
			Title:   "Available Modules",
			Summary: "Each module follows the shanraq.Module contract and can be composed independently.",
			Items: []DocItem{
				{
					Title:       "Authentication",
					Description: "JWT-based signup/signin/profile endpoints under /auth with refresh token rotation and password reset flows.",
					Link:        "/auth/signup",
				},
				{
					Title:       "Background Jobs",
					Description: "PostgreSQL-backed queue with retry semantics and operator console.",
					Link:        "/jobs",
				},
				{
					Title:       "Telemetry & Health",
					Description: "Prometheus metrics, health, and readiness endpoints.",
					Link:        "/metrics",
				},
				{
					Title:       "Job Context Helpers",
					Description: "Access worker index and attempt metadata via jobs.InfoFromContext(ctx).",
					Command:     "meta, _ := jobs.InfoFromContext(ctx)",
				},
			},
		},
		{
			Title:   "Operator Console",
			Summary: "Navigate the Bootstrap dashboard to monitor system health and manage jobs.",
			Items: []DocItem{
				{
					Title:       "Landing Page",
					Description: "Marketing-ready carousel with feature highlights and live stats.",
					Link:        "/",
				},
				{
					Title:       "Console Dashboard",
					Description: "Queue explorer, job submission modal, and system health cards.",
					Link:        "/console",
				},
				{
					Title:       "Jobs API",
					Description: "REST endpoints for enqueueing, retrying, and cancelling jobs.",
					Link:        "/jobs",
				},
			},
		},
	}

	if err := m.renderer.Render(w, "docs.html", *data); err != nil {
		m.rt.Logger.Error("render docs", zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
