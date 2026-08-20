package articles

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/modules/media"
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
	rt            *shanraq.Runtime
	store         *Store
	listings      *ListingStore
	ads           *AdStore
	mods          *ModStore
	refs          *ReferralStore
	reagents      *AgentStore
	pay           *PaymentStore
	paySettings   *PaymentSettingsStore
	flags         *ServiceFlags
	geo           *GeoStore
	comments      *CommentStore
	favs          *FavoriteStore
	admin         *AdminStore
	users         *auth.Store // account administration (list / edit / delete)
	content       *ContentStore
	predictions   *PredictionStore
	tariffs       *TariffStore
	metrics       *Metrics
	geoip         *geoIP
	excludeEmails map[string]bool
	ratings       *ratings.Store
	jobs          *jobs.Store
	auth          *auth.Module
	ai            *ai.Module
	syndicate     *syndicate.Module
	media         *media.Module
	mailer        Mailer
	tmpl          *template.Template
	validator     *validate.Validator
	infobar       *InfoBar
}

// New builds the articles module. It depends on auth (browser sessions), ai
// (writing assistant), syndicate (publish → external channels), and a mailer
// (listing expiry reminders).
func New(authModule *auth.Module, aiModule *ai.Module, syndicateModule *syndicate.Module, mediaModule *media.Module, mailer Mailer) *Module {
	return &Module{auth: authModule, ai: aiModule, syndicate: syndicateModule, media: mediaModule, mailer: mailer}
}

func (m *Module) Name() string { return "articles" }

// Init parses templates and wires the store.
func (m *Module) Init(ctx context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	m.listings = NewListingStore(rt.DB)
	m.ads = NewAdStore(rt.DB)
	m.mods = NewModStore(rt.DB)
	m.refs = NewReferralStore(rt.DB)
	m.reagents = NewAgentStore(rt.DB)
	m.pay = NewPaymentStore(rt.DB)
	// Which acquirer is live (and whether payments are on) is a runtime choice in
	// the admin panel — secrets stay in config. Defaults come from config so a
	// fresh DB has a starting point; the effective provider is built per request
	// from these settings (m.paymentProvider), disabled until an adapter lands.
	m.paySettings = NewPaymentSettingsStore(rt.DB, PaymentSettings{Provider: rt.Config.Payments.Provider})
	if _, err := m.paySettings.Load(ctx); err != nil {
		rt.Logger.Warn("load payment settings", zap.Error(err))
	}
	m.flags = NewServiceFlags(rt.DB)
	if err := m.flags.Load(ctx); err != nil {
		// Non-fatal: an unloaded cache defaults every service to available, so
		// the site still works; log and continue.
		rt.Logger.Warn("load service flags", zap.Error(err))
	}
	m.geo = NewGeoStore(rt.DB)
	m.comments = NewCommentStore(rt.DB)
	m.favs = NewFavoriteStore(rt.DB)
	m.admin = NewAdminStore(rt.DB)
	m.users = auth.NewStore(rt.DB)
	m.content = NewContentStore(rt.DB)
	m.predictions = NewPredictionStore(rt.DB)
	// Fill the editable-pages table from the built-in defaults on first boot;
	// idempotent and best-effort, so it never blocks startup.
	m.seedContentPages(ctx)
	m.tariffs = NewTariffStore(rt.DB)
	// Seed the rate card from the built-in defaults and load it into the cache;
	// best-effort, so a failure just serves the built-in prices.
	if err := m.tariffs.Load(ctx); err != nil {
		rt.Logger.Warn("load tariffs", zap.Error(err))
	}
	m.metrics = NewMetrics(rt.DB, rt.Logger)
	// Optional GeoIP: buckets visits by country so the audience panel can tell
	// domestic (KZ) readers from genuine foreign ones. The optional ASN database
	// separates hosting/cloud/VPN traffic ("datacenter") from real readers. Off
	// when no country DB is set.
	m.geoip = openGeoIP(rt.Config.Analytics.GeoIPDB, rt.Config.Analytics.GeoIPASNDB)
	if m.geoip != nil {
		rt.Logger.Info("country analytics enabled",
			zap.String("geoip_db", rt.Config.Analytics.GeoIPDB),
			zap.Bool("datacenter_filter", m.geoip.hasASN()))
	}
	// Emails whose traffic is excluded from analytics (the owner's test account,
	// etc.) — staff roles are excluded automatically. Comma-separated, from env.
	m.excludeEmails = map[string]bool{}
	for _, e := range strings.Split(rt.Config.Analytics.ExcludeEmails, ",") {
		if e = strings.ToLower(strings.TrimSpace(e)); e != "" {
			m.excludeEmails[e] = true
		}
	}
	m.ratings = ratings.NewStore(rt.DB)
	m.jobs = jobs.NewStore(rt.DB)
	m.validator = validate.New()
	m.infobar = NewInfoBar(rt.Logger, socialLinks(rt.Config.Social), rt.Config.Social.GitHub)

	tmpl, err := template.New("articles").Funcs(templateFuncs()).ParseFS(templateFiles, "templates/*.html")
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

	// Everything below lives in one inline group so the CSRF guard can be a
	// group middleware — the shared mux already has routes from other modules,
	// and chi forbids Use() after routes on the same mux.
	// Payment provider callbacks are server-to-server; they carry no Origin
	// header, so they must not go through the CSRF-guarded browser group.
	r.Post("/pay/webhook/{provider}", m.handlePaymentWebhook)
	r.Group(m.browserRoutes)
}

// browserRoutes registers all reader/studio/admin endpoints under a single
// group whose first middleware is the same-origin CSRF guard (the session
// cookie is already SameSite=Lax; this is defense in depth).
func (m *Module) browserRoutes(r chi.Router) {
	r.Use(m.verifyOrigin)
	// Soft-load the session so the maintenance guard can recognize staff and let
	// them through a global takedown; it never blocks a request on its own.
	r.Use(m.auth.LoadSession)
	// Global maintenance switch: when the site is down, everything below serves
	// a 503 maintenance page except staff and the admin/login recovery routes.
	r.Use(m.maintenanceGuard)

	// SEO endpoints (no session needed).
	r.Get("/robots.txt", m.handleRobots)
	r.Get("/llms.txt", m.handleLLMS)
	r.Get("/sitemap.xml", m.handleSitemap)
	r.Get("/sitemap-listings.xml", m.handleSitemapListings)
	r.Get("/sitemap-news.xml", m.handleSitemapNews)

	// Printed-campaign short links (/q/rd …). No session and no tracking of
	// their own: the hop is invisible to pageKind, so the arrival is counted
	// once, on the page it lands on, under the campaign's own label.
	r.Get("/q/{code}", m.handlePosterLink)

	// Public reader (session loaded softly so the header can show Studio link).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.LoadSession)
		// Aggregate audience counters (guest vs registered). Placed after the
		// soft session load so it can tell the two apart; records nothing else.
		r.Use(m.trackTraffic)
		r.Post("/api/track", m.handleTrack)
		r.Get("/", m.handleHome)
		// Uptime monitors and HTTP clients probe with HEAD, and chi answers 405
		// unless the method is registered. Go's ResponseWriter discards the body
		// for a HEAD automatically, so reusing the GET handler returns exactly
		// the right headers and nothing else.
		r.Head("/", m.handleHome)
		r.Get("/read", m.handleReadRedirect)
		r.Get("/read/{slug}", m.handleArticle)
		// The reader's half of moderation. A POST because it changes something,
		// and same-origin-checked with the rest of the browser surface.
		r.Post("/read/{slug}/report", m.handleArticleReport)
		r.Post("/read/{slug}/vote", m.handleVote)
		r.Post("/read/{slug}/comment", m.handleComment)
		r.Post("/read/{slug}/comment/{id}/delete", m.handleCommentDelete)
		r.Post("/read/{slug}/progress", m.handleReadProgress)
		r.Get("/author/{id}", m.handleAuthor)
		r.Get("/predictions", m.handlePredictions)
		r.Get("/about", m.handleStaticPage("about"))
		r.Get("/guide", m.handleStaticPage("guide"))
		r.Get("/formatting", m.handleStaticPage("formatting"))
		r.Get("/pricing", m.handleStaticPage("pricing"))
		r.Get("/support", m.handleStaticPage("support"))
		r.Get("/privacy", m.handleStaticPage("privacy"))
		r.Get("/terms", m.handleStaticPage("terms"))
		r.Get("/api/geo/roots", m.handleGeoRoots)
		r.Get("/api/geo/children", m.handleGeoChildren)
		r.Get("/api/geo/path", m.handleGeoPath)
		r.Get("/api/geocode", m.handleGeocode)
		r.Get("/api/geocode/reverse", m.handleGeocodeReverse)
		r.Get("/api/listings/map", m.handleListingsMap)
		r.Get("/api/listings/pins", m.handleListingPins)
		r.Get("/listings", m.handleListings)
		r.Get("/listings/new", m.handleListingNew)
		r.Post("/listings/new", m.handleListingCreate)
		r.Get("/listings/my", m.handleMyListings)
		r.Get("/listings/{id}/edit", m.handleListingEdit)
		r.Post("/listings/{id}/edit", m.handleListingUpdate)
		r.Post("/listings/{id}/delete", m.handleListingDelete)
		r.Post("/listings/{id}/extend", m.handleListingExtend)
		r.Post("/listings/{id}/promote", m.handleListingPromote)
		r.Post("/listings/{id}/promote-free", m.handleListingPromoteFree)
		r.Post("/listings/{id}/feature", m.handleListingFeature)
		r.Post("/listings/{id}/banner", m.handleListingBanner)
		r.Post("/listings/{id}/contact", m.handleListingContact)
		r.Get("/listings/{id}", m.handleListingView)
		r.Get("/agent/{id}", m.handleAgentPublic)
	})

	// Studio auth pages (public).
	r.Get("/studio/login", m.handleLoginPage)
	r.Post("/studio/login", m.handleLoginSubmit)
	r.Get("/studio/register", m.handleRegisterPage)
	r.Post("/studio/register", m.handleRegisterSubmit)
	r.Post("/studio/logout", m.handleLogout)

	// Studio (authenticated author cabinet).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.RequireSession("/studio/login", "user", "operator", "admin"))
		r.Get("/studio", m.handleDashboard)
		r.Get("/studio/profile", m.handleProfile)
		r.Post("/studio/bio", m.handleBioSave)
		r.Post("/studio/avatar", m.handleAvatarUpload)
		r.Post("/studio/avatar/delete", m.handleAvatarDelete)
		r.Post("/studio/delete", m.handleDeleteAccount)
		r.Get("/studio/author", m.handleAuthorVerifyPage)
		r.Post("/studio/author/name", m.handleAuthorName)
		r.Post("/studio/author/phone", m.handleAuthorPhone)
		r.Post("/studio/author/confirm", m.handleAuthorConfirm)
		r.Get("/studio/new", m.handleEditorNew)
		r.Get("/studio/consent", m.handleConsent)
		r.Get("/studio/invite", m.handleInvite)
		r.Get("/studio/moderation", m.handleMyModeration)
		r.Post("/studio/moderation/{id}/appeal", m.handleFileAppeal)
		r.Post("/studio/consent", m.handleConsentSubmit)
		r.Post("/studio/new", m.handleCreate)
		r.Get("/studio/a/{id}", m.handleEditorEdit)
		r.Post("/studio/a/{id}", m.handleUpdate)
		r.Post("/studio/a/{id}/publish", m.handlePublish)
		r.Post("/studio/a/{id}/unpublish", m.handleUnpublish)
		r.Post("/studio/a/{id}/delete", m.handleDeleteDraft)
		r.Post("/studio/a/{id}/translate", m.handleTranslate)
		r.Get("/favorites", m.handleFavorites)
		// Advertiser cabinet (Phase 0b MVP — order capture, billing later).
		r.Get("/agent", m.handleAgentCabinet)
		r.Post("/agent", m.handleAgentSave)
		r.Get("/advertise", m.handleAdvertise)
		r.Post("/advertise/company", m.handleAdvertiseCompany)
		r.Post("/advertise/order", m.handleAdvertiseOrder)
		r.Get("/advertise/availability", m.handleAdsAvailability)
		r.Post("/favorites/{type}/{id}", m.handleFavoriteToggle)
		r.Post("/listings/{id}/report", m.handleListingReport)
	})

	// Admin control panel (staff roles only).
	r.Group(func(r chi.Router) {
		r.Use(m.auth.RequireSession("/studio/login", adminRoles...))
		r.Get("/admin", m.handleAdmin)
		r.Post("/admin/roles", m.handleAdminAssignRole)
		r.Post("/admin/users/{id}", m.handleAdminUserUpdate)
		r.Post("/admin/users/{id}/role", m.handleAdminUserRole)
		r.Post("/admin/users/{id}/delete", m.handleAdminUserDelete)
		r.Post("/admin/services", m.handleAdminServiceFlag)
		r.Post("/admin/ai", m.handleAdminAI)
		r.Post("/admin/agents/{id}/decide", m.handleAdminAgentDecide)
		r.Post("/admin/comments/{id}/hide", m.handleAdminHideComment)
		r.Post("/admin/appeals/{id}/resolve", m.handleAdminResolveAppeal)
		r.Post("/admin/articles/{id}/decide", m.handleAdminDecideArticle)
		r.Get("/admin/predictions", m.handleAdminPredictions)
		r.Get("/admin/predictions/{id}", m.handleAdminPredictions)
		r.Post("/admin/predictions", m.handleAdminPredictionSave)
		r.Post("/admin/predictions/{id}/delete", m.handleAdminPredictionDelete)
		r.Get("/admin/pages", m.handleAdminPages)
		r.Get("/admin/pages/{key}", m.handleAdminPageEdit)
		r.Post("/admin/pages/{key}", m.handleAdminPageSave)
		r.Post("/admin/payments", m.handleAdminPayments)
		r.Get("/admin/tariffs", m.handleAdminTariffs)
		r.Post("/admin/tariffs", m.handleAdminTariffsSave)
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
