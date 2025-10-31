package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
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
		"percent":      formatPercent,
		"relativeTime": relativeTime,
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

func formatPercent(value float64) string {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return fmt.Sprintf("%.0f%%", value*100)
}

func relativeTime(t time.Time) string {
	if t.IsZero() {
		return "Not scheduled"
	}
	delta := time.Until(t)
	abs := delta
	if delta < 0 {
		abs = -delta
	}
	if abs < 10*time.Second {
		return "just now"
	}
	var unit string
	var value int
	switch {
	case abs < time.Minute:
		unit = "second"
		value = int(abs.Seconds())
	case abs < time.Hour:
		unit = "minute"
		value = int(abs.Minutes())
	case abs < 24*time.Hour:
		unit = "hour"
		value = int(abs.Hours())
	default:
		unit = "day"
		value = int(abs.Hours()/24 + 0.5)
	}
	if value > 1 {
		unit += "s"
	}
	phrase := fmt.Sprintf("%d %s", value, unit)
	if delta < 0 {
		return phrase + " ago"
	}
	return "in " + phrase
}
