package site

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"shanraq.org/web"
)

//go:embed templates/*.html
var templateFiles embed.FS

// Renderer is the site's single template set. Every module puts its own pages
// and helpers into it, so a page from one module can draw the header of
// another: the frame is defined once, in pkg/site, and the modules fill it.
//
// Parsing is deferred to Build, which the application calls after every module
// has had its turn. A template can then use any helper any module registered,
// whatever order the modules were registered in -- html/template resolves
// function names at parse time, so a template parsed before its module's
// helpers were known would fail for no reason a reader could see.
type Renderer struct {
	mu      sync.RWMutex
	funcs   template.FuncMap
	sources []templateSource
	tmpl    *template.Template
}

type templateSource struct {
	fsys     fs.FS
	patterns []string
}

// NewRenderer returns a renderer holding the site's own frame: the helpers in
// BaseFuncs and the shared partials (header, sidebar, footer).
func NewRenderer() *Renderer {
	r := &Renderer{funcs: template.FuncMap{}}
	r.Add(BaseFuncs(), templateFiles, "templates/*.html")
	return r
}

// Add contributes a module's helpers and templates. Template names are global:
// a module prefixes its own ("shop_book", not "book") so two modules cannot
// quietly define the same page.
func (r *Renderer) Add(funcs template.FuncMap, fsys fs.FS, patterns ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, fn := range funcs {
		r.funcs[name] = fn
	}
	if fsys != nil && len(patterns) > 0 {
		r.sources = append(r.sources, templateSource{fsys: fsys, patterns: patterns})
	}
}

// Build parses everything contributed so far. It is called once, at startup:
// a failure here is a broken template or a missing helper, which must stop the
// application rather than surface as a 500 on one page a month later.
func (r *Renderer) Build() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tmpl := template.New("shanraq").Funcs(r.funcs)
	for _, s := range r.sources {
		var err error
		if tmpl, err = tmpl.ParseFS(s.fsys, s.patterns...); err != nil {
			return err
		}
	}
	r.tmpl = tmpl
	return nil
}

// Execute writes one template to any writer. The error is returned rather than
// handled: the caller knows whether it is serving a page or filling a buffer.
func (r *Renderer) Execute(w io.Writer, name string, data any) error {
	r.mu.RLock()
	tmpl := r.tmpl
	r.mu.RUnlock()
	if tmpl == nil {
		return fmt.Errorf("site: templates not built (missing Build)")
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

// Render serves one template as an HTML page.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.Execute(w, name, data)
}

// Lookup reports whether a template of this name exists, for a caller that
// chooses a template by key and must not hand a stray name to Execute.
func (r *Renderer) Lookup(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tmpl != nil && r.tmpl.Lookup(name) != nil
}

// BaseFuncs are the helpers the frame itself needs, and the ones every module
// would otherwise write again: the dictionary, the languages, the icon set and
// the small formatting verbs. Anything that knows what a page is about belongs
// in that module's own map.
func BaseFuncs() template.FuncMap {
	return template.FuncMap{
		"t": T,
		// Every stylesheet and script must go through this: a bare path keeps
		// serving from cache for a day after a deploy, so new markup lands on
		// old CSS. See web.AssetURL.
		"asset":    web.AssetURL,
		"htmlLang": HTMLLang,
		"ogLocale": OGLocale,
		"langs":    func() []string { return Langs },
		"label":    func(l string) string { return LangLabels[l] },
		"langName": func(l string) string { return LangNames[l] },
		// The tabs and the translate button must name languages the same way,
		// or the reader has to work out that "3 языка" and "Қазақша" are about
		// the same thing.
		"langNames": func() map[string]string { return LangNames },
		"langList": func(ls []string) string {
			out := make([]string, 0, len(ls))
			for _, l := range ls {
				out = append(out, LangNames[l])
			}
			return strings.Join(out, ", ")
		},
		"liveSocial":  LiveSocial, // social profiles that aren't "#" placeholders
		"curSymbol":   CurSymbol,
		"svcOff":      ServiceOff, // is a service's entry link disabled?
		"svcMsg":      ServiceMsg, // its localized "unavailable" tooltip
		"categories":  func() []string { return Categories },
		"subcats":     func(cat string) []string { return Subcats(cat) },
		"icon":        Icon,
		"roomIcon":    RoomIcon,
		"amenityIcon": AmenityIcon,
		"catIcon":     CatIcon,
		// Brand mark for an analytics row, "" for rows that name no brand.
		"brandicon": BrandIcon,
		"withUTM":   WithUTM,
		"dict":      Dict,
		"firstN":    FirstStrings,
		// Templates count from zero and people count from one; screen-reader
		// labels are read by people.
		"inc": func(i int) int { return i + 1 },
		"mul": func(a, b int) int { return a * b },
		"sub": func(a, b int) int { return a - b },
		// Fills the {n} of a counted label. Not printf: the same string is read
		// by the browser, where a %d would be a stray placeholder in the markup.
		"count": func(s string, n int) string { return strings.ReplaceAll(s, "{n}", strconv.Itoa(n)) },
		// A cover that is a drawing rather than a photograph. The hero prints a
		// headline across the picture, and a diagram's own labels fight it.
		"isvector":  func(s string) bool { return strings.HasSuffix(strings.ToLower(s), ".svg") },
		"hasSuffix": strings.HasSuffix,
		"year":      func() int { return time.Now().Year() },
		"fmtDate": func(t time.Time) string {
			if t.IsZero() {
				return "—"
			}
			return t.Format("02.01.06")
		},
		"fmtDatePtr": func(t *time.Time) string {
			if t == nil || t.IsZero() {
				return "—"
			}
			return t.Format("02.01.06")
		},
	}
}

// WithUTM tags a share link with its channel, so a reader who passes an article
// on through WhatsApp arrives as WhatsApp rather than as "direct".
//
// This is not decoration. Messengers and in-app browsers strip the Referer, so
// without the tag every shared link lands in the direct bucket and the panel
// cannot tell a channel that works from one that does not — which is the whole
// question at this stage. utmSource maps the value back to a known label; a
// source it does not recognise is ignored rather than stored.
func WithUTM(rawURL, source string) string {
	if rawURL == "" || source == "" {
		return rawURL
	}
	sep := "?"
	if strings.Contains(rawURL, "?") {
		sep = "&"
	}
	return rawURL + sep + "utm_source=" + url.QueryEscape(source) + "&utm_medium=share"
}
