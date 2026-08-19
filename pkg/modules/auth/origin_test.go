package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// The guard on every cookie-authenticated endpoint: a form on another site must
// not be able to drive a logged-in staff member's session.
func TestSameOriginOnly(t *testing.T) {
	reached := false
	h := SameOriginOnly(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}))

	cases := []struct {
		name    string
		host    string
		origin  string
		referer string
		want    int
	}{
		{"same origin", "shanraq.org", "https://shanraq.org", "", http.StatusOK},
		{"same origin with port", "localhost:8080", "http://localhost:8080", "", http.StatusOK},
		{"other site", "shanraq.org", "https://evil.example", "", http.StatusForbidden},
		{"sandboxed iframe sends null", "shanraq.org", "null", "", http.StatusForbidden},
		{"no origin, same-site referer", "shanraq.org", "", "https://shanraq.org/console", http.StatusOK},
		{"no origin, foreign referer", "shanraq.org", "", "https://evil.example/x", http.StatusForbidden},
		{"neither header", "shanraq.org", "", "", http.StatusForbidden},
		{"look-alike host", "shanraq.org", "https://shanraq.org.evil.example", "", http.StatusForbidden},
		{"origin is a prefix of the host", "shanraq.org", "https://shanraq.or", "", http.StatusForbidden},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reached = false
			r := httptest.NewRequest(http.MethodPost, "/console/jobs", nil)
			r.Host = c.host
			if c.origin != "" {
				r.Header.Set("Origin", c.origin)
			}
			if c.referer != "" {
				r.Header.Set("Referer", c.referer)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, r)
			if rec.Code != c.want {
				t.Errorf("status = %d, want %d", rec.Code, c.want)
			}
			if reached != (c.want == http.StatusOK) {
				t.Errorf("handler reached = %v, want %v", reached, c.want == http.StatusOK)
			}
		})
	}
}
