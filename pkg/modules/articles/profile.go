package articles

import (
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
)

// ProfilePage backs the cabinet profile/settings screen: avatar, identity, and
// the account-deletion danger zone.
type ProfilePage struct {
	Base
	Email         string
	FirstName     string
	LastName      string
	PhoneVerified bool
	EmailVerified bool
	Notice        string
}

// handleProfile renders the user's profile & settings page.
func (m *Module) handleProfile(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	lang := m.resolveLang(w, r)
	page := ProfilePage{Base: m.base(r, T(lang, "prof.title"), lang)}
	page.Notice = noticeText(lang, r.URL.Query().Get("ok"))

	page.FirstName, page.LastName, page.PhoneVerified = m.auth.AuthorIdentity(r.Context(), authorID)
	page.EmailVerified = m.auth.IsEmailVerified(r.Context(), authorID)
	if claims, ok := auth.ClaimsFromContext(r.Context()); ok {
		page.Email = claims.Email
	}
	m.render(w, "studio_profile", page)
}

// handleAvatarUpload accepts a profile photo, processes it (square, EXIF-stripped,
// no watermark), stores it, and points the user's avatar at it.
func (m *Module) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if m.media == nil {
		http.Redirect(w, r, "/studio/profile?ok=avatar_bad", http.StatusSeeOther)
		return
	}
	limit := m.media.MaxUploadBytes()
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	if err := r.ParseMultipartForm(limit); err != nil {
		http.Redirect(w, r, "/studio/profile?ok=avatar_big", http.StatusSeeOther)
		return
	}
	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Redirect(w, r, "/studio/profile?ok=avatar_bad", http.StatusSeeOther)
		return
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		http.Redirect(w, r, "/studio/profile?ok=avatar_bad", http.StatusSeeOther)
		return
	}
	url, err := m.media.ProcessAndSaveAvatar(r.Context(), raw)
	if err != nil {
		m.rt.Logger.Warn("avatar process", zap.Error(err))
		http.Redirect(w, r, "/studio/profile?ok=avatar_bad", http.StatusSeeOther)
		return
	}
	if err := m.auth.SetAvatar(r.Context(), authorID, url); err != nil {
		m.rt.Logger.Error("set avatar", zap.Error(err))
		http.Redirect(w, r, "/studio/profile?ok=avatar_bad", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/studio/profile?ok=avatar_set", http.StatusSeeOther)
}

// handleAvatarDelete removes the user's avatar (reverting to the default).
func (m *Module) handleAvatarDelete(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := m.auth.SetAvatar(r.Context(), authorID, ""); err != nil {
		m.rt.Logger.Error("clear avatar", zap.Error(err))
	}
	http.Redirect(w, r, "/studio/profile?ok=avatar_cleared", http.StatusSeeOther)
}

// handleDeleteAccount permanently erases the user's account and all their
// content after confirming their password. This is the user's right-to-erasure
// action; it clears the session and returns to the home page.
func (m *Module) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	// Confirm intent: the correct password must be re-entered.
	if !m.auth.CheckPassword(r.Context(), authorID, strings.TrimSpace(r.FormValue("password"))) {
		http.Redirect(w, r, "/studio/profile?ok=del_badpass", http.StatusSeeOther)
		return
	}
	if err := m.auth.DeleteAccount(r.Context(), authorID); err != nil {
		m.rt.Logger.Error("delete account", zap.Error(err))
		http.Redirect(w, r, "/studio/profile?ok=del_failed", http.StatusSeeOther)
		return
	}
	auth.ClearSessionCookie(w, r)
	http.Redirect(w, r, "/?ok=account_deleted", http.StatusSeeOther)
}

// noticeText maps an ?ok= code to a localized message for the profile page.
func noticeText(lang, code string) string {
	if code == "" {
		return ""
	}
	msg := T(lang, "prof."+code)
	if strings.HasPrefix(msg, "prof.") {
		return "" // unknown code
	}
	return msg
}
