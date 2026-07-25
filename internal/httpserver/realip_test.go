package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func req(remote, xff, xrealip string) *http.Request {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = remote
	if xff != "" {
		r.Header.Set("X-Forwarded-For", xff)
	}
	if xrealip != "" {
		r.Header.Set("X-Real-IP", xrealip)
	}
	return r
}

func TestRealClientIP(t *testing.T) {
	trusted := parseTrustedProxies([]string{"127.0.0.1/32", "10.0.0.0/8"})

	// Untrusted peer: forwarded headers are ignored entirely (spoof-proof).
	if got := realClientIP(req("198.51.100.50:4000", "1.2.3.4", "5.6.7.8"), trusted); got != "" {
		t.Errorf("untrusted peer should yield no forwarded IP, got %q", got)
	}
	// Trusted peer + X-Real-IP wins.
	if got := realClientIP(req("127.0.0.1:5000", "1.2.3.4", "203.0.113.5"), trusted); got != "203.0.113.5" {
		t.Errorf("X-Real-IP from trusted proxy = %q, want 203.0.113.5", got)
	}
	// Trusted peer + XFF: take the right-most NON-trusted hop, ignoring a
	// spoofed left-most client entry.
	if got := realClientIP(req("127.0.0.1:5000", "9.9.9.9, 203.0.113.7, 10.1.2.3", ""), trusted); got != "203.0.113.7" {
		t.Errorf("XFF right-most untrusted = %q, want 203.0.113.7", got)
	}
	// Trusted peer but every XFF hop is trusted → nothing to promote.
	if got := realClientIP(req("127.0.0.1:5000", "10.0.0.1, 10.0.0.2", ""), trusted); got != "" {
		t.Errorf("all-trusted XFF should yield no client IP, got %q", got)
	}
}

func TestTrustedRealIPMiddleware(t *testing.T) {
	capture := func(dst *string) http.Handler {
		return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { *dst = r.RemoteAddr })
	}

	// No trusted proxies configured → forwarded headers never rewrite RemoteAddr.
	var seen string
	trustedRealIP(nil)(capture(&seen)).ServeHTTP(httptest.NewRecorder(), req("198.51.100.9:3000", "1.2.3.4", ""))
	if seen != "198.51.100.9:3000" {
		t.Errorf("no trusted proxies: RemoteAddr must be untouched, got %q", seen)
	}

	// With a trusted proxy, the header is believed and RemoteAddr is rewritten.
	trustedRealIP([]string{"127.0.0.1/32"})(capture(&seen)).ServeHTTP(httptest.NewRecorder(), req("127.0.0.1:3000", "203.0.113.9", ""))
	if seen != "203.0.113.9" {
		t.Errorf("trusted proxy XFF should set RemoteAddr, got %q", seen)
	}
}
