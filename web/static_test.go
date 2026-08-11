package web

import (
	"regexp"
	"strings"
	"testing"
)

// A deploy is only visible to a returning reader if the asset URL changes with
// the file: the assets go out with a day of Cache-Control and no ETag.
func TestAssetURLCarriesAContentHash(t *testing.T) {
	got := AssetURL("/static/css/shanraq.css")
	if !regexp.MustCompile(`^/static/css/shanraq\.css\?v=[0-9a-f]{10}$`).MatchString(got) {
		t.Fatalf("AssetURL = %q, want a 10-hex-digit ?v=", got)
	}
	// Stable within a build, or every page would bust the cache on each render.
	if again := AssetURL("/static/css/shanraq.css"); again != got {
		t.Errorf("not stable: %q then %q", got, again)
	}
	// Different files must not share a version.
	if other := AssetURL("/static/vendor/leaflet.css"); strings.TrimPrefix(other, "/static/vendor/leaflet.css?v=") ==
		strings.TrimPrefix(got, "/static/css/shanraq.css?v=") {
		t.Error("two different files hashed the same")
	}
}

// An unknown path must still render a usable link rather than break the page.
func TestAssetURLUnknownPathPassesThrough(t *testing.T) {
	if got := AssetURL("/static/nope.css"); got != "/static/nope.css" {
		t.Errorf("AssetURL(unknown) = %q", got)
	}
}
