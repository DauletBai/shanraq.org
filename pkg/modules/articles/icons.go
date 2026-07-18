package articles

import "html/template"

// Shanraq's own line-icon set — drawn in-house in one consistent style so we
// don't depend on a third-party pack. Every glyph is a 24×24 viewBox, no fill,
// 1.7px currentColor stroke with round joins, sized in `em` so it follows the
// surrounding text. The map holds only the inner markup; icon() wraps it.
var iconPaths = map[string]string{
	// ---- room types (real estate) ----
	"room_living": `<path d="M3 11v-1a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v1"/>` +
		`<path d="M2 12a2 2 0 0 1 4 0v3h12v-3a2 2 0 0 1 4 0v4a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2z"/>` +
		`<path d="M5 19v2M19 19v2"/>`,
	"room_bedroom": `<path d="M2 18v-3a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v3"/>` +
		`<path d="M2 18v3M22 18v3"/><path d="M2 16h20"/>` +
		`<rect x="5" y="9.5" width="6" height="3.7" rx="1.2"/>`,
	"room_kitchen": `<path d="M4 11h16v3a5 5 0 0 1-5 5H9a5 5 0 0 1-5-5z"/>` +
		`<path d="M2 11h2M20 11h2"/>` +
		`<path d="M9 8c0-1.2 1-1.2 1-2.5M13 8c0-1.2 1-1.2 1-2.5"/>`,
	"room_bathroom": `<path d="M3 12h18v3a4 4 0 0 1-4 4H7a4 4 0 0 1-4-4z"/>` +
		`<path d="M5 12V7.5A2 2 0 0 1 9 7.5"/><path d="M9 7.2v.2"/>` +
		`<path d="M6.5 19l-1 2M17.5 19l1 2"/>`,
	"room_wc": `<path d="M6 3h8a1 1 0 0 1 1 1v3a1 1 0 0 1-1 1H6z"/>` +
		`<path d="M6 8h10v1.5a5 5 0 0 1-5 5 4 4 0 0 1-4-4z"/>` +
		`<path d="M10 14.5 9 21M13 14.5 14 21"/>`,
	"room_hallway": `<rect x="6" y="3" width="12" height="18" rx="1.2"/>` +
		`<path d="M14.5 12h.01"/>`,
	"room_balcony": `<path d="M3 8h18"/><path d="M3 12h18"/><path d="M3 20h18"/>` +
		`<path d="M6 8v12M10 8v12M14 8v12M18 8v12"/>`,
	"room_loggia": `<rect x="3" y="4" width="18" height="16" rx="1.2"/>` +
		`<path d="M3 13h18"/><path d="M8 13v7M12 13v7M16 13v7"/>`,
	"room_other": `<circle cx="5" cy="12" r="1.3" fill="currentColor" stroke="none"/>` +
		`<circle cx="12" cy="12" r="1.3" fill="currentColor" stroke="none"/>` +
		`<circle cx="19" cy="12" r="1.3" fill="currentColor" stroke="none"/>`,

	// ---- article / meta ----
	"eye":      `<path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12z"/><circle cx="12" cy="12" r="3"/>`,
	"calendar": `<rect x="3" y="5" width="18" height="16" rx="2"/><path d="M3 9h18M8 3v4M16 3v4"/>`,
	"clock":    `<circle cx="12" cy="12" r="9"/><path d="M12 7.5V12l3 2"/>`,
	"comment":  `<path d="M4 5h16a1 1 0 0 1 1 1v9a1 1 0 0 1-1 1H9l-4 4v-4H4a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1z"/>`,
	"heart":    `<path d="M12 20s-7-4.6-7-10a4 4 0 0 1 7-2.6A4 4 0 0 1 19 10c0 5.4-7 10-7 10z"/>`,
	"thumb_up": `<path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h10.3a2 2 0 0 0 2-1.7l1.4-9A2 2 0 0 0 18.7 9H14z"/>` +
		`<path d="M7 22H4a2 2 0 0 1-2-2v-9a2 2 0 0 1 2-2h3"/>`,
	"thumb_down": `<path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H6.7a2 2 0 0 0-2 1.7l-1.4 9A2 2 0 0 0 5.3 15H10z"/>` +
		`<path d="M17 2h3a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-3"/>`,

	// ---- weather (info bar) ----
	"wx_sun":   `<circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M2 12h2M20 12h2M5 5l1.4 1.4M17.6 17.6 19 19M19 5l-1.4 1.4M6.4 17.6 5 19"/>`,
	"wx_cloud": `<path d="M7 18a4 4 0 0 1 0-8 5 5 0 0 1 9.6-1.3A3.5 3.5 0 0 1 18 18z"/>`,
	"wx_rain":  `<path d="M7 15a4 4 0 0 1 0-8 5 5 0 0 1 9.6-1.3A3.5 3.5 0 0 1 18 15z"/><path d="M8 18l-1 3M12 18l-1 3M16 18l-1 3"/>`,
	"wx_snow":  `<path d="M7 15a4 4 0 0 1 0-8 5 5 0 0 1 9.6-1.3A3.5 3.5 0 0 1 18 15z"/><path d="M8 19h.01M12 20h.01M16 19h.01"/>`,
	"wx_fog":   `<path d="M7 13a4 4 0 0 1 0-8 5 5 0 0 1 9.6-1.3A3.5 3.5 0 0 1 18 13z"/><path d="M5 17h14M7 21h10"/>`,
	"wx_storm": `<path d="M7 14a4 4 0 0 1 0-8 5 5 0 0 1 9.6-1.3A3.5 3.5 0 0 1 18 14z"/><path d="M12 13l-2.2 4H12l-1.4 4"/>`,

	// ---- amenities (real estate) ----
	"am_air_conditioner": `<rect x="3" y="5" width="18" height="7" rx="1.6"/><path d="M6 9h9"/>` +
		`<path d="M7 15.5c1-.3 1-1.2 2-1.5M12 16.5c1-.3 1-1.2 2-1.5M16.5 15.5c1-.3 1-1.2 2-1.5"/>`,
	"am_pool": `<path d="M2 16c2 0 2 1.6 4 1.6s2-1.6 4-1.6 2 1.6 4 1.6 2-1.6 4-1.6"/>` +
		`<path d="M2 20c2 0 2 1.6 4 1.6s2-1.6 4-1.6 2 1.6 4 1.6 2-1.6 4-1.6"/>` +
		`<path d="M8 16V6a2 2 0 0 1 4 0"/><path d="M16 16V6"/>`,
	"am_parking":   `<rect x="4" y="4" width="16" height="16" rx="3"/><path d="M9.5 16V8H13a2.5 2.5 0 0 1 0 5H9.5"/>`,
	"am_garage":    `<path d="M3 21V9.5L12 4l9 5.5V21"/><rect x="7" y="13" width="10" height="8"/><path d="M7 16.5h10"/>`,
	"am_furniture": `<path d="M6 12V8a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v4"/><path d="M4 13a2 2 0 0 1 4 0v3h8v-3a2 2 0 0 1 4 0v3a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2z"/><path d="M6 20v1M18 20v1"/>`,
	"am_fridge":    `<rect x="6" y="3" width="12" height="18" rx="2"/><path d="M6 10h12"/><path d="M9 6v2M9 13v3"/>`,
	"am_washer":    `<rect x="4" y="3" width="16" height="18" rx="2"/><circle cx="12" cy="13" r="4"/><path d="M8 6h.01M11 6h.01"/>`,
	"am_internet": `<path d="M4.5 11.5a11 11 0 0 1 15 0"/><path d="M7.5 15a7 7 0 0 1 9 0"/>` +
		`<path d="M10.5 18.3a3 3 0 0 1 3 0"/><circle cx="12" cy="20.3" r="0.7" fill="currentColor" stroke="none"/>`,
	"am_tv":       `<rect x="3" y="5" width="18" height="12" rx="2"/><path d="M8 21h8M12 17v4"/>`,
	"am_security": `<path d="M12 3l7 3v5c0 4.2-3 7.4-7 9-4-1.6-7-4.8-7-9V6z"/><path d="M9 12l2 2 4-4"/>`,
	"am_elevator": `<rect x="4" y="3" width="16" height="18" rx="2"/><path d="M9 10l1.5-1.5L12 10"/><path d="M9 14l1.5 1.5L12 14"/><path d="M15.5 8v8"/>`,
	"am_heating":  `<path d="M4 8a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1M4 17a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1"/><path d="M6 7v11M10 7v11M14 7v11M18 7v11"/>`,
	"am_hot_water": `<path d="M12 4s5.5 5.5 5.5 9.5A5.5 5.5 0 0 1 6.5 13.5C6.5 9.5 12 4 12 4z"/>` +
		`<path d="M9.5 14a2.5 2.5 0 0 0 2.5 2.5"/>`,
	"am_plastic_windows": `<rect x="4" y="3" width="16" height="18" rx="1"/><path d="M12 3v18M4 12h16"/><path d="M14.5 11v2"/>`,
	"am_playground":      `<path d="M4 5h16"/><path d="M5 5 4 20M19 5l1 15"/><path d="M9 5l1 9M15 5l-1 9"/><path d="M9.5 14h5"/>`,
	"am_gas":             `<path d="M12 22a6 6 0 0 0 6-6c0-4-3-5-4-9-1 3-3 3-4 2-1 3-4 4-4 7a6 6 0 0 0 6 6z"/>`,
}

// icon returns the inline SVG for a named glyph, or empty if unknown.
func icon(name string) template.HTML {
	p, ok := iconPaths[name]
	if !ok {
		return ""
	}
	return template.HTML(`<svg class="ic" viewBox="0 0 24 24" width="1em" height="1em" fill="none" ` +
		`stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` +
		p + `</svg>`)
}

// roomIcon returns the icon for a room-type key (e.g. "bedroom").
func roomIcon(roomType string) template.HTML { return icon("room_" + roomType) }

// amenityIcon returns the icon for an amenity key (e.g. "parking").
func amenityIcon(key string) template.HTML { return icon("am_" + key) }
