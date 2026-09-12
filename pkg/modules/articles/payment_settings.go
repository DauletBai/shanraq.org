package articles

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/payments"
)

// handleAdminPayments switches the active acquirer and payments on/off from the
// admin panel. Secret keys are never touched here — they live in the server
// config; this only records the choice, taking effect at once.
func (m *Module) handleAdminPayments(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if m.payments == nil {
		http.Redirect(w, r, "/admin?ok=pay_bad", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		// A truncated/corrupt POST must not be read as "enabled=false" and quietly
		// switch payments off.
		http.Redirect(w, r, "/admin?ok=pay_bad", http.StatusSeeOther)
		return
	}
	in := payments.Settings{
		Enabled:  r.FormValue("enabled") == "on",
		Provider: strings.TrimSpace(r.FormValue("provider")),
	}
	var by *uuid.UUID
	if claims != nil {
		if id, err := uuid.Parse(claims.Subject); err == nil {
			by = &id
		}
	}
	if err := m.payments.SaveAdmin(r.Context(), in, by); err != nil {
		m.rt.Logger.Warn("update payment settings", zap.Error(err))
		http.Redirect(w, r, "/admin?ok=pay_bad", http.StatusSeeOther)
		return
	}
	m.rt.Logger.Info("payment settings updated", zap.Bool("enabled", in.Enabled), zap.String("provider", in.Provider))
	http.Redirect(w, r, "/admin?ok=pay_set", http.StatusSeeOther)
}
