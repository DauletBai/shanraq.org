package articles

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// Ad is one creative in the sidebar ad carousel. Demo ads are self-made,
// clearly-labeled placeholders — no real brand, logo, or photo — so nothing
// implies a commercial relationship that does not exist. Real ads would later
// come from config or a paid-placement store; the template contract stays the
// same.
type Ad struct {
	Image string // /static/... illustration
	Title string
	Price string
	Desc  string
	URL   string // click target ("#" for demo)
}

// demoAds returns three localized placeholder car ads for the sidebar slot,
// used until real paid placements are wired up.
func demoAds(lang string) []Ad {
	return []Ad{
		{Image: "/static/demo/ads/car-sedan.svg", URL: "#",
			Title: T(lang, "ad.sedan_title"), Desc: T(lang, "ad.sedan_desc"), Price: "11 990 000 ₸"},
		{Image: "/static/demo/ads/car-suv.svg", URL: "#",
			Title: T(lang, "ad.suv_title"), Desc: T(lang, "ad.suv_desc"), Price: "15 490 000 ₸"},
		{Image: "/static/demo/ads/car-hatch.svg", URL: "#",
			Title: T(lang, "ad.hatch_title"), Desc: T(lang, "ad.hatch_desc"), Price: "9 890 000 ₸"},
	}
}

// adZoneFor maps a request to the inventory zone (and rubric) it belongs to,
// mirroring the zones sold in the advertiser cabinet.
func adZoneFor(r *http.Request) (zone, rubric string) {
	p := r.URL.Path
	switch {
	case strings.HasPrefix(p, "/listings"):
		return "realestate", ""
	case strings.HasPrefix(p, "/read/"):
		return "articles", ""
	case p == "/":
		if cat := r.URL.Query().Get("cat"); cat != "" && IsCategory(cat) {
			return "rubric", cat
		}
		return "home", ""
	}
	return "", ""
}

// sidebarAds serves the paid placements booked for this page's zone. Demo
// creatives only fill the slot while nothing is sold for it.
func (m *Module) sidebarAds(r *http.Request, lang string) []Ad {
	// In production a slot with nothing sold shows nothing, not placeholder
	// cars: a demo advert on a live page reads as a real listing that failed
	// to load, or worse as a commercial relationship that does not exist.
	demoOK := !strings.EqualFold(m.rt.Config.Environment, "production")
	fallback := func() []Ad {
		if demoOK {
			return demoAds(lang)
		}
		return nil
	}
	zone, rubric := adZoneFor(r)
	if zone == "" || m.ads == nil {
		return fallback()
	}
	orders, err := m.ads.ActiveByZone(r.Context(), zone, rubric, lang, adSlotCapacity)
	if err != nil {
		m.rt.Logger.Warn("sidebar ads", zap.Error(err))
		return fallback()
	}
	if len(orders) == 0 {
		return fallback()
	}
	out := make([]Ad, 0, len(orders))
	for _, o := range orders {
		url := o.TargetURL
		if url == "" {
			url = "#"
		}
		out = append(out, Ad{Image: o.ImageURL, Title: o.Title, Desc: o.Body, Price: o.CTA, URL: url})
	}
	return out
}
