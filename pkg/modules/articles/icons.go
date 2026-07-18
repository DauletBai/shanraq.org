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
