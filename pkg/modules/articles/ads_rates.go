package articles

import (
	"encoding/json"
	"html/template"
	"strings"
)

// Ad rate card, per-surface model.
//
// A "surface" is one place a banner can appear, and it is the unit of sale. An
// advertiser picks the surfaces they actually want — a sportswear seller takes
// the Sport rubric, a bookseller takes Culture — and pays per surface. This is
// what makes the inventory profitable: the same page set is sold to many
// contextual advertisers instead of one flat "everywhere" package that locks
// the whole site for a single low fee.
//
// Surfaces:
//
//	home          — front-page sidebar
//	realestate    — real-estate section sidebar
//	articles      — sidebar on every article page
//	rubric:<cat>  — sidebar on one rubric feed (e.g. rubric:sport)
//
// Each surface has adSlotCapacity rotating slots, so up to three advertisers
// can share a surface for the same dates.
const (
	surfaceHome       = "home"
	surfaceRealestate = "realestate"
	surfaceArticles   = "articles"
	surfaceRubricPfx  = "rubric:"
)

// adSurfacePrice is the price of ONE surface for a period, in tenge. Uniform
// across surfaces so the arithmetic is predictable for the buyer: one surface
// for 30 days is 4 990 ₸, and each further surface adds the same.
var adSurfacePrice = map[int]int64{3: 990, 7: 1990, 14: 3490, 30: 4990}

// adDurationDays are the bookable periods (days).
var adDurationDays = []int{3, 7, 14, 30}

var adDurationSet = func() map[int]bool {
	m := make(map[int]bool, len(adDurationDays))
	for _, d := range adDurationDays {
		m[d] = true
	}
	return m
}()

const (
	// adSlotCapacity is how many advertisers share one surface by rotation.
	adSlotCapacity = 3
	// Breadth discount: buying more surfaces earns a discount, but full
	// coverage by one advertiser stays substantial because it is still charged
	// per surface. This rewards reach without giving the site away.
	adDiscountHalfPct = 10 // from half the surfaces up
	adDiscountAllPct  = 20 // all surfaces
)

// AdSurface is one bookable place, as shown in the cabinet.
type AdSurface struct {
	Code  string `json:"code"`
	Kind  string `json:"kind"`  // home | realestate | articles | rubric
	Cat   string `json:"cat"`   // category slug for rubric surfaces, else ""
	Price int64  `json:"price"` // 30-day price, for the checkbox label
}

// AdSurfaces returns every bookable surface, in display order: the fixed
// placements first, then one per rubric.
func AdSurfaces() []AdSurface {
	out := []AdSurface{
		{Code: surfaceHome, Kind: "home", Price: adSurfacePrice[30]},
		{Code: surfaceRealestate, Kind: "realestate", Price: adSurfacePrice[30]},
		{Code: surfaceArticles, Kind: "articles", Price: adSurfacePrice[30]},
	}
	for _, c := range Categories {
		out = append(out, AdSurface{Code: surfaceRubricPfx + c, Kind: "rubric", Cat: c, Price: adSurfacePrice[30]})
	}
	return out
}

// adSurfaceCount is the number of surfaces, used for the discount thresholds.
func adSurfaceCount() int { return 3 + len(Categories) }

var adSurfaceSet = func() map[string]bool {
	m := map[string]bool{}
	for _, s := range AdSurfaces() {
		m[s.Code] = true
	}
	return m
}()

func isAdSurface(code string) bool { return adSurfaceSet[code] }

// AdDurations exposes the periods to templates.
func AdDurations() []int { return adDurationDays }

// AdSlotCapacity exposes the rotation size to templates.
func AdSlotCapacity() int { return adSlotCapacity }

// adDiscountPct returns the breadth discount for a given number of selected
// surfaces: full coverage earns the most, half or more earns some, below that
// none.
func adDiscountPct(n int) int {
	total := adSurfaceCount()
	switch {
	case n >= total:
		return adDiscountAllPct
	case n*2 >= total:
		return adDiscountHalfPct
	default:
		return 0
	}
}

// AdOrderPricing is the full breakdown for a set of surfaces over a period,
// so the buyer sees exactly what they pay for and why.
type AdOrderPricing struct {
	Surfaces    int   `json:"surfaces"`
	PerSurface  int64 `json:"per_surface"`
	Subtotal    int64 `json:"subtotal"`
	DiscountPct int   `json:"discount_pct"`
	Discount    int64 `json:"discount"`
	Total       int64 `json:"total"`
}

// AdOrderTotal prices a selection of surfaces for a period. Unknown surfaces are
// ignored, and an unknown period falls back to the shortest.
func AdOrderTotal(surfaces []string, days int) AdOrderPricing {
	per, ok := adSurfacePrice[days]
	if !ok {
		per = adSurfacePrice[3]
	}
	n := 0
	for _, s := range surfaces {
		if isAdSurface(s) {
			n++
		}
	}
	sub := per * int64(n)
	pct := adDiscountPct(n)
	disc := sub * int64(pct) / 100
	return AdOrderPricing{
		Surfaces: n, PerSurface: per, Subtotal: sub,
		DiscountPct: pct, Discount: disc, Total: sub - disc,
	}
}

// surfaceLabelKey maps a surface code to its i18n key, so the template can name
// it: fixed surfaces have their own key, rubric surfaces reuse the cat.* keys.
func surfaceLabelKey(code string) string {
	if strings.HasPrefix(code, surfaceRubricPfx) {
		return "cat." + strings.TrimPrefix(code, surfaceRubricPfx)
	}
	return "adv.surf_" + code
}

// SurfaceLabelKey is the template helper for the above.
func SurfaceLabelKey(code string) string { return surfaceLabelKey(code) }

// AdRatesJSON hands the surface price and discount rules to the cabinet JS, so
// the running total updates as the buyer ticks surfaces without a round trip.
func AdRatesJSON() template.JS {
	payload := map[string]any{
		"perSurface":   adSurfacePrice,
		"surfaceCount": adSurfaceCount(),
		"discountHalf": adDiscountHalfPct,
		"discountAll":  adDiscountAllPct,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return template.JS("{}")
	}
	return template.JS(b)
}

// adRubricSurface builds the surface code for a category page.
func adRubricSurface(cat string) string { return surfaceRubricPfx + cat }

// AdSurfacePriceOf is the template helper for the rate card: the price of one
// surface for a period.
func AdSurfacePriceOf(days int) int64 {
	if p, ok := adSurfacePrice[days]; ok {
		return p
	}
	return adSurfacePrice[3]
}
