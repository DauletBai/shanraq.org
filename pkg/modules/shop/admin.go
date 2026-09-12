package shop

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/site"
)

// adminPage lists what the site sells, with every editable field of each
// product on the same page: there is one product today, and a separate list
// and form would be two clicks to change a price.
type adminPage struct {
	site.Base
	Items    []adminItem
	Statuses []string
	Notice   string
	Bad      bool
}

type adminItem struct {
	P     Product
	Leads int
	URL   string
}

func (m *Module) handleAdmin(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canEdit(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	items, err := m.store.List(r.Context(), true)
	if err != nil {
		m.logger().Error("shop admin list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := adminPage{
		Base:     m.rt.Site.Base(r, site.T(lang, "shop.admin_title"), lang),
		Statuses: Statuses,
	}
	page.NoIndex = true
	for _, p := range items {
		n, err := m.store.LeadCount(r.Context(), p.ID)
		if err != nil {
			m.logger().Warn("shop lead count", zap.Error(err))
		}
		page.Items = append(page.Items, adminItem{P: p, Leads: n, URL: "/shop/" + p.Slug})
	}
	switch r.URL.Query().Get("ok") {
	case "saved":
		page.Notice = site.T(lang, "shop.admin_saved")
	case "bad":
		page.Notice, page.Bad = site.T(lang, "shop.admin_bad"), true
	}
	m.render(w, "shop_admin", page)
}

// handleAdminSave writes one product back. What is sold and for how much is
// leadership's call, so an editor who can publish articles cannot change it.
func (m *Module) handleAdminSave(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canEdit(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	slug := chi.URLParam(r, "slug")
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
		return
	}
	p, err := m.store.BySlug(r.Context(), slug, true)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	status := strings.TrimSpace(r.FormValue("status"))
	if !validStatus(status) {
		http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
		return
	}
	// Each package has its own price, and a price of nothing means "not
	// announced yet" rather than "free".
	prices := map[string]int64{}
	for _, pl := range p.Plans {
		v, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("price_"+pl.Code)), 10, 64)
		if err != nil || v < 0 {
			http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
			return
		}
		prices[pl.Code] = v
	}
	// Selling something for nothing is a mistake nobody means to make, and it
	// would be made silently: the page would show a buy button and take no money.
	if status == StatusSelling && !anyPriced(prices) {
		http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
		return
	}
	p.Status = status
	p.Edition = strings.TrimSpace(r.FormValue("edition"))
	p.CoverURL = strings.TrimSpace(r.FormValue("cover_url"))
	p.PreviewURL = strings.TrimSpace(r.FormValue("preview_url"))
	for _, l := range site.Langs {
		p.Title[l] = strings.TrimSpace(r.FormValue("title_" + l))
		p.Summary[l] = strings.TrimSpace(r.FormValue("summary_" + l))
		p.Body[l] = r.FormValue("body_" + l)
	}
	if err := m.store.Save(r.Context(), p); err != nil {
		m.logger().Error("shop save", zap.Error(err))
		http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
		return
	}
	for _, pl := range p.Plans {
		pl.Price = prices[pl.Code]
		for _, l := range site.Langs {
			pl.Title[l] = strings.TrimSpace(r.FormValue("plan_title_" + pl.Code + "_" + l))
			pl.Includes[l] = r.FormValue("plan_includes_" + pl.Code + "_" + l)
		}
		if err := m.store.SavePlan(r.Context(), pl); err != nil {
			m.logger().Error("shop save plan", zap.Error(err))
			http.Redirect(w, r, "/admin/shop?ok=bad", http.StatusSeeOther)
			return
		}
	}
	m.logger().Info("shop product saved", zap.String("slug", slug), zap.String("status", status))
	http.Redirect(w, r, "/admin/shop?ok=saved", http.StatusSeeOther)
}

// anyPriced reports whether at least one package has a price on it.
func anyPriced(prices map[string]int64) bool {
	for _, v := range prices {
		if v > 0 {
			return true
		}
	}
	return false
}

func validStatus(s string) bool {
	for _, v := range Statuses {
		if v == s {
			return true
		}
	}
	return false
}
