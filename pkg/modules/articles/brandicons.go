package articles

import "html/template"

// Brand marks for the analytics cards, so a row is recognisable before it is
// read. Drawn here rather than fetched: our own CSP blocks a CDN, and a logo
// pack would be a dependency for twenty small shapes.
//
// These are the brands' colours carrying their initial, not traced logos. That
// is a deliberate limit. Several of these marks genuinely are a coloured disc
// with a letter — Yandex, Facebook, Bing, Samsung, LinkedIn — and for those it
// is exact. For the rest a tracing done from memory and never looked at would
// be a worse thing to ship than a clean letter: at 16px the colour carries the
// recognition anyway. Where a logo is pure geometry the geometry is used —
// Windows' four squares, YouTube's triangle, Telegram's plane, Opera's ellipse,
// Chrome's ring.
//
// Colours are hard-coded rather than themed: a Chrome mark that changed hue in
// dark mode would stop being a Chrome mark.
var brandMarks = map[string]string{
	// Search and social.
	"google":     disc("#4285f4", "G"),
	"yandex":     disc("#fc3f1d", "Я"),
	"bing":       disc("#008373", "b"),
	"duckduckgo": disc("#de5833", "d"),
	"facebook":   disc("#1877f2", "f"),
	"linkedin":   plate("#0a66c2", "in"),
	"twitter":    disc("#000000", "X"),
	"youtube": `<rect y="2.5" width="16" height="11" rx="3.2" fill="#ff0000"/>` +
		`<path d="M6.5 5.7 11 8l-4.5 2.3z" fill="#fff"/>`,
	"telegram": `<circle cx="8" cy="8" r="8" fill="#29a9eb"/>` +
		`<path d="M12.6 4.6 11 12c-.1.5-.4.6-.9.4l-2.4-1.8-1.2 1.1c-.1.1-.2.2-.5.2l.2-2.4 4.4-4c.2-.2 0-.3-.3-.1L4.9 8.9 2.6 8.2c-.5-.2-.5-.5.1-.7l9.3-3.6c.4-.2.7.1.6.7z" fill="#fff"/>`,
	"whatsapp": disc("#25d366", "W"),
	"instagram": plate("#e1306c", "") +
		`<rect x="3.6" y="3.6" width="8.8" height="8.8" rx="2.9" fill="none" stroke="#fff" stroke-width="1.3"/>` +
		`<circle cx="8" cy="8" r="2.1" fill="none" stroke="#fff" stroke-width="1.3"/>`,

	// Browsers. Chrome keeps its ring because a ring is a shape, not a tracing.
	"chrome": `<circle cx="8" cy="8" r="8" fill="#4285f4"/>` +
		`<circle cx="8" cy="8" r="4.4" fill="#fff"/><circle cx="8" cy="8" r="3.2" fill="#4285f4"/>`,
	"safari":  disc("#1e90ff", "S"),
	"firefox": disc("#ff7139", "F"),
	"edge":    disc("#0f7abc", "e"),
	"opera":   `<circle cx="8" cy="8" r="8" fill="#ff1b2d"/><ellipse cx="8" cy="8" rx="2.5" ry="4.5" fill="#fff"/>`,
	"samsung": disc("#1428a0", "S"),

	// Operating systems. Windows is four squares and nothing else, so it is exact.
	"windows": `<rect x="1" y="1.4" width="6.3" height="6.1" fill="#0078d4"/>` +
		`<rect x="8.7" y="1" width="6.3" height="6.5" fill="#0078d4"/>` +
		`<rect x="1" y="8.9" width="6.3" height="6.1" fill="#0078d4"/>` +
		`<rect x="8.7" y="8.5" width="6.3" height="6.5" fill="#0078d4"/>`,
	"android":  disc("#3ddc84", "A"),
	"linux":    disc("#2b2b2b", "L"),
	"chromeos": disc("#4285f4", "C"),
}

// The three Apple entries share a colour but not a letter: Applebot, iOS and
// macOS turn up in different cards, and a row of identical black dots would say
// less than the words beside them.
func init() {
	brandMarks["apple"] = disc("#1d1d1f", "A") // Applebot, in the bots card
	brandMarks["ios"] = disc("#1d1d1f", "i")   // iOS, in the OS card
	brandMarks["macos"] = disc("#555558", "M") // macOS, same card as iOS
}

// plate draws a rounded square carrying a letter, for the brands whose mark is
// square rather than round.
func plate(fill, letter string) string {
	out := `<rect width="16" height="16" rx="4" fill="` + fill + `"/>`
	if letter != "" {
		out += label(letter)
	}
	return out
}

// disc draws a filled circle carrying a letter — the actual shape of several of
// these logos, and the honest fallback for the few that cannot be traced.
func disc(fill, letter string) string {
	out := `<circle cx="8" cy="8" r="8" fill="` + fill + `"/>`
	if letter != "" {
		out += label(letter)
	}
	return out
}

func label(text string) string {
	return `<text x="8" y="8" fill="#fff" font-size="9" font-weight="700" font-family="Helvetica,Arial,sans-serif"` +
		` text-anchor="middle" dominant-baseline="central">` + text + `</text>`
}

// brandIcon returns the mark for a metric row's slug, or an empty string for
// rows that name no brand — "direct", "other", "mobile" and the like. The
// template reserves the column either way, so the labels stay aligned.
func brandIcon(slug string) template.HTML {
	mark, ok := brandMarks[slug]
	if !ok {
		return ""
	}
	return template.HTML(`<svg class="brandico" viewBox="0 0 16 16" aria-hidden="true">` + mark + `</svg>`)
}
