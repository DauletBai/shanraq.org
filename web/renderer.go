package web

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed views/*.html views/partials/*.html
var files embed.FS

// Renderer wraps html/template parsing with Bootstrap-ready layouts.
type Renderer struct {
	tmpl *template.Template
}

// NewRenderer parses embedded templates and returns a ready-to-use renderer.
func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"statusColor": func(status string) string {
			switch strings.ToLower(status) {
			case "pending":
				return "warning"
			case "running":
				return "primary"
			case "retry":
				return "info"
			case "failed":
				return "danger"
			case "done":
				return "success"
			default:
				return "secondary"
			}
		},
	}
	t, err := template.New("layout.html").Funcs(funcs).ParseFS(files, "views/*.html", "views/partials/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{tmpl: t}, nil
}

// Render writes the named template to the response writer.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.tmpl.ExecuteTemplate(w, name, data)
}

// Template exposes the parsed template tree.
func (r *Renderer) Template() *template.Template {
	return r.tmpl
}
