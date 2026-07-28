package articles

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// The admin page editor lets project leadership update the info & legal pages
// (Privacy, Terms, Pricing, …) without a code change. The public reader serves
// from content_pages with a built-in fallback; these handlers write to it.

type adminPagesList struct {
	Base
	Items []adminPageItem
}

type adminPageItem struct{ Key, Name string }

type adminPageEditView struct {
	Base
	Key        string
	Name       string
	Notice     string
	Error      string
	LastEdited string
	LastEditor string
	Langs      []adminPageLangView
}

type adminPageLangView struct{ Code, Label, Title, Body string }

var pageEditLangs = []struct{ Code, Label string }{
	{"kz", "Қазақша"},
	{"ru", "Русский"},
	{"en", "English"},
}

func isEditablePage(key string) bool {
	for _, k := range editablePageKeys {
		if k == key {
			return true
		}
	}
	return false
}

// editablePageName is a page's effective title in the admin's UI language — used
// for the list and the editor heading.
func (m *Module) editablePageName(r *http.Request, key, lang string) string {
	t, _ := m.pageContent(r.Context(), key, lang)
	return t
}

func (m *Module) handleAdminPages(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	page := adminPagesList{Base: m.base(r, T(lang, "pages.title"), lang)}
	for _, key := range editablePageKeys {
		page.Items = append(page.Items, adminPageItem{Key: key, Name: m.editablePageName(r, key, lang)})
	}
	m.render(w, "admin_pages", page)
}

func (m *Module) handleAdminPageEdit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	key := chi.URLParam(r, "key")
	if !isEditablePage(key) {
		http.NotFound(w, r)
		return
	}
	view := adminPageEditView{
		Base: m.base(r, T(lang, "pages.edit_title"), lang),
		Key:  key,
		Name: m.editablePageName(r, key, lang),
	}
	if r.URL.Query().Get("ok") == "1" {
		view.Notice = T(lang, "pages.saved")
	}
	if when, by, ok, _ := m.content.LastEdited(r.Context(), key); ok && !when.IsZero() {
		view.LastEdited = when.Format("2006-01-02 15:04")
		view.LastEditor = by
	}
	for _, l := range pageEditLangs {
		t, b := m.pageContent(r.Context(), key, l.Code)
		view.Langs = append(view.Langs, adminPageLangView{Code: l.Code, Label: l.Label, Title: t, Body: b})
	}
	m.render(w, "admin_page_edit", view)
}

func (m *Module) handleAdminPageSave(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	key := chi.URLParam(r, "key")
	if !isEditablePage(key) {
		http.NotFound(w, r)
		return
	}
	editor, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login?reason=session_expired", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	lang := m.resolveLang(w, r)
	// Collect every language first so validation can reject the whole save: a
	// missing title or body must never overwrite a live legal page with a blank.
	langs := make([]adminPageLangView, 0, len(pageEditLangs))
	valid := true
	for _, l := range pageEditLangs {
		lv := adminPageLangView{
			Code:  l.Code,
			Label: l.Label,
			Title: strings.TrimSpace(r.FormValue("title_" + l.Code)),
			Body:  r.FormValue("body_" + l.Code),
		}
		if lv.Title == "" || strings.TrimSpace(lv.Body) == "" {
			valid = false
		}
		langs = append(langs, lv)
	}
	if !valid {
		// Re-render with what the operator typed plus an error — nothing is lost
		// and a blank page can't be saved.
		m.render(w, "admin_page_edit", adminPageEditView{
			Base:  m.base(r, T(lang, "pages.edit_title"), lang),
			Key:   key,
			Name:  m.editablePageName(r, key, lang),
			Error: T(lang, "pages.err_required"),
			Langs: langs,
		})
		return
	}
	if err := m.content.UpsertMany(r.Context(), key, langs, editor); err != nil {
		m.rt.Logger.Error("save content page", zap.String("key", key), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	m.rt.Logger.Info("content page edited", zap.String("key", key), zap.String("editor", editor.String()))
	http.Redirect(w, r, "/admin/pages/"+key+"?ok=1", http.StatusSeeOther)
}
