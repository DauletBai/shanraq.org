package media

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSameOriginOnly pins the CSRF guard on the upload endpoints. They
// authenticate by session cookie, which makes them the browser surface: without
// this, a form on another site could spend a logged-in author's storage quota
// and land an attacker's file under our domain, to be served back to readers.
func TestSameOriginOnly(t *testing.T) {
	reached := false
	h := sameOriginOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}))

	cases := []struct {
		name    string
		origin  string
		referer string
		want    int
	}{
		{"same origin", "https://shanraq.org", "", http.StatusOK},
		{"same origin with port in host", "http://localhost:8080", "", http.StatusOK},
		{"other site", "https://evil.example", "", http.StatusForbidden},
		{"origin null (sandboxed iframe)", "null", "", http.StatusForbidden},
		{"no origin, same-site referer", "", "https://shanraq.org/studio/new", http.StatusOK},
		{"no origin, foreign referer", "", "https://evil.example/x", http.StatusForbidden},
		{"neither header", "", "", http.StatusForbidden},
		{"look-alike host", "https://shanraq.org.evil.example", "", http.StatusForbidden},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reached = false
			r := httptest.NewRequest(http.MethodPost, "/media/upload", nil)
			r.Host = "shanraq.org"
			if c.name == "same origin with port in host" {
				r.Host = "localhost:8080"
			}
			if c.origin != "" {
				r.Header.Set("Origin", c.origin)
			}
			if c.referer != "" {
				r.Header.Set("Referer", c.referer)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != c.want {
				t.Errorf("status = %d, want %d", w.Code, c.want)
			}
			if reached != (c.want == http.StatusOK) {
				t.Errorf("handler reached = %v, want %v", reached, c.want == http.StatusOK)
			}
		})
	}
}
