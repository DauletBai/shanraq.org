package site

import (
	"html/template"
	"time"
)

// Base carries fields shared by every page (consumed by the header/footer
// partials via Go's embedded-field promotion). The whole UI renders in Lang.
type Base struct {
	// Nonce authorises this page's inline scripts. The policy admits a script
	// tag only if it carries the value minted for this response, so a tag an
	// attacker injects into the markup has nothing to put here.
	Nonce     string
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

	// SEO fields (populated by the module's base(); pages may override).
	SiteURL  string // absolute origin, e.g. https://shanraq.org
	Path     string // request path, no query (used for nav active state)
	CanonURL string // relative canonical path+query for THIS language,
	//                        including only whitelisted indexable filters
	Desc    string        // meta description
	OGImage string        // absolute image URL for social previews
	OGType  string        // "website" | "article"
	JSONLD  template.HTML // structured data (schema.org), injected verbatim
	SiteLD  template.HTML // the site's own card: who publishes this and where
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

// ServiceOff reports whether a service's entry point should be disabled in
// the UI. An unknown/unconfigured code is treated as available, so a missing
// flag never hides a link. Exposed to templates as "svcOff": it greys out
// links/buttons that lead to a service the admin turned off or set to maintenance.
func ServiceOff(svc map[string]ServiceView, code string) bool {
	if v, ok := svc[code]; ok {
		return !v.On
	}
	return false
}

// ServiceMsg returns the localized "temporarily unavailable / by invitation"
// notice for a service, for use as a tooltip on the disabled entry point.
// Exposed to templates as "svcMsg".
func ServiceMsg(svc map[string]ServiceView, code string) string {
	if v, ok := svc[code]; ok {
		return v.Msg
	}
	return ""
}

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

// Ad is one creative in the sidebar slot. It is either a booked placement from
// the advertiser cabinet or a house slide selling the slot itself; the corner
// ribbon says which, so nothing ever implies a commercial relationship that does
// not exist.
type Ad struct {
	Image string // /static/... illustration
	Title string
	Price string
	Desc  string
	URL   string // click target
	// House marks a slide that advertises the slot itself rather than a paying
	// customer. Both kinds carry the "Реклама" ribbon — a house slide is still
	// advertising — but a house slide takes the quieter graphite panel and keeps
	// rel="sponsored" off a link that points back into our own site.
	House bool
}

// Rate is one currency's KZT rate for the info bar. Main+Last split the value
// so the last digit can be dropped on narrow phones (Main shown, Last hidden).
type Rate struct {
	Code string // USD / EUR / RUB
	Main string // value without its last character, e.g. "469.8"
	Last string // the last character, e.g. "3"
	Dir  string // "up" | "down" | ""
}

// SocialLink is one configured social profile shown in the bar.
type SocialLink struct {
	Name string // telegram | instagram | youtube | facebook (icon key)
	URL  string
}

// LiveSocial keeps only the profiles that actually lead somewhere. A "#" entry
// is a placeholder for a network we have not opened yet — acceptable as a hint
// in the header strip, never on a card whose entire purpose is to be clicked.
func LiveSocial(links []SocialLink) []SocialLink {
	out := make([]SocialLink, 0, len(links))
	for _, l := range links {
		if l.URL != "" && l.URL != "#" {
			out = append(out, l)
		}
	}
	return out
}

// InfoBarData is the per-request snapshot handed to templates.
type InfoBarData struct {
	Today string
	// TodayISO is the same day in machine form, for the link into the archive.
	TodayISO     string
	WeatherIcon  string // icon key, e.g. "wx_sun" ("" when unavailable)
	WeatherTemp  string // e.g. "+25°"
	WeatherPress string // atmospheric pressure, e.g. "742 мм" ("" when unavailable)
	// WeatherPlace names the town the temperature was taken in, when it is the
	// reader's own rather than the default city. Empty keeps the cell silent.
	WeatherPlace string
	Rates        []Rate // empty when unavailable
	Social       []SocialLink
	GitHub       string // repository URL, rendered in the footer only ("" hides it)
}

// CurSymbol maps a currency code to its symbol for the compact bar.
func CurSymbol(code string) string {
	switch code {
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "RUB":
		return "₽"
	default:
		return code
	}
}
