package articles

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// CoursesPage lists the published courses.
type CoursesPage struct {
	Base
	List []*Series
}

// CoursePage is one course's map: every lesson in order, with the reader's
// position in the whole thing visible at a glance.
type CoursePage struct {
	Base
	Series  *Series
	Items   []*SeriesItem
	Lessons int
	Minutes int
	// Done is how many of them this reader has passed; zero for a visitor.
	Done int
}

// handleCourses serves /courses.
func (m *Module) handleCourses(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	page := CoursesPage{Base: m.base(r, T(lang, "course.index_title"), lang)}
	page.Desc = T(lang, "course.index_lead")

	list, err := m.series.List(r.Context())
	if err != nil {
		m.rt.Logger.Error("series list", zap.Error(err))
	}
	page.List = list
	m.render(w, "courses", page)
}

// handleCourse serves /course/{slug} -- the map of one course.
//
// The whole course is visible from the first lesson on purpose: a reader who
// can see where the road ends decides to walk it, and one who is shown a single
// step at a time has nothing to decide about.
func (m *Module) handleCourse(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	slug := strings.TrimSpace(chi.URLParam(r, "slug"))

	s, err := m.series.BySlug(r.Context(), slug, lang)
	if err != nil {
		if errors.Is(err, ErrSeriesNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("series by slug", zap.String("slug", slug), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	page := CoursePage{
		Base:    m.base(r, s.TitleIn(lang), lang),
		Series:  s,
		Lessons: s.Lessons(),
		Minutes: s.Minutes(),
	}
	if d := s.SummaryIn(lang); d != "" {
		page.Desc = d
	}
	if s.CoverURL != "" {
		page.OGImage = s.CoverURL
	}
	// Unpublished lessons are dropped rather than greyed out: a course whose map
	// is half placeholders reads as abandoned, and the reader cannot tell an
	// empty row from a broken link.
	for _, it := range s.Items {
		if it.Published {
			page.Items = append(page.Items, it)
		}
	}

	// A reader who is signed in sees which lessons they have already been
	// through. Marked by the exercise being accepted, not by the page having
	// been opened: the map should show what was done, not what was visited.
	if uid, ok := m.authorID(r); ok && len(page.Items) > 0 {
		ids := make([]uuid.UUID, 0, len(page.Items))
		for _, it := range page.Items {
			ids = append(ids, it.ArticleID)
		}
		if done, err := m.progress.PassedIn(r.Context(), uid, ids); err == nil {
			for _, it := range page.Items {
				it.Passed = done[it.ArticleID]
				if it.Passed {
					page.Done++
				}
			}
		} else {
			m.rt.Logger.Warn("course progress", zap.Error(err))
		}
	}
	m.render(w, "course", page)
}

// ---- admin ----

type adminSeriesPage struct {
	Base
	Items    []*Series
	Editing  *Series
	Langs    []adminPageLangView
	Statuses []string
	Articles []predArticleOption
	Notice   string
	Error    string
}

// handleAdminSeries serves the course editor.
func (m *Module) handleAdminSeries(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	page := adminSeriesPage{
		Base:     m.base(r, T(lang, "course.admin_title"), lang),
		Statuses: []string{SeriesDraft, SeriesPublished},
	}
	for _, l := range pageEditLangs {
		page.Langs = append(page.Langs, adminPageLangView{Code: l.Code, Label: l.Label})
	}
	switch r.URL.Query().Get("saved") {
	case "1":
		page.Notice = T(lang, "course.saved")
	case "deleted":
		page.Notice = T(lang, "course.deleted")
	}
	if r.URL.Query().Get("err") != "" {
		page.Error = T(lang, "course.err_slug")
	}

	var err error
	if page.Items, err = m.series.ListAll(r.Context()); err != nil {
		m.rt.Logger.Error("admin series list", zap.Error(err))
	}
	// The picker offers drafts too: a lesson has to take its place in the course
	// before it goes live, or publishing it would drop it at the end.
	if arts, err := m.series.Candidates(r.Context(), lang, 300); err == nil {
		page.Articles = arts
	} else {
		m.rt.Logger.Error("series candidates", zap.Error(err))
	}
	if id, err := uuid.Parse(chi.URLParam(r, "id")); err == nil {
		if s, err := m.series.ByID(r.Context(), id, lang); err == nil {
			page.Editing = s
		} else if !errors.Is(err, ErrSeriesNotFound) {
			m.rt.Logger.Error("admin series get", zap.Error(err))
		}
	}
	m.render(w, "admin_series", page)
}

// handleAdminSeriesSave creates or updates a course.
func (m *Module) handleAdminSeriesSave(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	title, summary := map[string]string{}, map[string]string{}
	for _, l := range Langs {
		title[l] = r.FormValue("title_" + l)
		summary[l] = r.FormValue("summary_" + l)
	}
	var idp *uuid.UUID
	if id, err := uuid.Parse(r.FormValue("id")); err == nil {
		idp = &id
	}
	sid, err := m.series.Save(r.Context(), idp,
		r.FormValue("slug"), r.FormValue("cover_url"), r.FormValue("status"), title, summary)
	if err != nil {
		m.rt.Logger.Error("series save", zap.Error(err))
		http.Redirect(w, r, "/admin/series?err=1", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/series/"+sid.String()+"?saved=1", http.StatusSeeOther)
}

// handleAdminSeriesAttach adds a lesson to a course, or moves one already in it.
func (m *Module) handleAdminSeriesAttach(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	sid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	pos, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("position")))
	if err := m.series.AttachBySlug(r.Context(), sid, r.FormValue("article"), pos); err != nil {
		m.rt.Logger.Error("series attach", zap.Error(err))
	}
	http.Redirect(w, r, "/admin/series/"+sid.String()+"?saved=1", http.StatusSeeOther)
}

// handleAdminSeriesOrder saves every lesson's position in one submit, so a
// course can be re-ordered without one redirect per row.
func (m *Module) handleAdminSeriesOrder(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	sid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	for key, vals := range r.PostForm {
		rest, ok := strings.CutPrefix(key, "pos_")
		if !ok || len(vals) == 0 {
			continue
		}
		aid, err := uuid.Parse(rest)
		if err != nil {
			continue
		}
		pos, err := strconv.Atoi(strings.TrimSpace(vals[0]))
		if err != nil {
			continue
		}
		if err := m.series.Attach(r.Context(), sid, aid, pos); err != nil {
			m.rt.Logger.Error("series reorder", zap.Error(err))
		}
	}
	http.Redirect(w, r, "/admin/series/"+sid.String()+"?saved=1", http.StatusSeeOther)
}

// handleAdminSeriesDetach removes a lesson from a course, leaving the article.
func (m *Module) handleAdminSeriesDetach(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	sid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	if aid, err := uuid.Parse(chi.URLParam(r, "article")); err == nil {
		if err := m.series.Detach(r.Context(), sid, aid); err != nil {
			m.rt.Logger.Error("series detach", zap.Error(err))
		}
	}
	http.Redirect(w, r, "/admin/series/"+sid.String()+"?saved=1", http.StatusSeeOther)
}

// handleAdminSeriesDelete removes a course. Its lessons stay published.
func (m *Module) handleAdminSeriesDelete(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if id, err := uuid.Parse(chi.URLParam(r, "id")); err == nil {
		if err := m.series.Delete(r.Context(), id); err != nil {
			m.rt.Logger.Error("series delete", zap.Error(err))
		}
	}
	http.Redirect(w, r, "/admin/series?saved=deleted", http.StatusSeeOther)
}
