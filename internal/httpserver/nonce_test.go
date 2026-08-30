package httpserver

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var cspNonce = regexp.MustCompile(`script-src 'self' 'nonce-([A-Za-z0-9+/]+)'`)

// The policy and the page have to agree: a script tag is admitted only if it
// carries the value this response was given, so the handler must be able to
// read exactly what went into the header.
func TestTheNonceInTheHeaderIsTheOneTheHandlerGets(t *testing.T) {
	var seen string
	h := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = NonceFromContext(r.Context())
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	m := cspNonce.FindStringSubmatch(rec.Header().Get("Content-Security-Policy"))
	if m == nil {
		t.Fatalf("no script nonce in the policy: %q", rec.Header().Get("Content-Security-Policy"))
	}
	if seen == "" {
		t.Fatal("the handler was given no nonce, so every inline script would be refused")
	}
	if seen != m[1] {
		t.Errorf("the handler got %q while the header says %q", seen, m[1])
	}
	if len(seen) < 20 {
		t.Errorf("nonce %q is too short to be unguessable", seen)
	}
}

// Reuse would defeat it: a value that repeats can be read off one page and
// pasted into an injection on the next.
func TestEveryResponseGetsItsOwnNonce(t *testing.T) {
	h := securityHeaders(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		m := cspNonce.FindStringSubmatch(rec.Header().Get("Content-Security-Policy"))
		if m == nil {
			t.Fatal("a response went out without a script nonce")
		}
		if seen[m[1]] {
			t.Fatalf("nonce %q was issued twice", m[1])
		}
		seen[m[1]] = true
	}
}

// Scripts moved off 'unsafe-inline'; styles did not, and the reason is written
// where the policy is. This pins both so neither drifts by accident.
func TestScriptsAreNonceOnlyAndStylesAreNot(t *testing.T) {
	rec := httptest.NewRecorder()
	securityHeaders(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).
		ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	csp := rec.Header().Get("Content-Security-Policy")

	script := directive(csp, "script-src")
	if strings.Contains(script, "unsafe-inline") {
		t.Errorf("script-src still allows inline scripts: %q", script)
	}
	if strings.Contains(script, "unsafe-eval") {
		t.Errorf("script-src allows eval: %q", script)
	}
	// Two hundred style attributes carry chart geometry, and a nonce does not
	// cover an attribute. Dropping this needs those moved, not a header edit.
	if !strings.Contains(directive(csp, "style-src"), "unsafe-inline") {
		t.Error("style-src lost 'unsafe-inline'; every chart would lose its geometry")
	}
	for _, want := range []string{"object-src 'none'", "frame-ancestors 'none'", "base-uri 'self'", "form-action 'self'"} {
		if !strings.Contains(csp, want) {
			t.Errorf("the policy lost %q", want)
		}
	}
}

func directive(csp, name string) string {
	for _, part := range strings.Split(csp, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, name+" ") {
			return part
		}
	}
	return ""
}
