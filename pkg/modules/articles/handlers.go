package articles

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
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
	CanAuthor bool   // leadership who may publish without email/phone verification
	Avatar    string // current user's avatar URL ("" = none), for the header/cabinet
	ShowLangs bool
	Active    string // active section: "latest" | "top" | ""
	ActiveCat string // active category slug, or "" for All
	ActiveSub string // active subcategory slug, or ""
	LangLinks map[string]string

	// SidebarNews feeds the "latest news" carousel in the sidebar.
	SidebarNews []FeedItem

	// Ads feeds the sidebar ad carousel (demo placements for now).
	Ads []Ad

	// Info feeds the top info bar (date, weather, rates, social links).
	Info InfoBarData

	// NeedsMap loads Leaflet. Only the two pages that draw one set it: the
	// library and its stylesheet are ~270 KB and were being fetched on every
	// page of the site, including the home feed, which has no map.
	NeedsMap bool

	// Newsletter form feedback, set from ?subscribed= after the POST redirect.
	// It lives on Base rather than one page's context because the form sits in
	// the follow card, which the home sidebar and every article aside share.
	SubMsg string
	SubBad bool // the message is a failure, not a confirmation

	// SEO fields (populated by base(); pages may override).
	SiteURL  string // absolute origin, e.g. https://shanraq.org
	Path     string // request path, no query (used for nav active state)
	CanonURL string // relative canonical path+query for THIS language,
	//                        including only whitelisted indexable filters
	Desc    string        // meta description
	OGImage string        // absolute image URL for social previews
	OGType  string        // "website" | "article"
	JSONLD  template.HTML // structured data (schema.org), injected verbatim
	// NoIndex asks search engines to keep this page out of their index while
	// still following its links. Set for articles flagged non-indexable.
	NoIndex bool

	// Svc carries the operational state of each toggleable service, already
	// localized, so any template can show a maintenance notice and hide a paid
	// action without a funcmap. Keyed by service code (e.g. "listing_promo").
	Svc map[string]ServiceView
}

// ServiceView is a service's state as a template sees it: whether its paid
// action is available, and the localized notice to show when it is not.
type ServiceView struct {
	On  bool
	Msg string
}

// serviceLinkOff reports whether a service's entry point should be disabled in
// the UI. An unknown/unconfigured code is treated as available, so a missing
// flag never hides a link. Exposed to templates as "svcOff": it greys out
// links/buttons that lead to a service the admin turned off or set to maintenance.
func serviceLinkOff(svc map[string]ServiceView, code string) bool {
	if v, ok := svc[code]; ok {
		return !v.On
	}
	return false
}

// serviceLinkMsg returns the localized "temporarily unavailable / by invitation"
// notice for a service, for use as a tooltip on the disabled entry point.
// Exposed to templates as "svcMsg".
func serviceLinkMsg(svc map[string]ServiceView, code string) string {
	if v, ok := svc[code]; ok {
		return v.Msg
	}
	return ""
}

// base builds the shared page context. The language switcher points at the
// current path so switching language re-renders the same page fully localized.
func (m *Module) base(r *http.Request, title, lang string) Base {
	claims, authed := auth.ClaimsFromContext(r.Context())
	site := m.rt.Config.PublicBase()
	avatar := ""
	if authed && claims != nil {
		if id, err := uuid.Parse(claims.Subject); err == nil {
			avatar = m.auth.Avatar(r.Context(), id)
		}
	}
	subMsg, subBad := subscribeFeedback(r, lang)
	return Base{
		Title:     title,
		Lang:      lang,
		SubMsg:    subMsg,
		SubBad:    subBad,
		Authed:    authed,
		IsStaff:   authed && claims.HasAnyRole(adminRoles...),
		CanAuthor: canAuthorAsStaff(claims),
		Avatar:    avatar,
		ShowLangs: true,
		LangLinks: langLinks(r.URL.Path, seoFilterQuery(r)),
		SiteURL:   site,
		Path:      r.URL.Path,
		CanonURL:  canonURL(r.URL.Path, seoFilterQuery(r), lang),
		Desc:      T(lang, "seo.site_desc"),
		OGImage:   site + "/static/brand/og-cover.png",
		OGType:    "website",
		Info:      m.infobar.Snapshot(localizedDate(lang, siteNow()), siteNow().Format("2006-01-02")),
		Ads:       m.sidebarAds(r, lang),
		Svc:       m.serviceViews(r, lang),
	}
}

// subscribeFeedback turns the ?subscribed= marker left by the syndicate module's
// redirect into a line under the newsletter form. "pending" is the success path:
// the address is stored but silent until the reader opens the confirmation link,
// so the copy has to send them to their inbox rather than say "subscribed".
func subscribeFeedback(r *http.Request, lang string) (string, bool) {
	switch r.URL.Query().Get("subscribed") {
	case "pending":
		return T(lang, "sidebar.subscribe_pending"), false
	case "bad":
		return T(lang, "sidebar.subscribe_bad"), true
	case "err":
		return T(lang, "sidebar.subscribe_err"), true
	}
	return "", false
}

// serviceViews snapshots every known service's state for THIS viewer, localized.
// For invite_only, On is true when the viewer is invited, so an invited tester
// sees no banner and keeps the action while the public sees "by invitation".
func (m *Module) serviceViews(r *http.Request, lang string) map[string]ServiceView {
	if m.flags == nil {
		return nil
	}
	invited, invitedComputed := false, false
	out := make(map[string]ServiceView, len(knownServices))
	for _, f := range m.flags.All() {
		on := f.Available()
		msg := f.Message(lang)
		switch f.Status {
		case svcInviteOnly:
			if !invitedComputed {
				invited, invitedComputed = m.isInvited(r), true
			}
			if invited {
				on = true
			} else if msg == "" {
				msg = T(lang, "svc.invite_note")
			}
		default:
			if !on && msg == "" {
				msg = T(lang, "svc.closed_note")
			}
		}
		out[f.Code] = ServiceView{On: on, Msg: msg}
	}
	return out
}

// isInvited reports whether the current viewer may use invite_only functions:
// staff always, plus any user who joined through an invite link.
func (m *Module) isInvited(r *http.Request) bool {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		return false
	}
	if claims.HasAnyRole(adminRoles...) {
		return true
	}
	if id, ok := m.authorID(r); ok {
		return m.refs.IsReferred(r.Context(), id)
	}
	return false
}

// gateReason decides whether the viewer may perform a flagged free action and,
// when not, returns the localized notice explaining why (by-invitation vs
// temporarily-closed).
func (m *Module) gateReason(r *http.Request, code, lang string) (allowed bool, msg string) {
	f := m.flags.Flag(code)
	switch f.Status {
	case svcOn:
		return true, ""
	case svcInviteOnly:
		if m.isInvited(r) {
			return true, ""
		}
		if s := f.Message(lang); s != "" {
			return false, s
		}
		return false, T(lang, "svc.invite_note")
	default: // maintenance | off
		if s := f.Message(lang); s != "" {
			return false, s
		}
		return false, T(lang, "svc.closed_note")
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

// addressedTo returns the places whose material this reader should be shown:
// the place they gave in their profile and everything that contains it.
//
// Empty for a guest and for anybody who never filled the field in, and then the
// feeds carry only material written for everyone.
func (m *Module) addressedTo(r *http.Request) []uuid.UUID {
	id, ok := m.authorID(r)
	if !ok || m.geo == nil {
		return nil
	}
	place, err := m.geo.UserPlace(r.Context(), id)
	if err != nil || place == nil {
		return nil
	}
	places, err := m.geo.AddressedTo(r.Context(), *place)
	if err != nil {
		m.rt.Logger.Warn("addressed places", zap.Error(err))
		return nil
	}
	return places
}

// latestNews returns the newest published articles for the sidebar carousel.
func (m *Module) latestNews(r *http.Request, lang string, n int) []FeedItem {
	arts, err := m.store.ListPublished(r.Context(), "", "", "", n, 0, m.addressedTo(r))
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
			Slug: a.Slug, Title: tr.Title, AuthorName: name, AuthorID: a.AuthorID.String(), AIAuthor: aiAuthor,
			ServedLang: served, Category: a.Category, CoverURL: a.CoverURL,
		})
	}
	return out
}

// seoFilterQuery returns the indexable filters for the request as a stable,
// ordered query string (no lang), whitelisted per page. Only these filters
// survive into canonical/hreflang URLs, so arbitrary query combinations (search,
// paging, sort) never spawn duplicate indexable pages.
func seoFilterQuery(r *http.Request) string {
	q := r.URL.Query()
	var parts []string
	switch r.URL.Path {
	case "/":
		if c := q.Get("cat"); IsCategory(c) {
			parts = append(parts, "cat="+url.QueryEscape(c))
			// A subcategory is only meaningful under its category.
			if s := q.Get("sub"); IsSubcategory(s) {
				parts = append(parts, "sub="+url.QueryEscape(s))
			}
		}
		// Page two holds different articles from page one and has to say so;
		// canonicalising it to "/" tells a crawler the two are the same page and
		// that everything past the first twenty-one can be ignored. The "top"
		// ordering is left out on purpose: it is the same corpus in another
		// order, and it already points at the default view.
		if q.Get("sort") != "top" {
			if n := pageParam(r); n > 1 {
				parts = append(parts, "page="+strconv.Itoa(n))
			}
		}
	case "/listings":
		if d := q.Get("deal"); isDealType(d) {
			parts = append(parts, "deal="+url.QueryEscape(d))
		}
		if t := q.Get("type"); isPropertyType(t) {
			parts = append(parts, "type="+url.QueryEscape(t))
		}
	}
	return strings.Join(parts, "&")
}

// canonURL builds the canonical relative URL for a page in one language,
// preserving the whitelisted filters so /?cat=sport canonicalizes to itself
// (with its category), not to a bare "/".
func canonURL(path, filters, lang string) string {
	q := "lang=" + lang
	if filters != "" {
		q = filters + "&" + q
	}
	return path + "?" + q
}

// langLinks builds the per-language alternates for the current page, carrying
// the same whitelisted filters so switching language keeps the category/filter.
func langLinks(base, filters string) map[string]string {
	out := make(map[string]string, len(Langs))
	for _, l := range Langs {
		out[l] = canonURL(base, filters, l)
	}
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
	AuthorID       string // for the byline link to /author/{id}
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

	// OrgName is the verified organisation this was published on behalf of.
	// A card shows it instead of the person: on a place page the reader is
	// looking for the akimat and the utility, and "А. Смағұлова" hides exactly
	// the fact the whole feature exists to show. The person is still named in
	// full on the article itself.
	OrgName string
}

// withOrgs fills in the organisation behind each card, in one query for the
// whole feed rather than one per card.
func (m *Module) withOrgs(ctx context.Context, arts []*Article, items []FeedItem) []FeedItem {
	if m.orgs == nil || len(items) == 0 {
		return items
	}
	ids := make([]uuid.UUID, 0, len(arts))
	for _, a := range arts {
		ids = append(ids, a.AuthorID)
	}
	names, err := m.orgs.VerifiedNames(ctx, ids)
	if err != nil || len(names) == 0 {
		return items
	}
	for i := range items {
		id, err := uuid.Parse(items[i].AuthorID)
		if err != nil {
			continue
		}
		if name := names[id]; name != "" {
			items[i].OrgName = name
		}
	}
	return items
}

// feedItems turns articles into the cards a template renders, in the reader's
// language. Extracted when the article page needed the same cards for its
// "read next" block: one shape of card, built one way, so the two lists cannot
// drift apart.
func feedItems(arts []*Article, lang string) []FeedItem {
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
			AuthorID:       a.AuthorID.String(),
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
	return items
}

// HomePage is the template context for the portal home.
type HomePage struct {
	Base
	Featured   *FeedItem
	Posts      []FeedItem
	Recent     []FeedItem
	Subscribed bool

	// The feed shows 21 articles and there are several times that many. Until
	// now the rest were reachable only by search, a category filter that also
	// stopped at 21, or a direct link — a hundred pages with no path to them
	// from the site, which is also how a crawler decides they do not matter.
	Page    int
	PrevURL string // empty on the first page
	NextURL string // empty on the last

	// Notice is where a reader lands after their report hid an article: the
	// piece itself is no longer readable, so the acknowledgement has to appear
	// somewhere they can still see it.
	Notice string
}

// homePageSize is how many articles one page of the feed holds.
const homePageSize = 21

// pageParam reads ?page=, clamping to a sane range. Anything unreadable is
// page 1: a bad page number is a reason to show the feed, not an error.
func pageParam(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || n < 1 {
		return 1
	}
	if n > 1000 {
		return 1000
	}
	return n
}

// feedURL builds a link to another page of the same feed, keeping whatever the
// reader is currently filtering and sorting by.
func feedURL(r *http.Request, lang string, page int) string {
	q := url.Values{}
	q.Set("lang", lang)
	src := r.URL.Query()
	if c := src.Get("cat"); IsCategory(c) {
		q.Set("cat", c)
		if sub := src.Get("sub"); IsSubcategory(sub) {
			q.Set("sub", sub)
		}
	}
	if src.Get("sort") == "top" {
		q.Set("sort", "top")
	}
	if page > 1 {
		q.Set("page", strconv.Itoa(page))
	}
	return "/?" + q.Encode()
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
		// Effective content: the admin-editable DB override, falling back to the
		// built-in default so a page never blanks out.
		title, body := m.pageContent(r.Context(), key, lang)
		if title == "" {
			http.NotFound(w, r)
			return
		}
		page := StaticPage{Base: m.base(r, title, lang)}
		page.Body = RenderMarkdown(applyOperator(body, m.rt.Config.Operator, lang))
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

	pageNo := pageParam(r)
	// One more than the page holds: if it comes back, there is a next page.
	// Cheaper than a COUNT(*) over the published set on every home request.
	arts, err := m.store.ListPublished(r.Context(), sort, cat, sub,
		homePageSize+1, (pageNo-1)*homePageSize, m.addressedTo(r))
	if err != nil {
		m.rt.Logger.Error("home list", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	hasNext := len(arts) > homePageSize
	if hasNext {
		arts = arts[:homePageSize]
	}

	items := m.withOrgs(r.Context(), arts, feedItems(arts, lang))

	page := HomePage{Base: m.base(r, T(lang, "home.page_title"), lang)}
	page.Active = active
	page.ActiveCat = cat
	page.ActiveSub = sub
	page.Subscribed = r.URL.Query().Get("subscribed") == "ok"
	if r.URL.Query().Get("reported") == "hidden" {
		page.Notice = T(lang, "article.report_hidden")
	}
	page.Posts = items
	page.Recent = recentSlice(items, 5)
	page.Page = pageNo
	if pageNo > 1 {
		page.PrevURL = feedURL(r, lang, pageNo-1)
		// A page past the end is a real URL a crawler can reach — from a link
		// that was valid when the archive was longer, or from a guess. It must
		// not join the index as an empty page claiming to be about the site.
		if len(items) == 0 {
			page.NoIndex = true
		}
	}
	// A translation the reader's language is missing drops the article from
	// items, so a full page can render short — but the next page still exists.
	if hasNext {
		page.NextURL = feedURL(r, lang, pageNo+1)
	}
	page.SidebarNews = m.latestNews(r, lang, 6)
	m.render(w, "home", page)
}

// ArticlePage is the template context for a single article.
type ArticlePage struct {
	Base
	ArticleID     string
	Slug          string
	Title         string
	Summary       string
	AuthorName    string
	AuthorID      string // for the byline link to /author/{id}
	ServedLang    string
	RequestedLang string
	Body          interface{}
	Published     *time.Time
	Updated       *time.Time // last edit, for dateModified in the article JSON-LD
	Views         int64
	IsAI          bool
	AIAuthor      bool
	Translated    bool

	// OrgName, OrgKind and OrgOfficial describe the organisation this was
	// published on behalf of, when the account has a verified one. Empty
	// otherwise, and empty is the only thing an unverified application shows.
	OrgName        string
	OrgKind        string
	OrgOfficial    bool
	AvailableLangs []string

	// CiteLine is the ready-made reference a reader can copy, empty on the
	// articles we do not offer as a source. See citeLine.
	CiteLine string

	// Related are the pieces offered at the end of this one. Until they
	// existed a reader who finished an article had nowhere to go, and the
	// article itself was a leaf with no link pointing out of it.
	Related []FeedItem

	// Predictions are the forecasts made in this piece, with what became of
	// them. Empty for the articles that made none, which is most of them.
	Predictions []*Prediction
	PredScore   PredictionScore

	Category    string
	Subcategory string
	CoverURL    string
	Score       int
	UserVote    int // -1, 0, +1
	AuthorKarma int
	CanVote     bool // logged in and not the author
	IsAuthor    bool
	Recent      []FeedItem // reserved for sidebar
	Subscribed  bool
	Comments    []Comment
	IsFavorite  bool
	// CanReport is false for a guest and for the author: a guest has no standing
	// to weigh, and reporting your own article is not a thing.
	CanReport bool
	Reported  bool
	// Notice is the one-line feedback after an action on this page — a report
	// accepted, or the verified-email bar explaining why it was not.
	Notice        string
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
	//
	// Crawlers are skipped, and that is not a detail. They used to be counted:
	// two thirds of every "views" number on the dashboard were Googlebot, the
	// Facebook link scraper and AI crawlers. That inflated the counter itself,
	// and — worse — it was the denominator of the reading-depth funnel, whose
	// numerator only a real browser can produce (the beacon needs JavaScript).
	// So genuine 23% read-through was reported as 2%, and an author reads that
	// as "nobody finishes my articles". Crawler traffic is not lost: the
	// analytics panel counts it under "bots".
	counted := botLabel(r.UserAgent()) == ""
	if counted {
		if err := m.store.RecordView(r.Context(), a.ID, served); err != nil {
			m.rt.Logger.Warn("record view", zap.Error(err))
		}
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
	page.AuthorID = a.AuthorID.String()
	page.ServedLang = served
	page.RequestedLang = lang
	page.Body, page.TOC = RenderMarkdownTOC(tr.BodyMD)
	page.ReadingMin = readingMinutes(tr.BodyMD)
	page.Published = a.PublishedAt
	if !a.UpdatedAt.IsZero() {
		u := a.UpdatedAt
		page.Updated = &u
	}
	page.Views = a.ViewsCount
	if counted {
		page.Views++ // the row was just bumped; show the reader their own visit
	}
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
	page.CanReport = page.CanVote
	if viewer != uuid.Nil {
		page.IsFavorite = m.favs.IsFavorite(r.Context(), viewer, "article", a.ID)
		page.Reported = m.store.HasReportedArticle(r.Context(), a.ID, viewer)
	}

	switch {
	case r.URL.Query().Get("reported") == "ok":
		page.Notice = T(lang, "article.report_thanks")
	case r.URL.Query().Get("notice") == "verify":
		page.Notice = T(lang, "article.report_verify")
	}
	page.CommentReview = r.URL.Query().Get("comment") == "review"
	// The article page shows only its table of contents in the aside (no news
	// carousel / widgets), so SidebarNews is intentionally not populated here.
	if cs, err := m.comments.ListForArticle(r.Context(), a.ID, viewer); err == nil {
		page.Comments = cs
	} else {
		m.rt.Logger.Warn("load comments", zap.Error(err))
	}

	// Non-indexable articles still render in full — they are only kept out of
	// search, not out of the site. See migration 20251107009900.
	page.NoIndex = !a.Indexable
	// A verified organisation replaces the byline; the person stays underneath,
	// because somebody has to be answerable for the text.
	if org, err := m.orgs.VerifiedByUser(r.Context(), a.AuthorID); err == nil && org != nil {
		page.OrgName = org.Name
		page.OrgKind = org.KindLabelKey()
		page.OrgOfficial = org.Official()
	}
	if rel, err := m.store.RelatedPublished(r.Context(), a.ID, a.Category, a.Subcategory, 4, m.addressedTo(r)); err == nil {
		page.Related = m.withOrgs(r.Context(), rel, feedItems(rel, page.Lang))
	} else {
		m.rt.Logger.Warn("related articles", zap.Error(err))
	}
	if preds, err := m.predictions.ForArticle(r.Context(), page.Lang, a.ID); err == nil {
		page.Predictions = preds
		// The ledger's running accuracy travels with the block: "open" means
		// more when the reader can see how the resolved ones turned out.
		if len(preds) > 0 {
			if sc, serr := m.predictions.Score(r.Context()); serr == nil {
				page.PredScore = sc
			}
		}
	} else {
		m.rt.Logger.Warn("article predictions", zap.Error(err))
	}

	// Only the human-written articles offer a citation: the whole point of the
	// block is to make our name easy to credit, and the AI columns are the one
	// thing we do not want credited to it.
	if !page.NoIndex {
		page.CiteLine = citeLine(page.Lang, page.AuthorName, page.Title,
			page.SiteURL, "/read/"+page.Slug, page.Published)
	}
	if page.NoIndex {
		// The meta tag in the template speaks to search engines. This header also
		// reaches the AI crawlers, which do not parse a robots meta at all, and it
		// arrives on a HEAD request where there is no body to read. robots.txt
		// says the same thing a third time; between them one of the three lands.
		w.Header().Set("X-Robots-Tag", aiRobotsTag)
	}

	m.applyArticleSEO(&page)
	m.render(w, "article", page)
}

// handleVote records a reader's up/down vote (toggling off when the same
// direction is submitted twice), then returns to the article.
// handleCommentDelete lets a reader take back their own comment. Until now the
// only way to unsay something published under your real name was to write to an
// administrator, which is not a feature — it is a missing one.
func (m *Module) handleCommentDelete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	userID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := m.comments.Delete(r.Context(), id, userID); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("delete comment", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/read/"+slug+"#comments", http.StatusSeeOther)
}

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
	// Comments gate (staged launch). The form is hidden when closed; this is the
	// backstop for a direct POST.
	if ok, _ := m.gateReason(r, SvcComments, m.resolveLang(w, r)); !ok {
		http.Redirect(w, r, backTo, http.StatusSeeOther)
		return
	}
	a, err := m.store.GetPublishedBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// A comment is published as written. It used to pass a model first, which
	// decided whether anyone would see it; readers decide that now, by voting,
	// and a comment voted far enough down folds away instead of vanishing.
	if err := m.comments.Create(r.Context(), a.ID, userID, body); err != nil {
		m.rt.Logger.Error("create comment", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
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

// handleCommentVote records a reader's up or down vote on a comment.
//
// The same rules as voting on an article: one vote per reader, weighted by the
// voter's own karma, and clicking the same arrow twice takes the vote back.
// What it replaces is a model that read every comment and decided whether it
// could be seen at all.
func (m *Module) handleCommentVote(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	backTo := "/read/" + slug + "#comments"

	voter, ok := m.authorID(r)
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
	var value int
	switch r.FormValue("value") {
	case "1", "up":
		value = ratings.VoteUp
	case "-1", "down":
		value = ratings.VoteDown
	default:
		http.Redirect(w, r, backTo, http.StatusSeeOther)
		return
	}
	// Clicking the current direction again retracts the vote.
	if cur, err := m.ratings.CommentVote(r.Context(), id, voter); err == nil && cur == value {
		value = ratings.VoteNone
	}

	if _, err := m.ratings.VoteComment(r.Context(), id, voter, value); err != nil {
		switch {
		case errors.Is(err, ratings.ErrOwnComment), errors.Is(err, ratings.ErrNotFound):
			http.Redirect(w, r, backTo, http.StatusSeeOther)
			return
		}
		m.rt.Logger.Error("vote comment", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, backTo, http.StatusSeeOther)
}

// ---------- auth pages ----------

// FormPage backs the login and register screens.
type FormPage struct {
	Base
	Mode   string // login | register
	Email  string
	First  string // given name (register)
	Last   string // family name (register)
	Middle string // patronymic, optional (register)
	Error  string
	Notice string
	Ref    string
	Next   string // where to go after a successful login/registration

	// PlaceID keeps the chosen location when the form comes back with an
	// error, so a rejected password does not also cost the reader their
	// place selection.
	PlaceID string
}

// safeNext returns a post-login destination taken from the request, or "" if it
// is anything other than a path on this site.
//
// Only a single leading slash passes. "//evil.example" is a protocol-relative
// URL that browsers follow off-site, and a login page that forwards to an
// attacker-chosen address after authenticating is a phishing primitive — the
// victim sees our domain, our form, our certificate, and lands somewhere else
// already trusting the page. Backslashes and CR/LF are refused for the same
// family of reasons: some clients normalise "\\" to "/", and a newline in a
// Location header splits the response.
func safeNext(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") {
		return ""
	}
	if strings.ContainsAny(v, "\\\r\n") {
		return ""
	}
	return v
}

func (m *Module) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	page := FormPage{Base: m.base(r, T(lang, "form.login_title"), lang), Mode: "login"}
	switch r.URL.Query().Get("verified") {
	case "ok":
		page.Notice = T(lang, "form.verified_ok")
	case "invalid":
		page.Error = T(lang, "form.verified_invalid")
	}
	// A protected action (e.g. publishing a listing) bounced here because the
	// session had lapsed — say so, so the user understands the failure.
	if r.URL.Query().Get("reason") == "session_expired" {
		page.Notice = T(lang, "form.session_expired")
	}
	// Arriving here on the way to somewhere else — a printed QR code aimed at
	// the posting form is the case this was built for. Carry the destination
	// through the form, and say why the account is being asked for at all: a
	// stranger who scanned a poster about free listings has no reason to guess.
	page.Next = safeNext(r.URL.Query().Get("next"))
	if page.Next == "/listings/new" && page.Notice == "" {
		page.Notice = T(lang, "form.next_listing")
	}
	m.render(w, "form", page)
}

func (m *Module) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	m.render(w, "form", FormPage{
		Base: m.base(r, T(lang, "form.register_title"), lang),
		Mode: "register",
		Ref:  strings.TrimSpace(r.URL.Query().Get("ref")),
		Next: safeNext(r.URL.Query().Get("next")),
	})
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

	// A second factor the main entrance ignores is not a second factor. This
	// form has no challenge step, so when MFA is configured it must refuse
	// rather than hand out a session that skipped it. Today TOTP is off and
	// nobody sees this; the check exists so that turning it on cannot silently
	// leave the browser door open. Checked before the password so a closed door
	// cannot be used to probe which e-mails exist.
	if m.auth.MFAEnabled() {
		m.render(w, "form", FormPage{
			Base:  m.base(r, T(lang, "form.login_title"), lang),
			Mode:  "login",
			Email: email,
			Error: T(lang, "form.err_mfa_web"),
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
			Next:  safeNext(r.FormValue("next")),
		})
		return
	}
	auth.SetSessionCookie(w, r, token, m.auth.SessionTTL())
	m.rt.Logger.Info("studio login", zap.String("user_id", user.ID.String()))
	http.Redirect(w, r, afterAuth(r), http.StatusSeeOther)
}

func (m *Module) handleRegisterSubmit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	first := auth.NormalizePersonName(r.FormValue("first_name"))
	last := auth.NormalizePersonName(r.FormValue("last_name"))
	middle := auth.NormalizePersonName(r.FormValue("middle_name"))
	ref := strings.TrimSpace(r.FormValue("ref"))

	// One re-render path: whatever fails, the user gets their input back.
	place := strings.TrimSpace(r.FormValue("geo_node_id"))
	regFail := func(msg string) {
		m.render(w, "form", FormPage{
			Base: m.base(r, T(lang, "form.register_title"), lang), Mode: "register",
			Email: email, First: first, Last: last, Middle: middle, Ref: ref, Error: msg,
			PlaceID: place,
		})
	}
	switch m.flags.Flag(SvcRegistration).Status {
	case svcOn:
		// open to everyone
	case svcInviteOnly:
		if _, ok := m.refs.ReferrerByCode(r.Context(), ref); !ok {
			regFail(T(lang, "form.err_invite_only"))
			return
		}
	default: // maintenance | off
		msg := m.flags.Flag(SvcRegistration).Message(lang)
		if msg == "" {
			msg = T(lang, "svc.closed_note")
		}
		regFail(msg)
		return
	}

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
		regFail(T(lang, "form.err_email_invalid"))
		return
	}
	if err := auth.ValidatePassword(password); err != nil {
		regFail(T(lang, "form.err_password_rule"))
		return
	}
	// Real name is required for everyone — readers, authors and agents alike —
	// so attribution on comments and articles is a person, not an e-mail.
	if err := auth.ValidatePersonName(first); err != nil {
		regFail(T(lang, "form.err_first_name"))
		return
	}
	if err := auth.ValidatePersonName(last); err != nil {
		regFail(T(lang, "form.err_last_name"))
		return
	}
	if err := auth.ValidateOptionalPersonName(middle); err != nil {
		regFail(T(lang, "form.err_middle_name"))
		return
	}

	user, token, err := m.auth.RegisterPassword(r.Context(), email, password, first, last, middle)
	if err == nil && place != "" {
		// Best-effort: a place that will not parse or will not save must never
		// cost somebody the account they have just created. They can set it
		// again from the profile, and until they do their feed carries only
		// what was written for everyone.
		if id, perr := uuid.Parse(place); perr == nil {
			if serr := m.geo.SetUserPlace(r.Context(), user.ID, &id); serr != nil {
				m.rt.Logger.Warn("save place on register", zap.Error(serr))
			}
		}
	}
	if err != nil {
		msg := T(lang, "form.err_generic")
		if errors.Is(err, auth.ErrEmailExists) {
			msg = T(lang, "form.err_email_taken")
		} else if errors.Is(err, auth.ErrInvalidEmail) {
			msg = T(lang, "form.err_email_invalid")
		}
		regFail(msg)
		return
	}
	// Where this account registered from, resolved once from the request IP.
	// The admin register needs a country beside each name, and this is the only
	// moment it can be known without following anyone around: no per-visit
	// history is kept, and an unresolvable address simply leaves it blank.
	if m.geoip != nil {
		if cc := m.geoip.country(clientIP(r)); cc != "" {
			if err := m.users.SetSignupCountry(r.Context(), user.ID, cc); err != nil {
				m.rt.Logger.Warn("signup country", zap.Error(err))
			}
		}
	}
	// Record the consent the checkbox represents (append-only proof).
	if err := m.auth.RecordConsent(r.Context(), r, user.ID, "web"); err != nil {
		m.rt.Logger.Error("record consent (web)", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	// Send the email-verification link (best effort).
	if err := m.auth.IssueEmailVerification(r.Context(), user.ID, user.Email); err != nil {
		m.rt.Logger.Warn("issue email verification (web)", zap.String("user_id", user.ID.String()), zap.Error(err))
	}
	// Referral capture: if the registration carried an invite code, link the
	// new user to their referrer. Best-effort — a bad code must not fail signup.
	if code := strings.TrimSpace(r.FormValue("ref")); code != "" {
		if referrer, ok := m.refs.ReferrerByCode(r.Context(), code); ok {
			if err := m.refs.RecordReferral(r.Context(), referrer, user.ID); err != nil {
				m.rt.Logger.Warn("record referral", zap.Error(err))
			}
		}
	}
	auth.SetSessionCookie(w, r, token, m.auth.SessionTTL())
	m.rt.Logger.Info("studio register", zap.String("user_id", user.ID.String()))
	http.Redirect(w, r, afterAuth(r), http.StatusSeeOther)
}

// afterAuth is where a fresh session lands: the destination the user was
// heading for before the login wall, or the studio when there was none.
//
// Without this every route through the wall ended at /studio, which is the
// article studio — so someone who scanned a poster offering a free property
// listing finished registration staring at a publishing dashboard, with the
// thing they came to do nowhere in sight.
func afterAuth(r *http.Request) string {
	if n := safeNext(r.FormValue("next")); n != "" {
		return n
	}
	return "/studio"
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
	// L25 and L100 are the listeners inside the figures above: how many started
	// the recording and how many heard it out.
	L25, L100 int64
	// PFinish is what share of the readers who STARTED the article finished it.
	// Unlike the percentages above it has no view count in the denominator, so
	// it stays honest whatever the view counter is doing — and it answers the
	// question an author actually has: was the piece worth staying with?
	PFinish int
}

// analyticsSince is the day the view and reading-depth counters were reset to
// zero together, after crawler hits were found in the view counts. Shown in the
// studio so nobody reads a small number as a collapse in readership — and so the
// next person to look does not have to reconstruct why the history is short.
const analyticsSince = "08.08.2026"

// StudioPage is the dashboard context.
type StudioPage struct {
	Base
	Stats AuthorStats
	Karma int
	// Since is the date the counters start from, shown above the table.
	Since    string
	Articles []StudioRow
	// Outcome of the last publish attempt, so the author is told what happened
	// instead of being returned to an unchanged-looking dashboard.
	Notice string
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
			// The funnel counts everyone who got that far, by either route.
			row.D25 = depthTotal(d, 25)
			row.D50 = depthTotal(d, 50)
			row.D75 = depthTotal(d, 75)
			row.D100 = depthTotal(d, 100)
			row.P25 = pctOf(row.D25, row.Views)
			row.P50 = pctOf(row.D50, row.Views)
			row.P75 = pctOf(row.D75, row.Views)
			row.P100 = pctOf(row.D100, row.Views)
			row.PFinish = pctOf(row.D100, row.D25)
			// And separately: of those, how many were listening. Shown beside
			// the funnel rather than inside it, because it answers a different
			// question -- not how far, but by which route.
			if l := d[ModeListen]; l != nil {
				row.L25, row.L100 = l[25], l[100]
			}
		}
		rows = append(rows, row)
	}

	karma, err := m.ratings.AuthorKarma(r.Context(), authorID)
	if err != nil {
		m.rt.Logger.Warn("author karma", zap.Error(err))
	}

	page := StudioPage{
		Base:  m.base(r, T(lang, "studio.title"), lang),
		Since: analyticsSince,
	}
	switch r.URL.Query().Get("ok") {
	case "published":
		page.Notice = T(lang, "studio.n_published")
	case "in_review":
		page.Notice = T(lang, "studio.n_review")
	}
	// A deleted draft leaves no trace in the table, so say so explicitly —
	// otherwise the author cannot tell a successful delete from a silent failure.
	if r.URL.Query().Get("deleted") == "1" {
		page.Notice = T(lang, "studio.deleted_ok")
	}
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

	// PlaceID is the place this article was published for, empty for "everyone".
	PlaceID string

	// CanTranslate is whether the site offers to translate this article. It is
	// separate from AIEnabled because the assistant stayed on for moderation
	// while automatic translation was switched off: authors have models of
	// their own and translate in a minute what this took eight to do.
	CanTranslate bool

	// TargetLangs are the languages this article would be translated into —
	// every language except the original. The button used to promise "three
	// languages", which is one more than it does: the original is already
	// written.
	TargetLangs []string

	// TranslateState is what happened to the last translation run: "running",
	// "done", "failed" or empty when none was ever started. Without it the
	// author clicks the button, sees "started", and then sits in front of empty
	// tabs with no way to learn that the job failed three times.
	TranslateState string
	TranslateError string

	// TranslationIssues lists, per language, what the translation lost against
	// the original — counted, not read. An author who does not know the target
	// language cannot check the meaning, but they can be told that a number
	// went missing or a link disappeared.
	TranslationIssues map[string][]TranslationIssue
}

// handleTranslateStatus answers the editor's poll: is the translation still
// running, did it finish, did it fail.
//
// Without it the author has no way to learn that anything happened. The button
// said "a few seconds", the tabs stayed empty, and the only record of three
// failures was in a log the author cannot read.
func (m *Module) handleTranslateStatus(w http.ResponseWriter, r *http.Request) {
	authorID, ok := m.authorID(r)
	if !ok {
		http.Error(w, "auth required", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// The state of someone else's article is none of this caller's business.
	if _, err := m.store.GetByID(r.Context(), id, authorID); err != nil {
		http.NotFound(w, r)
		return
	}
	state, msg := m.translationRun(r.Context(), id.String())
	writeJSONObj(w, map[string]any{"state": state, "error": msg})
}

// otherLangs lists every language except the original, in the site's order.
func otherLangs(original string) []string {
	out := make([]string, 0, len(Langs)-1)
	for _, l := range Langs {
		if l != original {
			out = append(out, l)
		}
	}
	return out
}

// translationRun reports how the most recent translation of this article ended.
// Best-effort: the editor renders with or without it.
func (m *Module) translationRun(ctx context.Context, articleID string) (state, msg string) {
	if m.rt == nil || m.rt.DB == nil {
		return "", ""
	}
	var status string
	var lastErr *string
	err := m.rt.DB.QueryRow(ctx, `
		SELECT status::text, last_error FROM job_queue
		 WHERE name = $1 AND payload->>'article_id' = $2
		 ORDER BY created_at DESC LIMIT 1`, ai.JobTranslate, articleID).Scan(&status, &lastErr)
	if err != nil {
		return "", ""
	}
	switch status {
	case "pending", "running", "retry":
		return "running", ""
	case "done":
		return "done", ""
	case "failed":
		if lastErr != nil {
			return "failed", *lastErr
		}
		return "failed", ""
	}
	return "", ""
}

// aiNotice maps an ?ai= redirect flag to a localized message.
func aiNotice(lang, flag string) string {
	// Translation is the only AI action an author can start now: the co-editor
	// and the column drafter are gone, and with them their flags.
	switch flag {
	case "queued":
		return T(lang, "notice.ai_queued")
	case "off":
		return T(lang, "notice.ai_off")
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
	page.TargetLangs = otherLangs(lang)
	page.Category = "society"
	page.Subcategory = ""
	page.Status = "draft"
	page.Fields = emptyFields()
	page.AIEnabled = m.ai.Enabled()
	page.CanTranslate = m.ai.AutoTranslateEnabled()
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
	// The publish button sends the author back here when a local piece is
	// missing one of the two languages, so the reason has to be visible on
	// arrival rather than left for them to guess.
	// The publish button sends the author back here when a version is missing.
	// Which one is stashed in the query, because "not all languages" sends them
	// hunting through their own piece for the gap.
	if r.URL.Query().Get("err") == "languages" {
		e := &MissingLanguagesError{
			Placed:  r.URL.Query().Get("scope") == "local",
			Missing: strings.Split(r.URL.Query().Get("miss"), ","),
		}
		page.Error = e.Reason(lang)
	}
	page.ArticleID = a.ID.String()
	page.Slug = a.Slug
	page.Status = a.Status
	page.OriginalLang = a.OriginalLang
	page.TargetLangs = otherLangs(a.OriginalLang)
	page.TranslateState, page.TranslateError = m.translationRun(r.Context(), a.ID.String())
	if orig, ok := a.Translations[a.OriginalLang]; ok {
		issues := map[string][]TranslationIssue{}
		for _, l := range otherLangs(a.OriginalLang) {
			tr, ok := a.Translations[l]
			if !ok {
				continue
			}
			// All three fields, not just the body: a summary came back three
			// times its original length with the model's own prompt echoed into
			// it, and a body-only check would have called that clean.
			var found []TranslationIssue
			for _, part := range []struct {
				field, src, dst string
			}{
				{"title", orig.Title, tr.Title},
				{"summary", orig.Summary, tr.Summary},
				{"body", orig.BodyMD, tr.BodyMD},
			} {
				for _, issue := range compareTranslation(part.src, part.dst) {
					issue.Field = part.field
					found = append(found, issue)
				}
			}
			if len(found) > 0 {
				issues[l] = found
			}
		}
		if len(issues) > 0 {
			page.TranslationIssues = issues
		}
	}
	// "Started, refresh in a minute" and "finished" were shown side by side:
	// the first comes from the redirect that is still in the address bar, the
	// second from the state. Once the state can speak for itself, it does.
	if page.TranslateState != "" && page.TranslateState != "running" {
		page.Notice = ""
	}
	page.Category = a.Category
	page.Subcategory = a.Subcategory
	page.CoverURL = a.CoverURL
	page.Fields = fields
	page.AIEnabled = m.ai.Enabled()
	page.CanTranslate = m.ai.AutoTranslateEnabled()
	if node, err := m.store.ArticlePlace(r.Context(), a.ID); err == nil && node != nil {
		page.PlaceID = node.String()
	}
	page.Notice = aiNotice(lang, r.URL.Query().Get("ai"))
	m.render(w, "studio_editor", page)
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

	if !m.ai.AutoTranslateEnabled() {
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

// sanitizeCoverURL keeps a cover reference only when it is a safe image
// location: an app-relative upload path (/media/… or /static/…) or an http(s)
// URL. Anything else — including javascript:/data: schemes — is dropped, so the
// value can be placed in an <img src> without becoming an injection vector.
func sanitizeCoverURL(v string) string {
	v = strings.TrimSpace(v)
	switch {
	case v == "":
		return ""
	case strings.HasPrefix(v, "/media/") || strings.HasPrefix(v, "/static/"):
		return v
	case strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://"):
		return v
	}
	return ""
}

// parseEditorForm reads the three-language editor form into translation inputs.
func parseEditorForm(r *http.Request) (originalLang, category, subcategory, coverURL string, trs []TranslationInput) {
	originalLang = r.FormValue("original_lang")
	if !IsLang(originalLang) {
		originalLang = LangRU
	}
	category = NormalizeCategory(r.FormValue("category"))
	subcategory = NormalizeSubcategory(category, r.FormValue("subcategory"))
	coverURL = sanitizeCoverURL(r.FormValue("cover_url"))
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

// formPlace reads the place chosen in the editor. An unparseable value is
// treated as "no place" rather than as an error: a bad hidden field must not
// cost an author their text.
func formPlace(r *http.Request) *uuid.UUID {
	raw := strings.TrimSpace(r.FormValue("geo_node_id"))
	if raw == "" {
		return nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil
	}
	return &id
}

// placeAllowed refuses a place outside a verified organisation's territory.
//
// The rule only bites for organisations: a private author writes about wherever
// they like. An akimat of Kostanay publishing for Almaty is either a mistake or
// an abuse, and without this the verified badge becomes a pass to the whole
// country.
func (m *Module) placeAllowed(r *http.Request, author uuid.UUID, place *uuid.UUID) bool {
	if place == nil {
		return true
	}
	org, err := m.orgs.VerifiedByUser(r.Context(), author)
	if err != nil || org == nil {
		return true
	}
	ok, err := m.orgs.MayPublishFor(r.Context(), org, place)
	if err != nil {
		m.rt.Logger.Warn("org territory check", zap.Error(err))
		return true
	}
	return ok
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
	// Article submission gate (staged launch): open / invite-only / closed.
	// Leadership (the CEO/directors) bypass it so they can write during the beta.
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canAuthorAsStaff(claims) {
		if ok, msg := m.gateReason(r, SvcArticleSubmit, lang); !ok {
			m.reRenderEditor(w, r, true, "", "", originalLang, category, subcategory, coverURL, trs, msg)
			return
		}
	}
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
	if place := formPlace(r); m.placeAllowed(r, authorID, place) {
		if err := m.store.SetArticlePlace(r.Context(), id, authorID, place); err != nil {
			m.rt.Logger.Warn("set place on create", zap.Error(err))
		}
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

	if place := formPlace(r); m.placeAllowed(r, authorID, place) {
		if err := m.store.SetArticlePlace(r.Context(), id, authorID, place); err != nil {
			m.rt.Logger.Warn("set place on update", zap.Error(err))
		}
	}
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
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	// Leadership (the CEO/directors) publish under their own name without
	// email/phone verification or the consent step — their staff account already
	// establishes identity and they own the documents. Everyone else must have a
	// verified real-name identity and have accepted the current documents.
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canAuthorAsStaff(claims) {
		// Articles carry a real-name byline: publishing requires a verified author
		// identity (real name + verified phone).
		if !m.auth.CanPublish(r.Context(), authorID) {
			http.Redirect(w, r, "/studio/author?need=publish", http.StatusSeeOther)
			return
		}
		// Acknowledgment of the current documents and tariffs before publishing.
		// Fail-closed: a DB error must not let an unconsented article through.
		consented, err := m.auth.HasAuthorConsent(r.Context(), authorID)
		if err != nil {
			m.rt.Logger.Error("author consent check", zap.String("user_id", authorID.String()), zap.Error(err))
			http.Redirect(w, r, "/studio/consent", http.StatusSeeOther)
			return
		}
		if !consented {
			http.Redirect(w, r, "/studio/consent", http.StatusSeeOther)
			return
		}
	}
	// Publishing is now a submission, not a switch: the rules check runs first
	// and the article either goes out or comes back with the rules it failed.
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	lang := m.resolveLang(w, r)

	// Leadership publishes immediately — no review queue. They are the editorial
	// authority, so their piece goes live at once, logged as a staff publish
	// under their own name rather than the checker's, which is what the ledger
	// used to say.
	if canAuthorAsStaff(claims) {
		title := ""
		if a, e := m.store.GetByID(r.Context(), id, authorID); e == nil {
			if tr, _ := a.Translation(lang); tr != nil {
				title = tr.Title
			}
		}
		first, last, _ := m.auth.AuthorIdentity(r.Context(), authorID)
		if err := m.commitReview(r.Context(), id, authorID, title, "published", "approve", "rules_ok",
			humanActor(authorID, strings.TrimSpace(first+" "+last)), nil); err != nil {
			m.rt.Logger.Error("staff publish", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if err := m.syndicate.EnqueuePublish(r.Context(), m.jobs, id); err != nil {
			m.rt.Logger.Warn("enqueue publish", zap.Error(err))
		}
		http.Redirect(w, r, "/studio?ok=published", http.StatusSeeOther)
		return
	}

	// Reader moderation: with the pre-publication checker off, the author's
	// publish button publishes. Nobody stands between the author and the reader,
	// which is the whole point — and the reason the reporting path below has to
	// work, because it is now the only thing that does.
	//
	// When the checker is enabled the old route stands: an administrator who
	// turns it on has chosen pre-moderation, and an outage must not quietly
	// downgrade that choice to open publishing.
	if m.ai == nil || !m.ai.ReviewCheckEnabled() {
		if err := m.publishNow(r.Context(), id, authorID); err != nil {
			// A local notice missing one of the two languages is the author's
			// to fix, not an error to log: they are told which language and
			// sent back to the editor.
			var miss *MissingLanguagesError
			if errors.As(err, &miss) {
				scope := "all"
				if miss.Placed {
					scope = "local"
				}
				http.Redirect(w, r, "/studio/a/"+id.String()+
					"?err=languages&scope="+scope+"&miss="+strings.Join(miss.Missing, ","),
					http.StatusSeeOther)
				return
			}
			if errors.Is(err, ErrNotFound) {
				// Either not theirs, or readers have hidden it — the one state
				// the author cannot clear by pressing publish again.
				http.Redirect(w, r, "/studio/moderation?err=hidden", http.StatusSeeOther)
				return
			}
			m.rt.Logger.Error("publish", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if err := m.syndicate.EnqueuePublish(r.Context(), m.jobs, id); err != nil {
			m.rt.Logger.Warn("enqueue publish", zap.Error(err))
		}
		http.Redirect(w, r, "/studio?ok=published", http.StatusSeeOther)
		return
	}

	published, blocking, err := m.submitForReview(r.Context(), id, authorID, lang)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("submit for review", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	switch {
	case published:
		if err := m.syndicate.EnqueuePublish(r.Context(), m.jobs, id); err != nil {
			m.rt.Logger.Warn("enqueue publish", zap.Error(err))
		}
		http.Redirect(w, r, "/studio?ok=published", http.StatusSeeOther)
	case blocking > 0:
		http.Redirect(w, r, "/studio/moderation?ok=returned", http.StatusSeeOther)
	default:
		// The checker could not be reached; the article waits for a person.
		http.Redirect(w, r, "/studio?ok=in_review", http.StatusSeeOther)
	}
}

// ConsentPage backs the one-time author consent gate.
type ConsentPage struct {
	Base
	Error string
}

// handleConsent shows the author consent gate (documents + tariffs) required
// once before publishing.
func (m *Module) handleConsent(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if consented, err := m.auth.HasAuthorConsent(r.Context(), authorID); err == nil && consented {
		http.Redirect(w, r, "/studio", http.StatusSeeOther)
		return
	}
	m.render(w, "author_consent", ConsentPage{Base: m.base(r, T(lang, "consent.title"), lang)})
}

// handleConsentSubmit records the author's acknowledgment, then returns to the
// studio so they can publish.
func (m *Module) handleConsentSubmit(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	authorID, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	if r.FormValue("consent") != "on" {
		m.render(w, "author_consent", ConsentPage{Base: m.base(r, T(lang, "consent.title"), lang), Error: T(lang, "consent.error")})
		return
	}
	if err := m.auth.RecordAuthorConsent(r.Context(), r, authorID, "web"); err != nil {
		m.rt.Logger.Error("record author consent", zap.String("user_id", authorID.String()), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/studio", http.StatusSeeOther)
}

func (m *Module) handleUnpublish(w http.ResponseWriter, r *http.Request) {
	m.transition(w, r, "draft", "/studio")
}

// handleDeleteDraft removes an abandoned draft — the false start, the duplicate
// left behind by a lost session. Drafts only: see Store.DeleteDraft for why a
// published article has to be unpublished first.
func (m *Module) handleDeleteDraft(w http.ResponseWriter, r *http.Request) {
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
	if err := m.store.DeleteDraft(r.Context(), id, authorID); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		m.rt.Logger.Error("delete draft", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	m.rt.Logger.Info("draft deleted", zap.String("article_id", id.String()))
	http.Redirect(w, r, "/studio?deleted=1", http.StatusSeeOther)
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
	page.TargetLangs = otherLangs(originalLang)
	page.Category = category
	page.Subcategory = subcategory
	page.CoverURL = coverURL
	page.Status = "draft"
	page.Fields = fields
	page.AIEnabled = m.ai.Enabled()
	page.CanTranslate = m.ai.AutoTranslateEnabled()
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
