package articles

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func posterRouter() http.Handler {
	m := &Module{}
	r := chi.NewRouter()
	r.Get("/q/{code}", m.handlePosterLink)
	return r
}

func TestPosterLinkRedirectsTagged(t *testing.T) {
	for code, want := range map[string]string{
		"rd":  "/?utm_source=qr_rudny&utm_medium=qr",
		"kst": "/?utm_source=qr_kostanay&utm_medium=qr",
		"kch": "/?utm_source=qr_kachar&utm_medium=qr",
	} {
		rec := httptest.NewRecorder()
		posterRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/q/"+code, nil))
		if rec.Code != http.StatusFound {
			t.Errorf("/q/%s: status %d, want 302 (a 301 would be cached forever and pin the destination)", code, rec.Code)
		}
		if got := rec.Header().Get("Location"); got != want {
			t.Errorf("/q/%s → %q, want %q", code, got, want)
		}
	}
}

// A printed code outlives the campaign it advertised. Whatever happens to the
// registry, the person holding the poster must land somewhere real.
func TestPosterLinkUnknownCodeGoesHome(t *testing.T) {
	for _, code := range []string{"zz", "RD-typo", "0"} {
		rec := httptest.NewRecorder()
		posterRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/q/"+code, nil))
		if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/" {
			t.Errorf("/q/%s → %d %q, want 302 /", code, rec.Code, rec.Header().Get("Location"))
		}
	}
}

func TestPosterCodesAreCaseInsensitive(t *testing.T) {
	rec := httptest.NewRecorder()
	posterRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/q/RD", nil))
	if got := rec.Header().Get("Location"); !strings.Contains(got, "qr_rudny") {
		t.Errorf("/q/RD → %q; a code printed in capitals must still resolve", got)
	}
}

// The counter table takes a closed set of labels. Campaign labels are allowed
// in because posterTargets vouches for them — nothing else may invent one.
func TestUTMSourceAcceptsOnlyRegisteredCampaigns(t *testing.T) {
	for _, v := range []string{"qr_rudny", "qr_kostanay", "qr_kachar", "QR_RUDNY", " qr ", "qr"} {
		if utmSource(v) == "" {
			t.Errorf("utmSource(%q) = \"\", want a label", v)
		}
	}
	for _, v := range []string{"qr_almaty", "qr_", "poster", "qr_rudny_extra", "'; DROP TABLE"} {
		if got := utmSource(v); got != "" {
			t.Errorf("utmSource(%q) = %q, want \"\" — unregistered labels must not reach the counters", v, got)
		}
	}
}

// Without this the admin panel falls back to printing the raw label, so a
// campaign added in a hurry shows up as "qr_kachar" beside "Прямые".
func TestEveryPosterLabelHasPanelTitles(t *testing.T) {
	labels := []string{"qr"}
	for _, tg := range posterTargets {
		labels = append(labels, tg.Label)
	}
	for _, label := range labels {
		for _, lang := range []string{"kz", "ru", "en"} {
			key := "ag.source." + label
			if got := T(lang, key); got == key || got == "" {
				t.Errorf("%s [%s] has no panel title", key, lang)
			}
		}
	}
}

func TestSafeNextRefusesOffsiteDestinations(t *testing.T) {
	for _, v := range []string{
		"//evil.example",         // protocol-relative: browsers leave the site
		"https://evil.example",   // absolute
		"http://evil.example/x",  //
		"/\\evil.example",        // some clients fold the backslash to a slash
		"/ok\r\nLocation: /evil", // header splitting
		"listings/new",           // not rooted
		"",                       //
		"   ",                    //
	} {
		if got := safeNext(v); got != "" {
			t.Errorf("safeNext(%q) = %q, want \"\"", v, got)
		}
	}
	for _, v := range []string{"/listings/new", "/studio", "/listings?geo=abc", " /listings/new "} {
		if safeNext(v) == "" {
			t.Errorf("safeNext(%q) = \"\", want the path kept", v)
		}
	}
}
