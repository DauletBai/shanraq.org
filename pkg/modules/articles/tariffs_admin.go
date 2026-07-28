package articles

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// The admin tariff editor lets leadership retune the ad rate card, listing
// service prices and durations without a redeploy.

type tariffCell struct {
	Key   string
	Value int64
}

type tariffAdRow struct {
	Label string // format size, e.g. "970×250"
	Cells []tariffCell
}

type tariffField struct {
	Key   string
	Label string
	Value int64
}

type tariffsAdminView struct {
	Base
	Notice    string
	AdDays    []int
	AdFormats []tariffAdRow
	Weights   []tariffField
	Banner    []tariffField
	Services  []tariffField
	Durations []tariffField
}

func (m *Module) tariffsAdminView(lang string) tariffsAdminView {
	v := tariffsAdminView{AdDays: adDurationDays}
	for _, f := range AdFormats() {
		row := tariffAdRow{Label: f.Size}
		for _, d := range adDurationDays {
			key := "ad." + f.Code + "." + strconv.Itoa(d)
			row.Cells = append(row.Cells, tariffCell{Key: key, Value: tariffVal(key)})
		}
		v.AdFormats = append(v.AdFormats, row)
	}
	v.Weights = []tariffField{
		{"weight.high", T(lang, "tar.w_high"), tariffVal("weight.high")},
		{"weight.mid", T(lang, "tar.w_mid"), tariffVal("weight.mid")},
		{"weight.base", T(lang, "tar.w_base"), tariffVal("weight.base")},
	}
	for _, d := range BannerDays() {
		key := "banner." + strconv.Itoa(d)
		v.Banner = append(v.Banner, tariffField{Key: key, Label: strconv.Itoa(d), Value: tariffVal(key)})
	}
	v.Services = []tariffField{
		{"promote.price", T(lang, "tar.promote"), tariffVal("promote.price")},
		{"feature.price", T(lang, "tar.feature"), tariffVal("feature.price")},
	}
	v.Durations = []tariffField{
		{"listing.free_days", T(lang, "tar.free_days"), tariffVal("listing.free_days")},
		{"promote.days", T(lang, "tar.promote_days"), tariffVal("promote.days")},
		{"feature.days", T(lang, "tar.feature_days"), tariffVal("feature.days")},
	}
	return v
}

func (m *Module) handleAdminTariffs(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	view := m.tariffsAdminView(lang)
	view.Base = m.base(r, T(lang, "tar.title"), lang)
	if r.URL.Query().Get("ok") == "1" {
		view.Notice = T(lang, "tar.saved")
	}
	m.render(w, "admin_tariffs", view)
}

func (m *Module) handleAdminTariffsSave(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if m.tariffs == nil {
		http.Redirect(w, r, "/admin/tariffs", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	values := make(map[string]int64, len(tariffDefs))
	for _, d := range tariffDefs {
		raw := strings.TrimSpace(r.FormValue(d.key))
		if raw == "" {
			continue
		}
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			continue // ignore a non-numeric field rather than fail the whole save
		}
		values[d.key] = n
	}
	var by *uuid.UUID
	if claims != nil {
		if id, err := uuid.Parse(claims.Subject); err == nil {
			by = &id
		}
	}
	if err := m.tariffs.SaveMany(r.Context(), values, by); err != nil {
		m.rt.Logger.Warn("save tariffs", zap.Error(err))
	}
	m.rt.Logger.Info("tariffs updated", zap.Int("fields", len(values)))
	http.Redirect(w, r, "/admin/tariffs?ok=1", http.StatusSeeOther)
}
