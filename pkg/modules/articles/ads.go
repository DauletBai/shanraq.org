package articles

import (
	"net/http"
	"strings"

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
	// customer. The corner ribbon says "this space is free" instead of
	// "advertisement", because calling an unsold slot an advertisement would be
	// the one dishonest thing in an otherwise honest placeholder.
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
		{"free", "/advertise"},   // the direct offer
		{"beside", "/advertise"}, // the format: beside the text, never over it
		{"notrack", "/privacy"},  // no third-party tracking — a real difference
		{"price", "/advertise"},  // open rate card, no agency in between
		{"realty", "/advertise"}, // the audience nearest to a buying decision
	}
	out := make([]Ad, 0, len(slides))
	for _, s := range slides {
		out = append(out, Ad{
			URL:   s.cta,
			Title: T(lang, "house."+s.key+"_title"),
			Desc:  T(lang, "house."+s.key+"_desc"),
			Price: T(lang, "house."+s.key+"_cta"),
			House: true,
		})
	}
	return out
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
func (m *Module) sidebarAds(r *http.Request, lang string) []Ad {
	surface := adSurfaceFor(r)
	if surface == "" || m.ads == nil {
		return nil
	}
	orders, err := m.ads.ActiveBySurface(r.Context(), surface, lang, 12)
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
