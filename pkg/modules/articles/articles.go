package articles

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/modules/ratings"
	"shanraq.org/pkg/modules/syndicate"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/validate"
)

//go:embed templates/*.html
var templateFiles embed.FS

// Module implements the article publishing product surface: a public reader
// (feed + article pages with a language switcher) and an authenticated author
// cabinet ("studio") with a dashboard and a trilingual editor.
// Mailer sends e-mail (satisfied by the notifier module). Used for listing
// expiry reminders.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type Module struct {
	rt        *shanraq.Runtime
	store     *Store
	listings  *ListingStore
	geo       *GeoStore
	comments  *CommentStore
	favs      *FavoriteStore
	admin     *AdminStore
	ratings   *ratings.Store
	jobs      *jobs.Store
	auth      *auth.Module
	ai        *ai.Module
	syndicate *syndicate.Module
	mailer    Mailer
	tmpl      *template.Template
	validator *validate.Validator
}

// New builds the articles module. It depends on auth (browser sessions), ai
// (writing assistant), syndicate (publish → external channels), and a mailer
// (listing expiry reminders).
func New(authModule *auth.Module, aiModule *ai.Module, syndicateModule *syndicate.Module, mailer Mailer) *Module {
	return &Module{auth: authModule, ai: aiModule, syndicate: syndicateModule, mailer: mailer}
}

func (m *Module) Name() string { return "articles" }

// Init parses templates and wires the store.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	m.listings = NewListingStore(rt.DB)
	m.geo = NewGeoStore(rt.DB)
	m.comments = NewCommentStore(rt.DB)
	m.favs = NewFavoriteStore(rt.DB)
	m.admin = NewAdminStore(rt.DB)
	m.ratings = ratings.NewStore(rt.DB)
	m.jobs = jobs.NewStore(rt.DB)
	m.validator = validate.New()

	funcs := template.FuncMap{
		"t":                T,
		"label":            func(l string) string { return LangLabels[l] },
		"langName":         func(l string) string { return LangNames[l] },
		"langs":            func() []string { return Langs },
		"categories":       func() []string { return Categories },
		"editorCategories": func() []string { return append([]string{CategoryGeneral}, Categories...) },
		"subcats":          func(cat string) []string { return Subcats(cat) },
		"dealTypes":        func() []string { return DealTypes },
		"propertyTypes":    func() []string { return PropertyTypes },
		"money":            money,
		"ogLocale":         ogLocale,
		"year":             func() int { return time.Now().Year() },
		"markdown":         RenderMarkdown,
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
	tmpl, err := template.New("articles").Funcs(funcs).ParseFS(templateFiles, "templates/*.html")
	if err != nil {
		return err
	}
	m.tmpl = tmpl
	return nil
}

// Routes registers the reader and studio endpoints.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}

	// SEO endpoints (no session needed).
	r.Get("/robots.txt", m.handleRobots)
	r.Get("/sitemap.xml", m.handleSitemap)

	// Public reader (session loaded softly so the header can show Studio link).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.LoadSession)
		r.Get("/", m.handleHome)
		r.Get("/read", m.handleReadRedirect)
		r.Get("/read/{slug}", m.handleArticle)
		r.Post("/read/{slug}/vote", m.handleVote)
		r.Post("/read/{slug}/comment", m.handleComment)
		r.Get("/author/sana", m.handleAuthorSana)
		r.Get("/about", m.handleStaticPage("about"))
		r.Get("/guide", m.handleStaticPage("guide"))
		r.Get("/formatting", m.handleStaticPage("formatting"))
		r.Get("/pricing", m.handleStaticPage("pricing"))
		r.Get("/support", m.handleStaticPage("support"))
		r.Get("/api/geo/roots", m.handleGeoRoots)
		r.Get("/api/geo/children", m.handleGeoChildren)
		r.Get("/api/listings/map", m.handleListingsMap)
		r.Get("/listings", m.handleListings)
		r.Get("/listings/new", m.handleListingNew)
		r.Post("/listings/new", m.handleListingCreate)
		r.Get("/listings/my", m.handleMyListings)
		r.Post("/listings/{id}/extend", m.handleListingExtend)
		r.Post("/listings/{id}/promote", m.handleListingPromote)
		r.Post("/listings/{id}/feature", m.handleListingFeature)
		r.Get("/listings/{id}", m.handleListingView)
	})

	// Studio auth pages (public).
	r.Get("/studio/login", m.handleLoginPage)
	r.Post("/studio/login", m.handleLoginSubmit)
	r.Get("/studio/register", m.handleRegisterPage)
	r.Post("/studio/register", m.handleRegisterSubmit)
	r.Get("/studio/logout", m.handleLogout)

	// Studio (authenticated author cabinet).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.RequireSession("/studio/login", "user", "operator", "admin"))
		r.Get("/studio", m.handleDashboard)
		r.Get("/studio/new", m.handleEditorNew)
		r.Post("/studio/new", m.handleCreate)
		r.Get("/studio/a/{id}", m.handleEditorEdit)
		r.Post("/studio/a/{id}", m.handleUpdate)
		r.Post("/studio/a/{id}/publish", m.handlePublish)
		r.Post("/studio/a/{id}/unpublish", m.handleUnpublish)
		r.Post("/studio/a/{id}/improve", m.handleImprove)
		r.Post("/studio/a/{id}/translate", m.handleTranslate)
		r.Get("/favorites", m.handleFavorites)
		r.Post("/favorites/{type}/{id}", m.handleFavoriteToggle)
	})

	// Admin control panel (staff roles only).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.RequireSession("/studio/login", adminRoles...))
		r.Get("/admin", m.handleAdmin)
		r.Post("/admin/roles", m.handleAdminAssignRole)
		r.Post("/admin/comments/{id}/hide", m.handleAdminHideComment)
	})
}

func (m *Module) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := m.tmpl.ExecuteTemplate(w, name, data); err != nil {
		m.rt.Logger.Error("render article template", zap.String("template", name), zap.Error(err))
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
	shanraq.StarterModule
} = (*Module)(nil)
