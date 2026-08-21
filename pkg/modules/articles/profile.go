package articles

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/media"
)

// ProfilePage backs the cabinet profile/settings screen: avatar, identity, and
// the account-deletion danger zone.
type ProfilePage struct {
	Base
	Email         string
	FirstName     string
	LastName      string
	BioKZ         string
	BioRU         string
	BioEN         string
	CanPublish    bool // leadership: publish without email/phone verification
	PhoneVerified bool
	EmailVerified bool
	Notice        string

	// PlaceID and PlaceLabel are where the reader says they live. Empty means
	// they never said, and then everything is shown to them — which is what
	// the site did before places existed.
	PlaceID    string
	PlaceLabel string
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
	card := m.auth.AuthorCard(r.Context(), authorID)
	page.BioKZ, page.BioRU, page.BioEN = card.BioKZ, card.BioRU, card.BioEN
	if claims, ok := auth.ClaimsFromContext(r.Context()); ok {
		page.Email = claims.Email
		page.CanPublish = canAuthorAsStaff(claims)
	}
	if node, err := m.geo.UserPlace(r.Context(), authorID); err != nil {
		m.rt.Logger.Warn("read user place", zap.Error(err))
	} else if node != nil {
		page.PlaceID = node.String()
		if label, lerr := m.geo.PlaceLabel(r.Context(), *node, lang); lerr == nil {
			page.PlaceLabel = label
		}
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
	url, err := m.media.ProcessAndSaveAvatar(r.Context(), authorID, raw)
	if err != nil {
		// Being out of room is not a broken file, and telling someone their
		// photo is unreadable when it is their quota that is full sends them to
		// re-export it forever.
		if errors.Is(err, media.ErrQuotaExceeded) || errors.Is(err, media.ErrStoreFull) {
			http.Redirect(w, r, "/studio/profile?ok=avatar_quota", http.StatusSeeOther)
			return
		}
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

// handleBioSave stores the user's public author bio (shown on their author page).
func (m *Module) handleBioSave(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	kz := clip(strings.TrimSpace(r.FormValue("bio_kz")), 600)
	ru := clip(strings.TrimSpace(r.FormValue("bio_ru")), 600)
	en := clip(strings.TrimSpace(r.FormValue("bio_en")), 600)
	if err := m.auth.SetBios(r.Context(), authorID, kz, ru, en); err != nil {
		m.rt.Logger.Error("set bios", zap.Error(err))
		http.Redirect(w, r, "/studio/profile?ok=bio_bad", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/studio/profile?ok=bio_set", http.StatusSeeOther)
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
		// The site must not be left without anyone who can administer it, and
		// deleting your own account is the one route to that state the console's
		// own guard never sees. Say why, rather than reporting a generic failure
		// the administrator would reasonably read as a bug.
		if errors.Is(err, auth.ErrLastAdmin) {
			http.Redirect(w, r, "/studio/profile?ok=del_lastadmin", http.StatusSeeOther)
			return
		}
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

// handleProfilePlace saves where the reader says they live, or clears it when
// they empty the field.
//
// A place chosen at registration and then unchangeable would be a trap: people
// move, and somebody who picked the wrong line in a hurry must be able to fix
// it without writing to support.
func (m *Module) handleProfilePlace(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	raw := strings.TrimSpace(r.FormValue("geo_node_id"))
	var node *uuid.UUID
	if raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			http.Redirect(w, r, "/studio/profile?ok=place_bad", http.StatusSeeOther)
			return
		}
		node = &id
	}
	if err := m.geo.SetUserPlace(r.Context(), authorID, node); err != nil {
		m.rt.Logger.Error("save user place", zap.Error(err))
		http.Redirect(w, r, "/studio/profile?ok=place_bad", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/studio/profile?ok=place_saved", http.StatusSeeOther)
}
