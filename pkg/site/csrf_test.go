package site

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// The guard is the site's only defence against a form on someone else's page
// posting to ours with the reader's cookie, so each way of arriving is checked.
func TestSameOrigin(t *testing.T) {
	cases := []struct {
		name    string
		origin  string
		referer string
		want    bool
	}{
		{"our own origin", "https://shanraq.org", "", true},
		{"another site", "https://evil.example", "", false},
		{"opaque origin", "null", "", false},
		{"no origin, our referer", "", "https://shanraq.org/read/x", true},
		{"no origin, foreign referer", "", "https://evil.example/x", false},
		{"neither", "", "", false},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodPost, "https://shanraq.org/subscribe", nil)
		r.Host = "shanraq.org"
		if c.origin != "" {
			r.Header.Set("Origin", c.origin)
		}
		if c.referer != "" {
			r.Header.Set("Referer", c.referer)
		}
		if got := SameOrigin(r); got != c.want {
			t.Errorf("%s: SameOrigin = %v, want %v", c.name, got, c.want)
		}
	}
}

// A safe method passes; a state change from elsewhere is refused with 403 and
// never reaches the handler.
func TestVerifyOrigin(t *testing.T) {
	var reached bool
	h := VerifyOrigin(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}))

	reached = false
	rec := httptest.NewRecorder()
	get := httptest.NewRequest(http.MethodGet, "https://shanraq.org/", nil)
	h.ServeHTTP(rec, get)
	if !reached || rec.Code != http.StatusOK {
		t.Errorf("GET must pass through: reached=%v code=%d", reached, rec.Code)
	}

	reached = false
	rec = httptest.NewRecorder()
	post := httptest.NewRequest(http.MethodPost, "https://shanraq.org/subscribe", nil)
	post.Host = "shanraq.org"
	post.Header.Set("Origin", "https://evil.example")
	h.ServeHTTP(rec, post)
	if reached || rec.Code != http.StatusForbidden {
		t.Errorf("cross-origin POST = %d reached=%v, want 403 and no handler", rec.Code, reached)
	}
}
