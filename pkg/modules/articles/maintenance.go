package articles

import (
	"bytes"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// MaintenancePage backs the global maintenance screen.
type MaintenancePage struct {
	Lang    string
	Title   string
	Message string
}

// maintenanceGuard serves a maintenance page for the whole site when the global
// switch is not 'on'. Two doors stay open so we never lock ourselves out: the
// recovery routes (login/admin) and any staff session — staff can keep browsing
// to verify a change before flipping the site back on.
func (m *Module) maintenanceGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.flags == nil || m.flags.SiteUp() {
			next.ServeHTTP(w, r)
			return
		}
		if isMaintenanceExempt(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if claims, ok := auth.ClaimsFromContext(r.Context()); ok && claims.HasAnyRole(adminRoles...) {
			next.ServeHTTP(w, r)
			return
		}
		m.renderMaintenance(w, r)
	})
}

// isMaintenanceExempt keeps the recovery surface reachable during a global
// takedown: the admin panel (to switch back) and the login/logout it needs.
func isMaintenanceExempt(p string) bool {
	if p == "/studio/login" || p == "/studio/logout" {
		return true
	}
	return p == "/admin" || strings.HasPrefix(p, "/admin/")
}

func (m *Module) renderMaintenance(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	f := m.flags.SiteFlag()
	data := MaintenancePage{Lang: lang, Title: T(lang, "svc.site_down_title"), Message: f.Message(lang)}
	if data.Message == "" {
		data.Message = T(lang, "svc.site_down_default")
	}
	// Render to a buffer first so a template error does not leave us with a
	// half-written body under a 503 status.
	var buf bytes.Buffer
	if err := m.tmpl.ExecuteTemplate(&buf, "maintenance", data); err != nil {
		m.rt.Logger.Error("render maintenance", zap.Error(err))
		http.Error(w, "Site under maintenance", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Retry-After", "3600")
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = buf.WriteTo(w)
}
