package web_test

import (
	"net/http/httptest"
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
		JobStats: webui.JobStats{},
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

	rec = httptest.NewRecorder()
	data.PageID = "dashboard"
	if err := r.Render(rec, "dashboard.html", data); err != nil {
		t.Fatalf("Render(dashboard) error = %v", err)
	}

	body = rec.Body.String()
	assertContains(t, body, `id="dashboard-root"`)
	assertContains(t, body, "Queue Explorer")
	assertContains(t, body, `/static/js/dashboard.js`)
	assertNotContains(t, body, "No content provided.")
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
