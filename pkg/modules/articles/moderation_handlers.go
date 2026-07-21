package articles

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// MyModerationPage is the author's own record: every decision taken about
// their content, with the reason and the route to contest it.
type MyModerationPage struct {
	Base
	Actions []ModAction
	Saved   string
	Error   string
}

// handleMyModeration shows an author what was done to their content and why.
// Before this existed a comment could vanish with no explanation to anyone.
func (m *Module) handleMyModeration(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	page := MyModerationPage{Base: m.base(r, T(lang, "mod.my_title"), lang)}
	page.Saved = r.URL.Query().Get("ok")
	if q := r.URL.Query().Get("err"); q != "" {
		page.Error = T(lang, "mod.err_appeal")
	}
	acts, err := m.mods.ForSubject(r.Context(), uid, 100)
	if err != nil {
		m.rt.Logger.Error("my moderation", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page.Actions = acts
	m.render(w, "my_moderation", page)
}

// handleFileAppeal accepts an author's objection to one decision. The store
// checks ownership, so a forged action id cannot contest someone else's case.
func (m *Module) handleFileAppeal(w http.ResponseWriter, r *http.Request) {
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Redirect(w, r, "/studio/moderation?err=empty", http.StatusSeeOther)
		return
	}
	if err := m.mods.Appeal(r.Context(), id, uid, clip(body, 2000)); err != nil {
		// Not an internal error: the usual cause is a second appeal on the same
		// decision, or an action that is not appealable at all.
		m.rt.Logger.Info("appeal rejected", zap.Error(err))
		http.Redirect(w, r, "/studio/moderation?err=1", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/studio/moderation?ok=appeal", http.StatusSeeOther)
}

// handleAdminResolveAppeal closes an appeal. Overturning restores the content
// and writes the reversal into the ledger as its own entry.
func (m *Module) handleAdminResolveAppeal(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canModerate(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	mid, _ := uuid.Parse(claims.Subject)
	uphold := r.FormValue("decision") == "uphold"
	if err := m.mods.Resolve(r.Context(), id, mid, uphold, clip(strings.TrimSpace(r.FormValue("note")), 1000)); err != nil {
		m.rt.Logger.Error("resolve appeal", zap.Error(err))
		http.Redirect(w, r, "/admin?err=appeal", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin?ok=appeal_resolved", http.StatusSeeOther)
}
