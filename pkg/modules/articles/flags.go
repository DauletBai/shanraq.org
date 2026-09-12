package articles

import (
	"html/template"
	"strings"
)

// countryFlags holds small COLORED flag SVGs — an intentional exception to the
// monochrome red set, since a flag is meaningful only in its own colours.
var countryFlags = map[string]string{
	"Казахстан": `<rect width="24" height="16" rx="2" fill="#00AFCA"/><circle cx="13" cy="7" r="2.4" fill="#FEC50C"/>` +
		`<g stroke="#FEC50C" stroke-width=".7" stroke-linecap="round"><path d="M13 3.3v1M13 10.7v-1M8.7 7h1M17.3 7h-1M9.9 3.9l.7.7M16.1 10.1l-.7-.7M16.1 3.9l-.7.7M9.9 10.1l.7-.7"/></g>` +
		`<path d="M3 2.6v10.8" stroke="#FEC50C" stroke-width=".9"/>`,
	"Россия": `<rect width="24" height="16" rx="2" fill="#fff"/>` +
		`<path d="M0 5.33h24v5.34H0z" fill="#0039A6"/>` +
		`<path d="M0 10.67h24V14a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2z" fill="#D52B1E"/>`,
}

// countryAliases maps every name a country is stored under to the key the flag
// table uses. The cascade writes the country in the language the author was
// reading, so the same place arrives as Россия, Ресей or Russia — and a lookup
// on the Russian name alone found nothing for two of the three.
var countryAliases = map[string]string{
	"Kazakhstan": "Казахстан", "Қазақстан": "Казахстан",
	"Russia": "Россия", "Ресей": "Россия",
}

// countryMark draws the flag beside a country in the public statistics: the
// drawn one where we have it, an emoji flag otherwise.
//
// The emoji is the only way to cover ninety-odd countries without hand-drawing
// ninety-odd flags, but it renders on Windows as two letters instead of a flag.
// The two countries this audience actually looks for are drawn as SVG, so they
// keep their colours on every machine; the long tail degrades to a code, which
// is still the right answer.
func countryMark(code, title string) template.HTML {
	if svg := countryFlag(title); svg != "" {
		return svg
	}
	return template.HTML(template.HTMLEscapeString(countryFlagEmoji(code)))
}

// countryFlag returns the colored flag for a country name, or "" if unknown.
func countryFlag(country string) template.HTML {
	if canonical, ok := countryAliases[country]; ok {
		country = canonical
	}
	f, ok := countryFlags[country]
	if !ok || f == "" {
		return ""
	}
	return template.HTML(`<svg class="flag" viewBox="0 0 24 16" width="1.3em" height="0.87em" aria-hidden="true">` + f + `</svg>`)
}

// countryFlagEmoji turns a two-letter ISO country code into its flag, by the
// Unicode rule that a flag IS its country code written in regional-indicator
// letters. Derived, not looked up: every country the analytics can ever report
// gets a flag, with no table to maintain and no country silently missing one.
//
// The datacenter/VPN bucket has no country by definition, so it gets a cloud —
// it is hosting, not a place. Anything that is not exactly two ASCII letters
// gets nothing rather than a mystery glyph.
//
// Caveat worth knowing: Windows renders these as the two letters instead of a
// flag. That degrades to the country code, which is still the right answer.
func countryFlagEmoji(code string) string {
	if code == datacenterLabel {
		return "☁️"
	}
	if len(code) != 2 {
		return ""
	}
	var out []rune
	for _, c := range strings.ToUpper(code) {
		if c < 'A' || c > 'Z' {
			return ""
		}
		out = append(out, regionalIndicatorA+(c-'A'))
	}
	return string(out)
}

// regionalIndicatorA is U+1F1E6 REGIONAL INDICATOR SYMBOL LETTER A.
const regionalIndicatorA = '\U0001F1E6'
