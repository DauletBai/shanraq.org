package articles

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// Account administration. The panel could previously report how many accounts
// held each role and nothing more — not one name, not one address — so a
// complaint about a person could not even be looked up, let alone acted on.

// adminActor returns the signed-in administrator's id, and whether they are
// allowed to manage accounts at all.
func (m *Module) adminActor(r *http.Request) (uuid.UUID, bool) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// adminUserAction resolves the target account and the acting administrator,
// writing the response itself when either is wrong.
func (m *Module) adminUserAction(w http.ResponseWriter, r *http.Request) (target, actor uuid.UUID, ok bool) {
	actor, allowed := m.adminActor(r)
	if !allowed {
		http.Error(w, "forbidden", http.StatusForbidden)
		return uuid.Nil, uuid.Nil, false
	}
	target, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return uuid.Nil, uuid.Nil, false
	}
	_ = r.ParseForm()
	return target, actor, true
}

// backToUsers returns to the account register, preserving the search that was
// open so a correction made from a filtered list does not throw the list away.
func backToUsers(r *http.Request, notice string) string {
	q := strings.TrimSpace(r.FormValue("q"))
	back := "/admin?ok=" + notice
	if q != "" {
		back += "&uq=" + url.QueryEscape(q)
	}
	return back + "#users"
}

// handleAdminUserUpdate corrects the name an account posts under, and whether
// its e-mail counts as verified.
func (m *Module) handleAdminUserUpdate(w http.ResponseWriter, r *http.Request) {
	target, _, ok := m.adminUserAction(w, r)
	if !ok {
		return
	}
	first := auth.NormalizePersonName(r.FormValue("first_name"))
	last := auth.NormalizePersonName(r.FormValue("last_name"))
	middle := auth.NormalizePersonName(r.FormValue("middle_name"))
	// The same rules the account holder faces. An administrator correcting a
	// typo must not be able to write something into a byline that registration
	// would have refused.
	if err := auth.ValidatePersonName(first); err != nil {
		http.Redirect(w, r, backToUsers(r, "user_bad_name"), http.StatusSeeOther)
		return
	}
	if err := auth.ValidatePersonName(last); err != nil {
		http.Redirect(w, r, backToUsers(r, "user_bad_name"), http.StatusSeeOther)
		return
	}
	if err := auth.ValidateOptionalPersonName(middle); err != nil {
		http.Redirect(w, r, backToUsers(r, "user_bad_name"), http.StatusSeeOther)
		return
	}
	verified := r.FormValue("verified") == "on"
	if err := m.users.UpdateUserByAdmin(r.Context(), target, first, last, middle, verified); err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("admin update user", zap.Error(err))
		http.Redirect(w, r, backToUsers(r, "user_failed"), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, backToUsers(r, "user_saved"), http.StatusSeeOther)
}

// handleAdminUserRole changes an account's role.
func (m *Module) handleAdminUserRole(w http.ResponseWriter, r *http.Request) {
	target, actor, ok := m.adminUserAction(w, r)
	if !ok {
		return
	}
	role := strings.ToLower(strings.TrimSpace(r.FormValue("role")))
	if !contains(assignableRoles, role) {
		http.Redirect(w, r, backToUsers(r, "user_failed"), http.StatusSeeOther)
		return
	}
	// Demoting yourself locks you out of the panel you are standing in, and the
	// only way back is a database console. Refuse rather than let one careless
	// dropdown end administrative access to the site.
	if target == actor && !contains([]string{"admin", "director"}, role) {
		http.Redirect(w, r, backToUsers(r, "user_self_demote"), http.StatusSeeOther)
		return
	}
	if err := m.users.SetRoleByID(r.Context(), target, role); err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			http.NotFound(w, r)
		case errors.Is(err, auth.ErrLastAdmin):
			http.Redirect(w, r, backToUsers(r, "user_last_admin"), http.StatusSeeOther)
		default:
			m.rt.Logger.Error("admin set role", zap.Error(err))
			http.Redirect(w, r, backToUsers(r, "user_failed"), http.StatusSeeOther)
		}
		return
	}
	http.Redirect(w, r, backToUsers(r, "role_set"), http.StatusSeeOther)
}

// handleAdminUserDelete removes an account and everything that cascades from
// it: articles, listings, comments, votes, bookmarks.
func (m *Module) handleAdminUserDelete(w http.ResponseWriter, r *http.Request) {
	target, actor, ok := m.adminUserAction(w, r)
	if !ok {
		return
	}
	// Deleting your own account from the admin panel is almost never what was
	// meant, and it takes the panel with it. The profile page has a delete that
	// asks for the password — that is where an administrator leaving should go.
	if target == actor {
		http.Redirect(w, r, backToUsers(r, "user_self_delete"), http.StatusSeeOther)
		return
	}
	if err := m.users.DeleteUserByAdmin(r.Context(), target); err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			http.NotFound(w, r)
		case errors.Is(err, auth.ErrLastAdmin):
			http.Redirect(w, r, backToUsers(r, "user_last_admin"), http.StatusSeeOther)
		default:
			m.rt.Logger.Error("admin delete user", zap.Error(err))
			http.Redirect(w, r, backToUsers(r, "user_failed"), http.StatusSeeOther)
		}
		return
	}
	http.Redirect(w, r, backToUsers(r, "user_deleted"), http.StatusSeeOther)
}
