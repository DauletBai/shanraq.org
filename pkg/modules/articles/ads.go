package articles

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

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

// houseAds fills the sidebar slot when nothing is sold for this surface. Every
// slide sells the slot itself, which is the only honest thing an empty ad space
// can advertise.
//
// It replaced three invented car adverts with prices attached. Those implied
// advertisers who do not exist, and a made-up price beside a made-up car is the
// kind of prop that stops being obviously a prop the moment somebody screenshots
// it. A slot that says "this space is free" is worth more than one pretending to
// be sold.
//
// Each claim here is one the site can actually stand behind: the aside is
// sticky and never covers the text, no third-party ad network is loaded on any
// page, and /advertise really does show a rate card and a live availability
// calendar. Deliberately absent: audience figures. Most recorded traffic is
// crawlers and datacenter addresses, so any number quoted here would flatter.
func houseAds(lang string) []Ad {
	slides := []struct{ key, cta string }{
		// Our own product leads. It used to sit last of six, which meant the
		// footer -- which shows only the first -- never carried it at all, and
		// the carousel reached it after five turns. Unsold inventory is worth
		// more selling something we have than selling the space itself, and the
		// offer to advertise loses nothing by following one slide later.
		{"adam", "/adam"},
		{"free", "/advertise"},   // the direct offer
		{"beside", "/advertise"}, // the format: beside the text, never over it
		{"notrack", "/privacy"},  // no third-party tracking — a real difference
		{"price", "/advertise"},  // open rate card, no agency in between
		{"realty", "/advertise"}, // the audience nearest to a buying decision
	}
	out := make([]Ad, 0, len(slides))
	for _, s := range slides {
		out = append(out, houseSlide(lang, s.key, s.cta))
	}
	return out
}

// houseSlide builds one house slide from its string key and where it leads.
func houseSlide(lang, key, cta string) Ad {
	return Ad{
		URL:   cta,
		Title: T(lang, "house."+key+"_title"),
		Desc:  T(lang, "house."+key+"_desc"),
		Price: T(lang, "house."+key+"_cta"),
		House: true,
	}
}

// ownAds is what a page outside the sold surfaces carries.
//
// Only three surfaces are sold -- listings, articles, the home page -- so every
// other address returned no ads at all: no card in the aside, and a footer
// falling through to the generic "your ad could be here". That is most of the
// site by count. The thousand forecast pages alone carried nothing, and they
// are the pages a stranger arrives on.
//
// Nothing is sold there, so nothing is displaced: the space shows our own
// product. Two addresses are left out -- /adam, where the ad would point at the
// page the reader is already on, and /advertise, whose whole subject is selling
// the slot rather than filling it.
func ownAds(r *http.Request, lang string) []Ad {
	p := r.URL.Path
	if p == "/adam" || strings.HasPrefix(p, "/advertise") {
		return nil
	}
	return []Ad{houseSlide(lang, "adam", "/adam")}
}

// adSurfaceFor maps a request to the surface it belongs to, mirroring the
// surfaces sold in the advertiser cabinet.
func adSurfaceFor(r *http.Request) string {
	p := r.URL.Path
	switch {
	case strings.HasPrefix(p, "/listings"):
		return surfaceRealestate
	case strings.HasPrefix(p, "/read/"):
		return surfaceArticles
	case p == "/":
		if cat := r.URL.Query().Get("cat"); cat != "" && IsCategory(cat) {
			return adRubricSurface(cat)
		}
		return surfaceHome
	}
	return ""
}

// sidebarAds serves the paid placements booked for this page's zone, falling
// back to the house slides when none are running. It used to return nil on an
// unsold surface, which is why the slot never appeared at all and the sidebar
// filled the space with a carousel repeating the page's own headlines.
//
// A page with no surface at all -- a forecast, the rates, a static page -- gets
// our own product rather than nothing; see ownAds.
func (m *Module) sidebarAds(r *http.Request, lang string) []Ad {
	surface := adSurfaceFor(r)
	if surface == "" || m.ads == nil {
		return ownAds(r, lang)
	}
	orders, err := m.ads.ActiveBySurface(r.Context(), surface, lang, m.adPageGeo(r), 12)
	if err != nil {
		m.rt.Logger.Warn("sidebar ads", zap.Error(err))
		return houseAds(lang)
	}
	if len(orders) == 0 {
		return houseAds(lang) // unsold: the slot advertises itself
	}
	out := make([]Ad, 0, len(orders))
	for _, o := range orders {
		url := o.TargetURL
		if url == "" {
			url = "#"
		}
		// Price doubles as the button label. An order booked without one would
		// otherwise render a button reading just "→".
		cta := strings.TrimSpace(o.CTA)
		if cta == "" {
			cta = T(lang, "ad.learn_more")
		}
		out = append(out, Ad{Image: o.ImageURL, Title: o.Title, Desc: o.Body, Price: cta, URL: url})
	}
	return out
}

// adPageGeo is the page's own geography, used to decide which geo-targeted
// bookings belong on it.
//
// Only two kinds of page carry a place: the page of a place, and an article
// written for one. Everything else — the front page, a rubric, a listing —
// belongs to no particular geography and therefore carries only the bookings
// that were bought for none.
func (m *Module) adPageGeo(r *http.Request) uuid.UUID {
	if m.store == nil {
		return uuid.Nil
	}
	p := r.URL.Path
	var id uuid.UUID
	var err error
	switch {
	case strings.HasPrefix(p, "/place/"):
		slug := strings.Trim(strings.TrimPrefix(p, "/place/"), "/")
		if slug == "" || strings.Contains(slug, "/") {
			return uuid.Nil
		}
		id, err = m.store.PlaceIDBySlug(r.Context(), slug)
	case strings.HasPrefix(p, "/read/"):
		slug := strings.Trim(strings.TrimPrefix(p, "/read/"), "/")
		if slug == "" || strings.Contains(slug, "/") {
			return uuid.Nil
		}
		id, err = m.store.ArticlePlaceBySlug(r.Context(), slug)
	default:
		return uuid.Nil
	}
	if err != nil {
		m.rt.Logger.Warn("ad page geo", zap.Error(err))
		return uuid.Nil
	}
	return id
}
