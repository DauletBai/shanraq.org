package articles

import (
	"encoding/json"
	"html/template"
	"strings"
)

// Ad rate card, format × surface model.
//
// Two dimensions set the price. The BANNER FORMAT is its size and prominence —
// a top billboard is worth more than a small sidebar rectangle. The SURFACE is
// where it runs — the home page and the real-estate section carry more traffic
// than a niche rubric, so they cost more. An advertiser picks one format and
// the surfaces they want; the price is the format's rate times the sum of the
// surface weights.
//
// The card is sized so that, fully sold, the monthly inventory is about
// 5 000 000 ₸ — the ceiling the platform is built to grow into, not launch
// revenue. Most slots sell only as traffic arrives; the card scales to that
// without another rewrite.

// ---- surfaces ----
const (
	surfaceHome       = "home"
	surfaceRealestate = "realestate"
	surfaceArticles   = "articles"
	surfaceRubricPfx  = "rubric:"
)

// adSurfaceWeight10 is the surface's price weight ×10 (to keep integer math).
// Home and real estate carry the most traffic; the three broad rubrics are
// mid; niche rubrics are the base.
func adSurfaceWeight10(code string) int64 {
	switch code {
	case surfaceHome, surfaceRealestate:
		return 20
	case surfaceArticles,
		surfaceRubricPfx + "society", surfaceRubricPfx + "politics", surfaceRubricPfx + "economy":
		return 13
	default:
		return 10
	}
}

// AdSurface is one bookable place, as shown in the cabinet.
type AdSurface struct {
	Code   string `json:"code"`
	Kind   string `json:"kind"`
	Cat    string `json:"cat"`
	Weight int64  `json:"weight"` // ×10
}

// AdSurfaces returns every bookable surface, fixed placements first then rubrics.
func AdSurfaces() []AdSurface {
	out := []AdSurface{
		{Code: surfaceHome, Kind: "home", Weight: adSurfaceWeight10(surfaceHome)},
		{Code: surfaceRealestate, Kind: "realestate", Weight: adSurfaceWeight10(surfaceRealestate)},
		{Code: surfaceArticles, Kind: "articles", Weight: adSurfaceWeight10(surfaceArticles)},
	}
	for _, c := range Categories {
		code := surfaceRubricPfx + c
		out = append(out, AdSurface{Code: code, Kind: "rubric", Cat: c, Weight: adSurfaceWeight10(code)})
	}
	return out
}

var adSurfaceSet = func() map[string]bool {
	m := map[string]bool{}
	for _, s := range AdSurfaces() {
		m[s.Code] = true
	}
	return m
}()

func isAdSurface(code string) bool { return adSurfaceSet[code] }

// ---- formats ----

// AdFormat is a banner size/position, the other price dimension.
type AdFormat struct {
	Code     string `json:"code"`
	Size     string `json:"size"`     // e.g. "970×250"
	Slots    int    `json:"slots"`    // rotation capacity for this format on a surface
	Price30  int64  `json:"price30"`  // base 30-day price on a ×1.0 surface
	Vertical bool   `json:"vertical"` // tall (sidebar) vs wide (top)
}

// adFormatPrice is per-format, per-duration base price on a base (×1.0) surface.
var adFormatPrice = map[string]map[int]int64{
	// wide top billboard — the most prominent, one exclusive slot.
	"horizontal": {3: 18000, 7: 36000, 14: 60000, 30: 90000},
	// tall sidebar half-page — high impact, one slot.
	"vertical": {3: 14000, 7: 28000, 14: 47000, 30: 70000},
	// square — standard, one slot.
	"square": {3: 9000, 7: 18000, 14: 30000, 30: 45000},
	// medium rectangle — the workhorse, three rotating slots.
	"rectangle": {3: 8000, 7: 16000, 14: 27000, 30: 40000},
}

// adFormatSlots is how many advertisers share a format's position on a surface.
var adFormatSlots = map[string]int{"horizontal": 1, "vertical": 1, "square": 1, "rectangle": 3}

// AdFormats returns the banner formats in display order (most prominent first).
func AdFormats() []AdFormat {
	order := []struct {
		code, size string
		vert       bool
	}{
		{"horizontal", "970×250", false},
		{"vertical", "300×600", true},
		{"square", "300×300", true},
		{"rectangle", "300×250", true},
	}
	out := make([]AdFormat, 0, len(order))
	for _, f := range order {
		out = append(out, AdFormat{
			Code: f.code, Size: f.size, Slots: adFormatSlots[f.code],
			Price30: adFormatPrice[f.code][30], Vertical: f.vert,
		})
	}
	return out
}

func isAdFormat(code string) bool { _, ok := adFormatPrice[code]; return ok }

// AdFormatSlots exposes a format's rotation capacity (templates/serving).
func AdFormatSlots(format string) int {
	if n, ok := adFormatSlots[format]; ok {
		return n
	}
	return 1
}

// ---- durations ----
var adDurationDays = []int{3, 7, 14, 30}

var adDurationSet = func() map[int]bool {
	m := make(map[int]bool, len(adDurationDays))
	for _, d := range adDurationDays {
		m[d] = true
	}
	return m
}()

// AdDurations exposes the periods to templates.
func AdDurations() []int { return adDurationDays }

// ---- pricing ----

// AdOrderPricing is the breakdown for a format over a set of surfaces.
type AdOrderPricing struct {
	Format     string `json:"format"`
	Surfaces   int    `json:"surfaces"`
	Weight10   int64  `json:"weight10"`
	FormatRate int64  `json:"format_rate"`
	Total      int64  `json:"total"`
}

// AdOrderTotal prices a format + surface set for a period. Total is the format's
// rate for the period times the summed surface weights.
func AdOrderTotal(format string, surfaces []string, days int) AdOrderPricing {
	byDur, ok := adFormatPrice[format]
	if !ok {
		byDur = adFormatPrice["rectangle"]
		format = "rectangle"
	}
	rate, ok := byDur[days]
	if !ok {
		rate = byDur[3]
	}
	var w int64
	n := 0
	for _, s := range surfaces {
		if isAdSurface(s) {
			w += adSurfaceWeight10(s)
			n++
		}
	}
	return AdOrderPricing{
		Format: format, Surfaces: n, Weight10: w,
		FormatRate: rate, Total: rate * w / 10,
	}
}

// AdSurfaceFormatPrice is the 30-day price of one format on one surface — the
// per-surface figure shown next to each checkbox, which changes with the format.
func AdSurfaceFormatPrice(format, surface string, days int) int64 {
	byDur, ok := adFormatPrice[format]
	if !ok {
		return 0
	}
	rate := byDur[days]
	if rate == 0 {
		rate = byDur[30]
	}
	return rate * adSurfaceWeight10(surface) / 10
}

// surfaceLabelKey maps a surface code to its i18n key.
func surfaceLabelKey(code string) string {
	if strings.HasPrefix(code, surfaceRubricPfx) {
		return "cat." + strings.TrimPrefix(code, surfaceRubricPfx)
	}
	return "adv.surf_" + code
}

// SurfaceLabelKey is the template helper for the above.
func SurfaceLabelKey(code string) string { return surfaceLabelKey(code) }

// AdRatesJSON hands the format prices and surface weights to the cabinet JS so
// the running total and per-surface prices update without a round trip.
func AdRatesJSON() template.JS {
	weights := map[string]int64{}
	for _, s := range AdSurfaces() {
		weights[s.Code] = s.Weight
	}
	payload := map[string]any{
		"formatPrice": adFormatPrice,
		"weights":     weights,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return template.JS("{}")
	}
	return template.JS(b)
}

// adRubricSurface builds the surface code for a category page.
func adRubricSurface(cat string) string { return surfaceRubricPfx + cat }
