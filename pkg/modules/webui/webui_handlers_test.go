package webui

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap/zaptest"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/web"
)

type stubJobsData struct {
	metrics jobs.MetricsSnapshot
	counts  map[string]int
	recent  []jobs.Job
	err     error
}

func (s stubJobsData) Metrics(_ context.Context, _ *uuid.UUID) (jobs.MetricsSnapshot, error) {
	return s.metrics, s.err
}

func (s stubJobsData) CountByStatus(_ context.Context, _ *uuid.UUID) (map[string]int, error) {
	return s.counts, s.err
}

func (s stubJobsData) ListRecent(_ context.Context, _ int, _ *uuid.UUID) ([]jobs.Job, error) {
	return s.recent, s.err
}

func buildTestWebUIModule(t *testing.T) *Module {
	t.Helper()

	renderer, err := web.NewRenderer()
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}

	m := New(4, 2*time.Second)
	m.renderer = renderer
	m.rt = &shanraq.Runtime{
		Config: config.Config{Environment: "test"},
		Logger: zaptest.NewLogger(t),
	}
	m.aboutContent = &AboutContent{
		Headline:     "Test Framework",
		Subheadline:  "Testing the Shanraq console.",
		FeatureOne:   "Feature A",
		FeatureTwo:   "Feature B",
		FeatureThree: "Feature C",
	}
	m.jobsData = stubJobsData{
		metrics: jobs.MetricsSnapshot{
			Pending:         2,
			Retry:           1,
			Running:         1,
			Failed:          0,
			Done:            5,
			DoneLastHour:    5,
			FailedLastHour:  0,
			NextScheduled:   time.Now().Add(10 * time.Minute),
			NextScheduledOk: true,
		},
		counts: map[string]int{
			"pending": 2,
			"retry":   1,
			"running": 1,
			"failed":  0,
			"done":    5,
		},
		recent: []jobs.Job{
			{
				ID:        uuid.New(),
				Name:      "demo-job",
				Status:    "done",
				Payload:   []byte(`{"demo":true}`),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	return m
}

func TestWebUIRoutes_HomeAndDashboard(t *testing.T) {
	module := buildTestWebUIModule(t)
	router := chi.NewRouter()
	module.Routes(router)

	tests := []struct {
		name           string
		route          string
		expectStatus   int
		expectContains string
	}{
		{
			name:           "home",
			route:          "/",
			expectStatus:   http.StatusOK,
			expectContains: `id="home-root"`,
		},
		{
			name:           "dashboard",
			route:          "/console",
			expectStatus:   http.StatusOK,
			expectContains: `id="dashboard-root"`,
		},
		{
			name:           "docs",
			route:          "/docs",
			expectStatus:   http.StatusOK,
			expectContains: "Framework Documentation",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectStatus {
				t.Fatalf("expected status %d, got %d", tt.expectStatus, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tt.expectContains) {
				t.Fatalf("expected response to contain %q", tt.expectContains)
			}
		})
	}
}

func TestWebUIRoutes_DashboardPartial(t *testing.T) {
	module := buildTestWebUIModule(t)
	router := chi.NewRouter()
	module.Routes(router)

	req := httptest.NewRequest(http.MethodGet, "/partials/dashboard", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Queue Explorer") {
		t.Fatalf("expected partial to contain queue explorer content")
	}
	if ctype := rec.Header().Get("Content-Type"); !strings.Contains(ctype, "text/html") {
		t.Fatalf("expected HTML content-type, got %s", ctype)
	}
}
