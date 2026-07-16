package httpserver

import (
	"net/http"
	"strings"
)

// securityHeaders sets baseline hardening headers on every response. The CSP
// still permits inline styles/scripts (the templates use them); it is scoped to
// same-origin resources and blocks plugins, framing, and cross-origin form
// posts — a meaningful floor short of a nonce-based policy.
func securityHeaders(next http.Handler) http.Handler {
	const csp = "default-src 'self'; " +
		"img-src 'self' data: https:; " +
		"style-src 'self' 'unsafe-inline'; " +
		"script-src 'self' 'unsafe-inline'; " +
		"font-src 'self' data:; " +
		"connect-src 'self'; " +
		"object-src 'none'; " +
		"base-uri 'self'; " +
		"frame-ancestors 'none'; " +
		"form-action 'self'"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy", csp)
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		if isHTTPS(r) {
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
