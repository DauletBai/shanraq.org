package articles

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsMaintenanceExempt(t *testing.T) {
	exempt := []string{"/studio/login", "/studio/logout", "/admin", "/admin/services", "/admin/roles"}
	for _, p := range exempt {
		if !isMaintenanceExempt(p) {
			t.Errorf("%s should be exempt (recovery route)", p)
		}
	}
	blocked := []string{"/", "/read/x", "/listings", "/advertise", "/studio", "/adminx", "/studio/new"}
	for _, p := range blocked {
		if isMaintenanceExempt(p) {
			t.Errorf("%s should NOT be exempt", p)
		}
	}
}

// TestMaintenanceGuard checks the two doors that prevent a lockout: when the
// site is down, an ordinary page returns 503, but the admin recovery route
// still passes through. No DB needed — the cache is set directly.
func TestMaintenanceGuard(t *testing.T) {
	m := &Module{tmpl: buildTemplates(t)}
	m.flags = &ServiceFlags{cache: map[string]ServiceFlag{
		SvcSite: {Code: SvcSite, Status: svcMaintenance, MessageRU: "Скоро вернёмся"},
	}}

	var reached bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true; w.WriteHeader(200) })
	guard := m.maintenanceGuard(next)

	// A public page is taken down.
	reached = false
	rec := httptest.NewRecorder()
	guard.ServeHTTP(rec, httptest.NewRequest("GET", "/listings", nil))
	if reached {
		t.Fatalf("public page should not reach the handler during maintenance")
	}
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Errorf("expected Retry-After header")
	}

	// The admin recovery route stays open, or we could never switch back.
	reached = false
	rec = httptest.NewRecorder()
	guard.ServeHTTP(rec, httptest.NewRequest("POST", "/admin/services", nil))
	if !reached {
		t.Fatalf("admin route must pass through during maintenance (anti-lockout)")
	}

	// When the site is up, everything passes.
	m.flags.cache[SvcSite] = ServiceFlag{Code: SvcSite, Status: svcOn}
	reached = false
	rec = httptest.NewRecorder()
	guard.ServeHTTP(rec, httptest.NewRequest("GET", "/listings", nil))
	if !reached {
		t.Fatalf("site up: page should pass through")
	}
}
