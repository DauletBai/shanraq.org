package articles

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/modules/ratings"
)

const langCookieName = "shanraq_lang"

// Base carries fields shared by every page (consumed by the header/footer
// partials via Go's embedded-field promotion). The whole UI renders in Lang.
type Base struct {
	Title     string
	Lang      string
	Authed    bool
	IsStaff   bool
	ShowLangs bool
	Active    string // active section: "latest" | "top" | ""
	ActiveCat string // active category slug, or "" for All
	ActiveSub string // active subcategory slug, or ""
	LangLinks map[string]string

	// SidebarNews feeds the "latest news" carousel in the sidebar.
	SidebarNews []FeedItem

	// SEO fields (populated by base(); pages may override).
	SiteURL string        // absolute origin, e.g. https://shanraq.org
	Path    string        // request path, no query
	Desc    string        // meta description
	OGImage string        // absolute image URL for social previews
	OGType  string        // "website" | "article"
	JSONLD  template.HTML // structured data (schema.org), injected verbatim
}

// base builds the shared page context. The language switcher points at the
// current path so switching language re-renders the same page fully localized.
func (m *Module) base(r *http.Request, title, lang string) Base {
	claims, authed := auth.ClaimsFromContext(r.Context())
	site := m.rt.Config.PublicBase()
	return Base{
		Title:     title,
		Lang:      lang,
		Authed:    authed,
		IsStaff:   authed && claims.HasAnyRole(adminRoles...),
		ShowLangs: true,
		LangLinks: langLinks(r.URL.Path, lang),
		SiteURL:   site,
		Path:      r.URL.Path,
		Desc:      T(lang, "seo.site_desc"),
		OGImage:   site + "/static/brand/shanraq.svg",
		OGType:    "website",
	}
}

// resolveLang picks the active language from ?lang=, then cookie, then default.
func (m *Module) resolveLang(w http.ResponseWriter, r *http.Request) string {
	if q := r.URL.Query().Get("lang"); IsLang(q) {
		http.SetCookie(w, &http.Cookie{Name: langCookieName, Value: q, Path: "/", MaxAge: 31536000, SameSite: http.SameSiteLaxMode})
		return q
	}
	if c, err := r.Cookie(langCookieName); err == nil && IsLang(c.Value) {
		return c.Value
	}
	return LangRU
}

// latestNews returns the newest published articles for the sidebar carousel.
func (m *Module) latestNews(r *http.Request, lang string, n int) []FeedItem {
	arts, err := m.store.ListPublished(r.Context(), "", "", "", n, 0)
	if err != nil {
		m.rt.Logger.Warn("sidebar news", zap.Error(err))
		return nil
	}
	out := make([]FeedItem, 0, len(arts))
	for _, a := range arts {
		tr, served := a.Translation(lang)
		if tr == nil {
			continue
		}
		name, aiAuthor := authorDisplay(a)
		out = append(out, FeedItem{
			Slug: a.Slug, Title: tr.Title, AuthorName: name, AIAuthor: aiAuthor,
			ServedLang: served, Category: a.Category, CoverURL: a.CoverURL,
		})
	}
	return out
}

func langLinks(base string, current string) map[string]string {
	out := make(map[string]string, len(Langs))
	for _, l := range Langs {
		out[l] = base + "?lang=" + l
	}
	_ = current
	return out
}

func (m *Module) authorID(r *http.Request) (uuid.UUID, bool) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// viewerID returns the current user's ID, or uuid.Nil for anonymous readers.
func (m *Module) viewerID(r *http.Request) uuid.UUID {
	if id, ok := m.authorID(r); ok {
		return id
	}
	return uuid.Nil
}

// ---------- public reader ----------

// FeedItem is one card in the feed.
type FeedItem struct {
	Slug           string
	Title          string
	Summary        string
	AuthorName     string
	ServedLang     string
	Category       string
	Subcategory    string
	CoverURL       string
	Published      *time.Time
	Views          int64
	Score          int
	IsAI           bool
	AIAuthor       bool
	AvailableLangs []string
}

// HomePage is the template context for the portal home.
type HomePage struct {
	Base
	Featured   *FeedItem
	Posts      []FeedItem
	Recent     []FeedItem
	Subscribed bool
}

// recentSlice returns up to n items for the sidebar "recent" list.
func recentSlice(items []FeedItem, n int) []FeedItem {
	if len(items) > n {
		return items[:n]
	}
	return items
}

// StaticPage backs the About / Guide / Support info pages.
type StaticPage struct {
	Base
	Body interface{}
}

// handleStaticPage renders a localized info page by key.
func (m *Module) handleStaticPage(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := m.resolveLang(w, r)
		c := staticContent(key, lang)
		if c.Title == "" {
			http.NotFound(w, r)
			return
		}
		page := StaticPage{Base: m.base(r, c.Title, lang)}
		page.Body = RenderMarkdown(c.Body)
		m.render(w, "page", page)
	}
}

// handleReadRedirect keeps the old /read URL working by sending it home.
func (m *Module) handleReadRedirect(w http.ResponseWriter, r *http.Request) {
	target := "/"
	if q := r.URL.RawQuery; q != "" {
		target += "?" + q
	}
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func (m *Module) handleHome(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)

	sort := "recent"
	active := "latest"
	if r.URL.Query().Get("sort") == "top" {
		sort = "top"
		active = "top"
	}
	cat := ""
	if c := r.URL.Query().Get("cat"); IsCategory(c) {
		cat = c
	}
	sub := ""
	if s := r.URL.Query().Get("sub"); IsSubcategory(s) {
		sub = s
		cat = SubcategoryParent(s) // a subcategory implies its parent category
	}

	arts, err := m.store.ListPublished(r.Context(), sort, cat, sub, 21, 0)
	if err != nil {
		m.rt.Logger.Error("home list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	items := make([]FeedItem, 0, len(arts))
	for _, a := range arts {
		tr, served := a.Translation(lang)
		if tr == nil {
			continue
		}
		summary := tr.Summary
		if summary == "" {
			summary = excerpt(stripMD(tr.BodyMD), 170)
		}
		authorName, aiAuthor := authorDisplay(a)
		items = append(items, FeedItem{
			Slug:           a.Slug,
			Title:          tr.Title,
			Summary:        summary,
			AuthorName:     authorName,
			AIAuthor:       aiAuthor,
			ServedLang:     served,
			Category:       a.Category,
			Subcategory:    a.Subcategory,
			CoverURL:       a.CoverURL,
			Published:      a.PublishedAt,
			Views:          a.ViewsCount,
			Score:          a.Score,
			IsAI:           tr.Source == "ai",
			AvailableLangs: a.AvailableLangs(),
		})
	}

	page := HomePage{Base: m.base(r, T(lang, "home.page_title"), lang)}
	page.Active = active
	page.ActiveCat = cat
	page.ActiveSub = sub
	page.Subscribed = r.URL.Query().Get("subscribed") == "ok"
	page.Posts = items
	page.Recent = recentSlice(items, 5)
	page.SidebarNews = m.latestNews(r, lang, 6)
	m.render(w, "home", page)
}

// ArticlePage is the template context for a single article.
type ArticlePage struct {
	Base
	ArticleID      string
	Slug           string
	Title          string
	Summary        string
	AuthorName     string
	ServedLang     string
	RequestedLang  string
	Body           interface{}
	Published      *time.Time
	Views          int64
	IsAI           bool
	AIAuthor       bool
	Translated     bool
	AvailableLangs []string

	Category      string
	Subcategory   string
	CoverURL      string
	Score         int
	UserVote      int // -1, 0, +1
	AuthorKarma   int
	CanVote       bool // logged in and not the author
	IsAuthor      bool
	Recent        []FeedItem // reserved for sidebar
	Subscribed    bool
	Comments      []Comment
	IsFavorite    bool
	TOC           []TOCItem
	ReadingMin    int
	CommentReview bool // the reader's comment was held for moderation
}

func (m *Module) handleArticle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	lang := m.resolveLang(w, r)

	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tr, served := a.Translation(lang)
	if tr == nil {
		http.NotFound(w, r)
		return
	}

	// Count the view asynchronously-ish; ignore errors (best effort analytics).
	if err := m.store.RecordView(r.Context(), a.ID, served); err != nil {
		m.rt.Logger.Warn("record view", zap.Error(err))
	}

	page := ArticlePage{Base: m.base(r, tr.Title, lang)}
	page.ArticleID = a.ID.String()
	page.Slug = a.Slug
	page.Category = a.Category
	page.Subcategory = a.Subcategory
	page.CoverURL = a.CoverURL
	page.Title = tr.Title
	page.Summary = tr.Summary
	page.AuthorName, page.AIAuthor = authorDisplay(a)
	page.ServedLang = served
	page.RequestedLang = lang
	page.Body, page.TOC = RenderMarkdownTOC(tr.BodyMD)
	page.ReadingMin = readingMinutes(tr.BodyMD)
	page.Published = a.PublishedAt
	page.Views = a.ViewsCount + 1
	page.IsAI = tr.Source == "ai"
	page.Translated = served != lang
	page.AvailableLangs = a.AvailableLangs()

	viewer := m.viewerID(r)
	if rating, err := m.ratings.ForArticle(r.Context(), a.ID, viewer); err == nil {
		page.Score = rating.Score
		page.UserVote = rating.UserVote
	} else {
		m.rt.Logger.Warn("article rating", zap.Error(err))
	}
	if karma, err := m.ratings.AuthorKarma(r.Context(), a.AuthorID); err == nil {
		page.AuthorKarma = karma
	}
	page.IsAuthor = viewer != uuid.Nil && viewer == a.AuthorID
	page.CanVote = viewer != uuid.Nil && !page.IsAuthor
	if viewer != uuid.Nil {
		page.IsFavorite = m.favs.IsFavorite(r.Context(), viewer, "article", a.ID)
	}

	page.CommentReview = r.URL.Query().Get("comment") == "review"
	// The article page shows only its table of contents in the aside (no news
	// carousel / widgets), so SidebarNews is intentionally not populated here.
	if cs, err := m.comments.ListForArticle(r.Context(), a.ID); err == nil {
		page.Comments = cs
	} else {
		m.rt.Logger.Warn("load comments", zap.Error(err))
	}

	m.applyArticleSEO(&page)
	m.render(w, "article", page)
}

// handleVote records a reader's up/down vote (toggling off when the same
// direction is submitted twice), then returns to the article.
func (m *Module) handleComment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	backTo := "/read/" + slug + "#comments"

	userID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Redirect(w, r, backTo, http.StatusSeeOther)
		return
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// AI pre-moderation (no-op when the assistant is off): a flagged comment is
	// filed as hidden into the same queue human reports feed. On any error we
	// fail open and publish — the human report button remains the backstop.
	status := "published"
	if m.ai != nil && m.ai.Enabled() {
		if v, err := m.ai.Moderate(r.Context(), "comment", body); err == nil && v.Flagged() {
			status = "hidden"
			m.rt.Logger.Info("ai moderator hid comment",
				zap.String("reason", v.Reason), zap.Float64("confidence", v.Confidence))
		}
	}
	if err := m.comments.CreateWithStatus(r.Context(), a.ID, userID, body, status); err != nil {
		m.rt.Logger.Error("create comment", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if status == "hidden" {
		backTo = "/read/" + slug + "?comment=review#comments"
	}
	http.Redirect(w, r, backTo, http.StatusSeeOther)
}

func (m *Module) handleVote(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	backTo := "/read/" + slug

	voter, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	value := ratings.VoteNone
	switch r.FormValue("value") {
	case "1", "up":
		value = ratings.VoteUp
	case "-1", "down":
		value = ratings.VoteDown
	default:
		http.Redirect(w, r, backTo, http.StatusSeeOther)
		return
	}

	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Clicking the current direction again retracts the vote.
	if cur, err := m.ratings.ForArticle(r.Context(), a.ID, voter); err == nil && cur.UserVote == value {
		value = ratings.VoteNone
	}

	if _, err := m.ratings.Vote(r.Context(), a.ID, voter, a.AuthorID, value); err != nil {
		if errors.Is(err, ratings.ErrSelfVote) {
			http.Redirect(w, r, backTo, http.StatusSeeOther)
			return
		}
		m.rt.Logger.Error("vote", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, backTo, http.StatusSeeOther)
}

// ---------- auth pages ----------

// FormPage backs the login and register screens.
type FormPage struct {
	Base
	Mode  string // login | register
	Email string
	Error string
}

func (m *Module) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	m.render(w, "form", FormPage{Base: m.base(r, T(lang, "form.login_title"), lang), Mode: "login"})
}

func (m *Module) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	m.render(w, "form", FormPage{Base: m.base(r, T(lang, "form.register_title"), lang), Mode: "register"})
}

func (m *Module) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if !m.auth.AllowAuthAttempt(r, "signin", email) {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.login_title"), lang),
			Mode:  "login",
			Email: email,
			Error: T(lang, "form.err_rate_limit"),
		})
		return
	}

	user, token, err := m.auth.LoginPassword(r.Context(), email, password)
	if err != nil {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.login_title"), lang),
			Mode:  "login",
			Email: email,
			Error: T(lang, "form.err_credentials"),
		})
		return
	}
	auth.SetSessionCookie(w, r, token, m.auth.SessionTTL())
	m.rt.Logger.Info("studio login", zap.String("user_id", user.ID.String()))
	http.Redirect(w, r, "/studio", http.StatusSeeOther)
}

func (m *Module) handleRegisterSubmit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if !m.auth.AllowAuthAttempt(r, "signup", email) {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.register_title"), lang),
			Mode:  "register",
			Email: email,
			Error: T(lang, "form.err_rate_limit"),
		})
		return
	}

	// KZ online-platform law: registration cannot complete without explicit
	// consent to the Terms and Privacy Policy.
	if r.FormValue("consent") != "on" {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.register_title"), lang),
			Mode:  "register",
			Email: email,
			Error: T(lang, "form.err_consent"),
		})
		return
	}

	if _, ok := auth.NormalizeEmail(email); !ok {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.register_title"), lang),
			Mode:  "register",
			Email: email,
			Error: T(lang, "form.err_email_invalid"),
		})
		return
	}
	if len(password) < 8 {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.register_title"), lang),
			Mode:  "register",
			Email: email,
			Error: T(lang, "form.err_short_pw"),
		})
		return
	}

	user, token, err := m.auth.RegisterPassword(r.Context(), email, password)
	if err != nil {
		msg := T(lang, "form.err_generic")
		if errors.Is(err, auth.ErrEmailExists) {
			msg = T(lang, "form.err_email_taken")
		} else if errors.Is(err, auth.ErrInvalidEmail) {
			msg = T(lang, "form.err_email_invalid")
		}
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.register_title"), lang),
			Mode:  "register",
			Email: email,
			Error: msg,
		})
		return
	}
	auth.SetSessionCookie(w, r, token, m.auth.SessionTTL())
	m.rt.Logger.Info("studio register", zap.String("user_id", user.ID.String()))
	http.Redirect(w, r, "/studio", http.StatusSeeOther)
}

func (m *Module) handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessionCookie(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------- studio ----------

// StudioRow is one row in the author's article table.
type StudioRow struct {
	ID      string
	Slug    string
	Title   string
	Status  string
	Updated time.Time
	Views   int64
	Langs   []string
	// Reading-depth funnel: reader counts and their share of views (percent).
	D25, D50, D75, D100 int64
	P25, P50, P75, P100 int
}

// StudioPage is the dashboard context.
type StudioPage struct {
	Base
	Stats    AuthorStats
	Karma    int
	Articles []StudioRow
}

func (m *Module) handleDashboard(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}

	stats, err := m.store.AuthorStats(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Error("author stats", zap.Error(err))
	}
	arts, err := m.store.ListByAuthor(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Error("author list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	depth, err := m.store.AuthorReadingDepth(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Warn("author reading depth", zap.Error(err))
	}

	lang := m.resolveLang(w, r)
	rows := make([]StudioRow, 0, len(arts))
	for _, a := range arts {
		title := T(lang, "studio.untitled")
		if tr, _ := a.Translation(a.OriginalLang); tr != nil && tr.Title != "" {
			title = tr.Title
		}
		row := StudioRow{
			ID:      a.ID.String(),
			Slug:    a.Slug,
			Title:   title,
			Status:  a.Status,
			Updated: a.UpdatedAt,
			Views:   a.ViewsCount,
			Langs:   a.AvailableLangs(),
		}
		if d := depth[row.ID]; d != nil {
			row.D25, row.D50, row.D75, row.D100 = d[25], d[50], d[75], d[100]
			row.P25 = pctOf(row.D25, row.Views)
			row.P50 = pctOf(row.D50, row.Views)
			row.P75 = pctOf(row.D75, row.Views)
			row.P100 = pctOf(row.D100, row.Views)
		}
		rows = append(rows, row)
	}

	karma, err := m.ratings.AuthorKarma(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Warn("author karma", zap.Error(err))
	}

	page := StudioPage{Base: m.base(r, T(lang, "studio.title"), lang)}
	page.Stats = stats
	page.Karma = karma
	page.Articles = rows
	m.render(w, "studio_dashboard", page)
}

// TranslationField holds editable fields for one language tab.
type TranslationField struct {
	Title   string
	Summary string
	BodyMD  string
	Source  string
}

// EditorPage backs the trilingual editor.
type EditorPage struct {
	Base
	IsNew        bool
	ArticleID    string
	Slug         string
	Status       string
	OriginalLang string
	Category     string
	Subcategory  string
	CoverURL     string
	Fields       map[string]TranslationField
	Error        string
	AIEnabled    bool
	Notice       string
}

// aiNotice maps an ?ai= redirect flag to a localized message.
func aiNotice(lang, flag string) string {
	switch flag {
	case "improved":
		return T(lang, "notice.ai_improved")
	case "queued":
		return T(lang, "notice.ai_queued")
	case "off":
		return T(lang, "notice.ai_off")
	case "drafted":
		return T(lang, "notice.ai_drafted")
	case "draft_skip":
		return T(lang, "notice.ai_draft_skip")
	case "draft_empty":
		return T(lang, "notice.ai_draft_empty")
	default:
		return ""
	}
}

func emptyFields() map[string]TranslationField {
	f := make(map[string]TranslationField, len(Langs))
	for _, l := range Langs {
		f[l] = TranslationField{}
	}
	return f
}

func (m *Module) handleEditorNew(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	page := EditorPage{Base: m.base(r, T(lang, "editor.new"), lang)}
	page.IsNew = true
	page.OriginalLang = lang
	page.Category = "society"
	page.Subcategory = ""
	page.Status = "draft"
	page.Fields = emptyFields()
	page.AIEnabled = m.ai.Enabled()
	m.render(w, "studio_editor", page)
}

func (m *Module) handleEditorEdit(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a, err := m.store.GetByID(r.Context(), id, authorID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fields := emptyFields()
	for _, l := range Langs {
		if tr, ok := a.Translations[l]; ok {
			fields[l] = TranslationField{Title: tr.Title, Summary: tr.Summary, BodyMD: tr.BodyMD, Source: tr.Source}
		}
	}

	lang := m.resolveLang(w, r)
	page := EditorPage{Base: m.base(r, T(lang, "editor.edit"), lang)}
	page.ArticleID = a.ID.String()
	page.Slug = a.Slug
	page.Status = a.Status
	page.OriginalLang = a.OriginalLang
	page.Category = a.Category
	page.Subcategory = a.Subcategory
	page.CoverURL = a.CoverURL
	page.Fields = fields
	page.AIEnabled = m.ai.Enabled()
	page.Notice = aiNotice(lang, r.URL.Query().Get("ai"))
	m.render(w, "studio_editor", page)
}

// handleImprove rewrites the original-language body with the AI co-editor and
// saves the result back to the draft.
func (m *Module) handleImprove(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	editURL := "/studio/a/" + id.String()

	if !m.ai.Enabled() {
		http.Redirect(w, r, editURL+"?ai=off", http.StatusSeeOther)
		return
	}

	a, err := m.store.GetByID(r.Context(), id, authorID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tr, ok := a.Translations[a.OriginalLang]
	if !ok || tr.BodyMD == "" {
		http.Redirect(w, r, editURL, http.StatusSeeOther)
		return
	}

	improved, err := m.ai.Improve(r.Context(), a.OriginalLang, tr.BodyMD)
	if err != nil {
		m.rt.Logger.Error("ai improve", zap.Error(err))
		http.Redirect(w, r, editURL+"?ai=off", http.StatusSeeOther)
		return
	}

	input := []TranslationInput{{
		Lang:    a.OriginalLang,
		Title:   tr.Title,
		Summary: tr.Summary,
		BodyMD:  improved,
		Source:  "human",
	}}
	if err := m.store.Update(r.Context(), id, authorID, a.OriginalLang, a.Category, a.Subcategory, a.CoverURL, input); err != nil {
		m.rt.Logger.Error("save improved", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, editURL+"?ai=improved", http.StatusSeeOther)
}

// handleDraft asks the Sana Qyran columnist to draft an evergreen column body
// from a short brief and fills it into the empty original-language body. It
// refuses to overwrite an existing body so the author never loses their work.
func (m *Module) handleDraft(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	editURL := "/studio/a/" + id.String()

	if !m.ai.Enabled() {
		http.Redirect(w, r, editURL+"?ai=off", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	brief := strings.TrimSpace(r.FormValue("brief"))

	a, err := m.store.GetByID(r.Context(), id, authorID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tr := a.Translations[a.OriginalLang]
	if brief == "" {
		brief = strings.TrimSpace(tr.Title)
	}
	if brief == "" {
		http.Redirect(w, r, editURL+"?ai=draft_empty", http.StatusSeeOther)
		return
	}
	// Never clobber existing writing — draft only into an empty body.
	if strings.TrimSpace(tr.BodyMD) != "" {
		http.Redirect(w, r, editURL+"?ai=draft_skip", http.StatusSeeOther)
		return
	}

	body, err := m.ai.DraftColumn(r.Context(), a.OriginalLang, brief)
	if err != nil || strings.TrimSpace(body) == "" {
		m.rt.Logger.Error("ai draft", zap.Error(err))
		http.Redirect(w, r, editURL+"?ai=off", http.StatusSeeOther)
		return
	}

	input := []TranslationInput{{
		Lang:    a.OriginalLang,
		Title:   tr.Title,
		Summary: tr.Summary,
		BodyMD:  body,
		Source:  "ai",
	}}
	if err := m.store.Update(r.Context(), id, authorID, a.OriginalLang, a.Category, a.Subcategory, a.CoverURL, input); err != nil {
		m.rt.Logger.Error("save draft", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, editURL+"?ai=drafted", http.StatusSeeOther)
}

// handleTranslate enqueues an async AI translation into the other languages.
func (m *Module) handleTranslate(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	editURL := "/studio/a/" + id.String()

	if !m.ai.Enabled() {
		http.Redirect(w, r, editURL+"?ai=off", http.StatusSeeOther)
		return
	}

	// Ensure the article belongs to this author before enqueuing work for it.
	if _, err := m.store.GetByID(r.Context(), id, authorID); err != nil {
		http.NotFound(w, r)
		return
	}

	payload, err := ai.EnqueuePayload(id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	job := jobs.Job{
		ID:          uuid.New(),
		UserID:      authorID,
		Name:        ai.JobTranslate,
		Payload:     payload,
		RunAt:       time.Now(),
		MaxAttempts: 3,
	}
	if err := m.jobs.Enqueue(r.Context(), job); err != nil {
		m.rt.Logger.Error("enqueue translate", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, editURL+"?ai=queued", http.StatusSeeOther)
}

// parseEditorForm reads the three-language editor form into translation inputs.
func parseEditorForm(r *http.Request) (originalLang, category, subcategory, coverURL string, trs []TranslationInput) {
	originalLang = r.FormValue("original_lang")
	if !IsLang(originalLang) {
		originalLang = LangRU
	}
	category = NormalizeCategory(r.FormValue("category"))
	subcategory = NormalizeSubcategory(category, r.FormValue("subcategory"))
	coverURL = strings.TrimSpace(r.FormValue("cover_url"))
	for _, l := range Langs {
		trs = append(trs, TranslationInput{
			Lang:    l,
			Title:   strings.TrimSpace(r.FormValue("title_" + l)),
			Summary: strings.TrimSpace(r.FormValue("summary_" + l)),
			BodyMD:  r.FormValue("body_" + l),
			Source:  "human",
		})
	}
	return originalLang, category, subcategory, coverURL, trs
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	originalLang, category, subcategory, coverURL, trs := parseEditorForm(r)

	lang := m.resolveLang(w, r)
	orig := findTR(trs, originalLang)
	if orig.Title == "" || orig.BodyMD == "" {
		m.reRenderEditor(w, r, true, "", "", originalLang, category, subcategory, coverURL, trs, T(lang, "editor.err_required"))
		return
	}

	slug, err := m.uniqueSlug(r, orig.Title)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	id, err := m.store.Create(r.Context(), authorID, slug, originalLang, category, subcategory, coverURL, trs)
	if err != nil {
		m.rt.Logger.Error("create article", zap.Error(err))
		m.reRenderEditor(w, r, true, "", "", originalLang, category, subcategory, coverURL, trs, T(lang, "editor.err_save"))
		return
	}
	http.Redirect(w, r, "/studio/a/"+id.String(), http.StatusSeeOther)
}

func (m *Module) handleUpdate(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	originalLang, category, subcategory, coverURL, trs := parseEditorForm(r)

	if err := m.store.Update(r.Context(), id, authorID, originalLang, category, subcategory, coverURL, trs); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("update article", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/studio/a/"+id.String(), http.StatusSeeOther)
}

func (m *Module) handlePublish(w http.ResponseWriter, r *http.Request) {
	m.transition(w, r, "published", "/studio")
}

func (m *Module) handleUnpublish(w http.ResponseWriter, r *http.Request) {
	m.transition(w, r, "draft", "/studio")
}

func (m *Module) transition(w http.ResponseWriter, r *http.Request, status, redirect string) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := m.store.SetStatus(r.Context(), id, authorID, status); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("transition", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// On publish, push the article out to external channels (Telegram) for
	// block-resilience. Best-effort: a syndication failure must not block the
	// publish itself.
	if status == "published" {
		if err := m.syndicate.EnqueuePublish(r.Context(), m.jobs, id); err != nil {
			m.rt.Logger.Warn("enqueue syndication", zap.Error(err))
		}
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (m *Module) reRenderEditor(w http.ResponseWriter, r *http.Request, isNew bool, id, slug, originalLang, category, subcategory, coverURL string, trs []TranslationInput, errMsg string) {
	fields := emptyFields()
	for _, tr := range trs {
		fields[tr.Lang] = TranslationField{Title: tr.Title, Summary: tr.Summary, BodyMD: tr.BodyMD, Source: tr.Source}
	}
	lang := m.resolveLang(w, r)
	page := EditorPage{Base: m.base(r, T(lang, "editor.new"), lang)}
	page.IsNew = isNew
	page.ArticleID = id
	page.Slug = slug
	page.OriginalLang = originalLang
	page.Category = category
	page.Subcategory = subcategory
	page.CoverURL = coverURL
	page.Status = "draft"
	page.Fields = fields
	page.AIEnabled = m.ai.Enabled()
	page.Error = errMsg
	m.render(w, "studio_editor", page)
}

func (m *Module) uniqueSlug(r *http.Request, title string) (string, error) {
	base := Slugify(title)
	slug := base
	exists, err := m.store.SlugExists(r.Context(), slug)
	if err != nil {
		return "", err
	}
	if exists {
		slug = base + "-" + uuid.NewString()[:6]
	}
	return slug, nil
}

func findTR(trs []TranslationInput, lang string) TranslationInput {
	for _, tr := range trs {
		if tr.Lang == lang {
			return tr
		}
	}
	return TranslationInput{Lang: lang}
}
