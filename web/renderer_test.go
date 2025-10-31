package web_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"shanraq.org/pkg/modules/webui"
	web "shanraq.org/web"
)

// TestRendererLayouts ensures that primary templates render successfully and expose brand assets.
func TestRendererLayouts(t *testing.T) {
	r, err := web.NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	fixedTime := time.Date(2025, 10, 31, 9, 0, 0, 0, time.UTC)
	web.SetNowFunc(func() time.Time { return fixedTime })
	defer web.SetNowFunc(nil)

	data := webui.DashboardData{
		FrameworkName:        "Shanraq",
		FrameworkDescription: "Modular Go framework",
		PageTitle:            "Shanraq Console",
		CurrentYear:          2025,
		Environment:          "test",
		ReloadAfter:          5 * time.Second,
		SidebarLinks: []webui.SidebarLink{
			{Label: "Home", Href: "/"},
			{Label: "Console", Href: "/console"},
		},
		JobStats: webui.JobStats{
			Pending:      2,
			Running:      1,
			Completed:    5,
			Failed:       1,
			Retrying:     1,
			Total:        9,
			Workers:      4,
			PollInterval: "2s",
		},
		QueueOverview: webui.QueueOverview{
			DoneLastHour:       12,
			FailedLastHour:     3,
			SuccessRate:        0.8,
			FailureRate:        0.2,
			NextScheduled:      fixedTime.Add(30 * time.Minute),
			NextScheduledValid: true,
		},
		LastUpdated: fixedTime,
	}

	rec := httptest.NewRecorder()
	data.PageID = "home"
	if err := r.Render(rec, "home.html", data); err != nil {
		t.Fatalf("Render(home) error = %v", err)
	}

	body := rec.Body.String()
	assertContains(t, body, `id="home-root"`)
	assertContains(t, body, `/static/brand/logo.svg`)
	assertContains(t, body, `/static/brand/favicon.svg`)
	assertContains(t, body, `class="code-block`)
	assertNotContains(t, body, "No content provided.")
	assertNotContains(t, body, `/static/js/dashboard.js`)
	assertSnapshot(t, "home.html", body)

	rec = httptest.NewRecorder()
	data.PageID = "dashboard"
	if err := r.Render(rec, "dashboard.html", data); err != nil {
		t.Fatalf("Render(dashboard) error = %v", err)
	}

	body = rec.Body.String()
	assertContains(t, body, `id="dashboard-root"`)
	assertContains(t, body, "Queue Explorer")
	assertContains(t, body, `/static/js/dashboard.js`)
	assertContains(t, body, "Throughput (last hour)")
	assertContains(t, body, "Failure rate")
	assertNotContains(t, body, "No content provided.")
	assertSnapshot(t, "dashboard.html", body)

	rec = httptest.NewRecorder()
	data.PageID = "docs"
	data.PageTitle = "Framework Documentation"
	data.Docs = []webui.DocSection{
		{
			Title:   "Quick Start",
			Summary: "Kick off the reference binary.",
			Items: []webui.DocItem{
				{Title: "Run", Description: "Start locally", Command: "go run ./cmd/app -config config.yaml"},
				{Title: "Console", Description: "Open dashboard", Link: "/console"},
			},
		},
	}
	if err := r.Render(rec, "docs.html", data); err != nil {
		t.Fatalf("Render(docs) error = %v", err)
	}

	body = rec.Body.String()
	assertContains(t, body, "Framework Documentation")
	assertContains(t, body, "Quick Start")
	assertContains(t, body, "Start locally")
	assertContains(t, body, "go run ./cmd/app -config config.yaml")
	assertSnapshot(t, "docs.html", body)
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected output to contain %q; got: %s", needle, snippet(haystack))
	}
}

func assertNotContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("expected output not to contain %q; got: %s", needle, snippet(haystack))
	}
}

func snippet(s string) string {
	const max = 400
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func assertSnapshot(t *testing.T, name, got string) {
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		path := filepath.Join("testdata", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create snapshot dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write snapshot: %v", err)
		}
	}
	path := filepath.Join("testdata", name)
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot %s: %v", name, err)
	}
	if string(want) != got {
		t.Fatalf("snapshot mismatch for %s", name)
	}
}
