package articles

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Short links for printed advertising.
//
// A QR code on paper is the one link we can never edit again: the sheet hangs
// in a bus shelter for months while the site keeps moving underneath it. So the
// printed code carries nothing but a two-letter key, and the destination lives
// here — a campaign can be repointed, or retired, without reprinting anything.
//
// Short also means sparse, and sparse is what gets a code read. Encoding the
// full tagged destination ("/listings/new?utm_source=qr_rudny&utm_medium=qr")
// needs roughly three times the modules of "/q/rd", and module size is exactly
// what decides whether a phone can resolve the code through scratched shelter
// glass from two metres away, in winter, at dusk. The redirect buys the
// tracking and the legibility at once, where encoding the tag directly trades
// one against the other.
type posterTarget struct {
	Path  string // where the scan lands
	Label string // analytics source label, and the utm_source it arrives under
}

// posterTargets is the whole campaign registry: printed key → destination.
//
// All three land on the front page. The owner's call, and the catalogue backs
// it up: at the time these were cut the site held five listings in the whole
// country and none at all in Rudny, so a code pointing straight at property
// search would have answered a poster with an empty result page. The front
// page has 108 articles on it and a menu — something to read for the visitor
// who was only curious, and a visible way through for the one who came to
// sell.
//
// The keys stay separate anyway. Pointing three codes at one page costs
// nothing and buys the only number that matters afterwards: which town
// answered. Repointing any of them later is a one-line edit here, with no
// reprinting.
var posterTargets = map[string]posterTarget{
	"rd":  {Path: "/", Label: "qr_rudny"},
	"kst": {Path: "/", Label: "qr_kostanay"},
	"kch": {Path: "/", Label: "qr_kachar"},
}

// handlePosterLink resolves a printed key and sends the scan on, tagged.
func (m *Module) handlePosterLink(w http.ResponseWriter, r *http.Request) {
	code := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "code")))
	t, ok := posterTargets[code]
	if !ok {
		// A printed link never earns a 404. A retired campaign, a smudged
		// digit or a misread still belongs to a real person holding a phone at
		// a bus stop, and the front page is a better answer than an error.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// 302, never 301. A permanent redirect is cached by the browser for good,
	// which would nail the destination down in exactly the place we went to
	// this trouble to keep loose.
	http.Redirect(w, r, t.Path+"?utm_source="+url.QueryEscape(t.Label)+"&utm_medium=qr",
		http.StatusFound)
}

// posterSource maps a utm_source back to a campaign label, but only if one of
// our own printed codes could have produced it.
//
// This is what lets utmSource stay a closed set while still growing: a new
// campaign adds a row to posterTargets and nothing else anywhere can write a
// fresh label into the counter table. An arbitrary ?utm_source=qr_whatever is
// still ignored.
func posterSource(v string) string {
	for _, t := range posterTargets {
		if t.Label == v {
			return t.Label
		}
	}
	return ""
}
