package shop

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

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
	{Slug: "why", Title: "Почему книга называется «Тартар и Я»"},
	{Slug: "preface", Title: "Предисловие. Не выбрать сторону, а искать след"},
	{Slug: "chapter-1", Title: "Глава 1. Как читать слово"},
	{Slug: "chapter-2", Title: "Глава 2. «Будь»"},
}

type tartarPage struct {
	site.Base
	ChapterTitle string
	Content      template.HTML
	Sections     []tartarSection
	Previous     *tartarSection
	Next         *tartarSection
}

// Keep one source of the excerpt's text: the already published, embedded HTML.
// The wrapper and its relative navigation are replaced with the site's frame.
func loadTartarExcerpt() (map[string]template.HTML, error) {
	content := make(map[string]template.HTML, len(tartarSections))
	for _, section := range tartarSections {
		name := "shop/tartar-i-ya-sample/" + section.Slug + ".html"
		data, err := fs.ReadFile(web.StaticFS(), name)
		if err != nil {
			return nil, fmt.Errorf("read Tartar excerpt %s: %w", name, err)
		}
		_, body, ok := strings.Cut(string(data), `<article class="reading">`)
		if !ok {
			return nil, fmt.Errorf("Tartar excerpt %s has no reading article", name)
		}
		body, _, ok = strings.Cut(body, "</article>")
		if !ok {
			return nil, fmt.Errorf("Tartar excerpt %s has no article end", name)
		}
		content[section.Slug] = template.HTML(body) // trusted, embedded publication asset
	}
	return content, nil
}

func (m *Module) handleTartar(w http.ResponseWriter, r *http.Request) {
	lang := site.ResolveLang(w, r)
	slug := chi.URLParam(r, "section")
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
