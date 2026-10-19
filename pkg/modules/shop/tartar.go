package shop

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"shanraq.org/pkg/site"
	"shanraq.org/web"
)

const (
	tartarPath     = "/book/tartar-i-ya"
	tartarCoverURL = "/static/shop/tartar-i-ya-sample/cover.jpg"
)

type tartarSection struct {
	Slug  string
	Title string
}

var tartarSections = []tartarSection{
	{Slug: "preface", Title: "Предисловие. Не выбрать сторону, а искать след"},
	{Slug: "prologue", Title: "Пролог. Как Казак стал Казахом"},
	{Slug: "chapter-1", Title: "Глава 1. Ру — Русь — Русский: как признак стал именем"},
	{Slug: "chapter-2", Title: "Глава 2. Бегство к свободной земле"},
	{Slug: "chapter-3", Title: "Глава 3. Бала — Лайық: имя инструмента как след"},
	{Slug: "chapter-4", Title: "Глава 4. Аткинсон: свидетель двора на шве переодевания"},
	{Slug: "chapter-5", Title: "Глава 5. Атаман, есаул и Арбат"},
}

type tartarPage struct {
	site.Base
	ChapterTitle string
	Content      template.HTML
	Sections     []tartarSection
	Previous     *tartarSection
	Next         *tartarSection
}

// The excerpt is extracted from the revised print layout and embedded as
// HTML fragments so each section uses the site's own page frame.
func loadTartarExcerpt() (map[string]template.HTML, error) {
	content := make(map[string]template.HTML, len(tartarSections))
	for _, section := range tartarSections {
		name := "shop/tartar-i-ya-sample/sections/" + section.Slug + ".html"
		data, err := fs.ReadFile(web.StaticFS(), name)
		if err != nil {
			return nil, fmt.Errorf("read Tartar excerpt %s: %w", name, err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("Tartar excerpt %s is empty", name)
		}
		content[section.Slug] = template.HTML(data) // trusted, embedded publication asset
	}
	return content, nil
}

func (m *Module) handleTartar(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	slug := chi.URLParam(r, "section")
	if slug == "why" {
		http.Redirect(w, r, tartarPath+"?lang="+lang, http.StatusMovedPermanently)
		return
	}
	page := tartarPage{
		Base:     m.rt.Site.Base(r, "Тартар и Я", lang),
		Sections: tartarSections,
	}
	page.Desc = site.T(lang, "tartar.sample_desc")
	page.OGImage = absolute(page.SiteURL, tartarCoverURL)
	if slug != "" {
		position := -1
		for i, section := range tartarSections {
			if section.Slug == slug {
				position = i
				page.ChapterTitle = section.Title
				page.Content = m.tartar[slug]
				break
			}
		}
		if position < 0 || page.Content == "" {
			http.NotFound(w, r)
			return
		}
		page.Title = page.ChapterTitle + " — Тартар и Я"
		if position == 0 {
			page.Previous = &tartarSection{Title: site.T(lang, "tartar.contents")}
		} else {
			page.Previous = &tartarSections[position-1]
		}
		if position+1 < len(tartarSections) {
			page.Next = &tartarSections[position+1]
		}
	}
	m.render(w, "tartar_sample", page)
}
