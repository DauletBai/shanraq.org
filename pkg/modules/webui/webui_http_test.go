package webui

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/web"
)

type stubJobsProvider struct {
	metrics jobs.MetricsSnapshot
	counts  map[string]int
	recent  []jobs.Job
}

func (s stubJobsProvider) Metrics(_ context.Context, _ *uuid.UUID) (jobs.MetricsSnapshot, error) {
	return s.metrics, nil
}

func (s stubJobsProvider) CountByStatus(_ context.Context, _ *uuid.UUID) (map[string]int, error) {
	return s.counts, nil
}

func (s stubJobsProvider) ListRecent(_ context.Context, _ int, _ *uuid.UUID) ([]jobs.Job, error) {
	return s.recent, nil
}

func newTestModule(t *testing.T) *Module {
	t.Helper()
	renderer, err := web.NewRenderer()
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}

	mod := New(4, 2*time.Second)
	mod.renderer = renderer
	mod.jobsData = stubJobsProvider{
		metrics: jobs.MetricsSnapshot{
			Pending:         2,
			Running:         1,
			Retry:           1,
			Failed:          1,
			Done:            5,
			DoneLastHour:    3,
			FailedLastHour:  1,
			Total:           10,
			NextScheduled:   time.Now().Add(5 * time.Minute),
			NextScheduledOk: true,
		},
		counts: map[string]int{
			"pending": 2,
			"running": 1,
			"retry":   1,
			"failed":  1,
			"done":    5,
		},
		recent: []jobs.Job{
			{
				Name:        "demo_job",
				Status:      "done",
				Attempts:    1,
				MaxAttempts: 3,
				RunAt:       time.Now(),
			},
		},
	}
	mod.aboutContent = &AboutContent{
		Headline:     "Test Headline",
		Subheadline:  "Test Subheadline",
		FeatureOne:   "Feature 1",
		FeatureTwo:   "Feature 2",
		FeatureThree: "Feature 3",
	}
	mod.rt = &shanraq.Runtime{
		Config: config.Config{Environment: "test"},
		Logger: zap.NewNop(),
		Router: chi.NewRouter(),
	}
	return mod
}

func TestRoutesHome(t *testing.T) {
	mod := newTestModule(t)
	router := chi.NewRouter()
	mod.Routes(router)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "id=\"home-root\"") {
		t.Fatalf("expected home content")
	}
}

func TestRoutesConsole(t *testing.T) {
	mod := newTestModule(t)
	router := chi.NewRouter()
	mod.Routes(router)

	req := httptest.NewRequest("GET", "/console", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "id=\"dashboard-root\"") {
		t.Fatalf("expected dashboard content")
	}
}

func TestRoutesDocs(t *testing.T) {
	mod := newTestModule(t)
	router := chi.NewRouter()
	mod.Routes(router)

	req := httptest.NewRequest("GET", "/docs", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Framework Documentation") {
		t.Fatalf("expected docs content")
	}
	if !strings.Contains(body, "Configuration") {
		t.Fatalf("expected configuration section")
	}
}
