package articles

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// PredictionsPage is the public ledger.
type PredictionsPage struct {
	Base
	Score PredictionScore
	Open  []*Prediction
	Done  []*Prediction
}

// handlePredictions serves /predictions: the scoreboard and the whole ledger.
//
// Open forecasts sit above settled ones because they are the only part a reader
// can still argue with — the settled half is the evidence that the open half is
// worth reading.
func (m *Module) handlePredictions(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	page := PredictionsPage{Base: m.base(r, T(lang, "pred.title"), lang)}
	page.Desc = T(lang, "pred.lead")

	list, err := m.predictions.List(r.Context(), lang)
	if err != nil {
		m.rt.Logger.Error("predictions list", zap.Error(err))
	}
	for _, p := range list {
		if p.Resolved() {
			page.Done = append(page.Done, p)
		} else {
			page.Open = append(page.Open, p)
		}
	}
	if sc, err := m.predictions.Score(r.Context()); err == nil {
		page.Score = sc
	} else {
		m.rt.Logger.Error("predictions score", zap.Error(err))
	}
	m.render(w, "predictions", page)
}

// adminPredictionsPage is the management screen: the whole ledger plus one
// editor. Editing happens on the same page as the list so an operator resolving
// a forecast can see the others without navigating away — the common case is
// "what is overdue", not "open this one row".
type adminPredictionsPage struct {
	Base
	Score    PredictionScore
	Items    []*Prediction
	Editing  *Prediction
	Langs    []adminPageLangView // reused for the language tab labels
	Statuses []string
	Articles []predArticleOption
	Notice   string
	Error    string
}

type predArticleOption struct {
	ID    string
	Title string
}

func (m *Module) handleAdminPredictions(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	page := adminPredictionsPage{
		Base:     m.base(r, T(lang, "pred.admin_title"), lang),
		Statuses: PredStatuses,
	}
	for _, l := range pageEditLangs {
		page.Langs = append(page.Langs, adminPageLangView{Code: l.Code, Label: l.Label})
	}
	switch r.URL.Query().Get("saved") {
	case "1":
		page.Notice = T(lang, "pred.saved")
	case "deleted":
		page.Notice = T(lang, "pred.deleted")
	}
	if e := r.URL.Query().Get("err"); e != "" {
		page.Error = T(lang, "pred.err_empty")
	}

	var err error
	if page.Items, err = m.predictions.List(r.Context(), lang); err != nil {
		m.rt.Logger.Error("admin predictions list", zap.Error(err))
	}
	if page.Score, err = m.predictions.Score(r.Context()); err != nil {
		m.rt.Logger.Error("admin predictions score", zap.Error(err))
	}
	// The article picker offers the pieces a forecast can be attached to.
	if arts, err := m.store.LLMSIndex(r.Context(), lang); err == nil {
		for _, a := range arts {
			page.Articles = append(page.Articles, predArticleOption{ID: a.Slug, Title: a.Title})
		}
	}
	if id, err := uuid.Parse(chi.URLParam(r, "id")); err == nil {
		if p, err := m.predictions.Get(r.Context(), lang, id); err == nil {
			page.Editing = p
		} else if !errors.Is(err, pgx.ErrNoRows) {
			m.rt.Logger.Error("admin prediction get", zap.Error(err))
		}
	}
	m.render(w, "admin_predictions", page)
}

// parseDate reads a date from an <input type="date">, returning nil for the
// empty string so "no deadline" survives the round trip.
func parseDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

func (m *Module) handleAdminPredictionSave(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	in := PredictionInput{
		Status:    r.FormValue("status"),
		SourceURL: strings.TrimSpace(r.FormValue("source_url")),
		Horizon:   parseDate(r.FormValue("horizon")),
		Statement: map[string]string{},
		Verdict:   map[string]string{},
	}
	if d := parseDate(r.FormValue("made_on")); d != nil {
		in.MadeOn = *d
	}
	in.ResolvedOn = parseDate(r.FormValue("resolved_on"))
	for _, l := range Langs {
		in.Statement[l] = r.FormValue("statement_" + l)
		in.Verdict[l] = r.FormValue("verdict_" + l)
	}
	// The picker submits a slug, because that is what an operator recognises.
	if slug := strings.TrimSpace(r.FormValue("article")); slug != "" {
		if a, err := m.store.GetPublishedBySlug(r.Context(), slug); err == nil && a != nil {
			id := a.ID
			in.ArticleID = &id
		}
	}
	id := uuid.Nil
	if parsed, err := uuid.Parse(r.FormValue("id")); err == nil {
		id = parsed
	}
	if _, err := m.predictions.Save(r.Context(), id, in); err != nil {
		if errors.Is(err, ErrPredictionEmpty) {
			http.Redirect(w, r, "/admin/predictions?err=empty&lang="+lang, http.StatusSeeOther)
			return
		}
		m.rt.Logger.Error("save prediction", zap.Error(err))
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/predictions?saved=1&lang="+lang, http.StatusSeeOther)
}

func (m *Module) handleAdminPredictionDelete(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := m.predictions.Delete(r.Context(), id); err != nil {
		m.rt.Logger.Error("delete prediction", zap.Error(err))
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/predictions?saved=deleted&lang="+lang, http.StatusSeeOther)
}
