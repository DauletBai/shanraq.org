package articles

import (
	"encoding/json"
	"html/template"
)

// Ad rate card. Placement zones map 1:1 to real inventory we actually own, so
// we never sell a slot that does not exist:
//
//	home       — sidebar on the front page
//	rubric     — sidebar on one chosen rubric feed
//	articles   — sidebar on every article page
//	realestate — sidebar of the real-estate section
//	all        — package: every zone above
//
// Prices are flat per period (no impression cap), which beats CPM for the
// buyer at our traffic level and is what regional KZ media do.
var adZoneKeys = []string{"home", "rubric", "articles", "realestate", "all"}

var adZoneSet = func() map[string]bool {
	m := make(map[string]bool, len(adZoneKeys))
	for _, z := range adZoneKeys {
		m[z] = true
	}
	return m
}()

// adDurationDays are the bookable periods (days). 14 and 30 already include a
// volume discount versus the daily rate.
var adDurationDays = []int{3, 7, 14, 30}

var adDurationSet = func() map[int]bool {
	m := make(map[int]bool, len(adDurationDays))
	for _, d := range adDurationDays {
		m[d] = true
	}
	return m
}()

// adPriceTable is zone → days → price in tenge.
var adPriceTable = map[string]map[int]int64{
	"home":       {3: 4990, 7: 9900, 14: 17900, 30: 33900},
	"realestate": {3: 3990, 7: 7900, 14: 14400, 30: 26900},
	"articles":   {3: 3490, 7: 6900, 14: 12500, 30: 23900},
	"rubric":     {3: 2490, 7: 4900, 14: 8900, 30: 16900},
	"all":        {3: 9900, 7: 19900, 14: 35900, 30: 67900},
}

const (
	// adSlotCapacity is how many advertisers share one zone's slot by rotation.
	adSlotCapacity = 3
	// adExclusivePct is the surcharge for taking a zone's whole rotation.
	adExclusivePct = 50
)

// AdZones / AdDurations expose the options to templates.
func AdZones() []string  { return adZoneKeys }
func AdDurations() []int { return adDurationDays }

// AdSlotCapacity exposes the rotation size to templates.
func AdSlotCapacity() int { return adSlotCapacity }

// AdPrice returns the price for a zone + period, plus the exclusivity
// surcharge. Unknown combinations fall back to the cheapest valid period.
func AdPrice(zone string, days int, exclusive bool) int64 {
	byDays, ok := adPriceTable[zone]
	if !ok {
		byDays = adPriceTable["articles"]
	}
	price, ok := byDays[days]
	if !ok {
		price = byDays[3]
	}
	if exclusive {
		price += price * adExclusivePct / 100
	}
	return price
}

// AdPriceOf is the template helper for the price grid (no exclusivity).
func AdPriceOf(zone string, days int) int64 { return AdPrice(zone, days, false) }

// AdPriceJSON hands the whole rate card to the cabinet's JS so the total
// updates as the buyer changes zone/period without a round trip.
func AdPriceJSON() template.JS {
	b, err := json.Marshal(adPriceTable)
	if err != nil {
		return template.JS("{}")
	}
	return template.JS(b)
}
