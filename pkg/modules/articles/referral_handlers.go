package articles

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InvitePage backs the referral cabinet page: the user's invite link, how many
// they have brought in, and their promotion-day credit.
type InvitePage struct {
	Base
	Stats  ReferralStats
	Link   string
	Reward int
	Saved  string
	Error  string
}

// handleInvite shows the author their invite link and referral progress.
func (m *Module) handleInvite(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	stats, err := m.refs.Stats(r.Context(), uid)
	if err != nil {
		m.rt.Logger.Error("referral stats", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := InvitePage{Base: m.base(r, T(lang, "inv.title"), lang)}
	page.Stats = stats
	page.Reward = referralRewardDays
	page.Link = m.rt.Config.PublicBase() + "/studio/register?ref=" + stats.Code
	page.Saved = r.URL.Query().Get("ok")
	if r.URL.Query().Get("err") == "credit" {
		page.Error = T(lang, "inv.err_credit")
	}
	m.render(w, "invite", page)
}

// handleListingPromoteFree spends referral credit to promote a listing instead
// of paying. The store re-checks the balance inside its transaction, so this
// cannot overspend even under concurrent clicks.
func (m *Module) handleListingPromoteFree(w http.ResponseWriter, r *http.Request) {
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
	if err := m.refs.SpendCredit(r.Context(), uid, id, referralRewardDays); err != nil {
		if errors.Is(err, ErrInsufficientCredit) {
			http.Redirect(w, r, "/studio/invite?err=credit", http.StatusSeeOther)
			return
		}
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("promote free", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/listings/my?ok=promoted", http.StatusSeeOther)
}
