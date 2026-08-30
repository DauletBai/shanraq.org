package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// nonceKey carries the per-request script nonce to whatever renders the page.
type nonceKey struct{}

// NonceFromContext returns the script nonce for this request, or "" outside one.
// A template that writes an inline script must put it on the tag, or the browser
// will refuse to run it.
func NonceFromContext(ctx context.Context) string {
	n, _ := ctx.Value(nonceKey{}).(string)
	return n
}

// securityHeaders sets baseline hardening headers on every response.
//
// Scripts run under a per-request nonce rather than 'unsafe-inline'. That is the
// half of the policy that matters: injected markup executes through a script
// tag, and a nonce it cannot guess is what stops one. The nonce is 16 random
// bytes, minted per response and never reused.
//
// Styles keep 'unsafe-inline', and deliberately. A nonce covers <style>
// elements, of which the templates have none; it does not cover style
// attributes, of which they have two hundred -- chart geometry, bar widths,
// label positions, computed per row. Removing those means moving every one into
// a class or a custom property, and the reward is small: an injected style
// attribute can deface a page, not run code. 'unsafe-hashes' would nominally
// close it while being weaker than what it replaced, so it is not used.
func securityHeaders(next http.Handler) http.Handler {
	const cspFmt = "default-src 'self'; " +
		"img-src 'self' data: https:; " +
		"style-src 'self' 'unsafe-inline'; " +
		"script-src 'self' 'nonce-%s'; " +
		"font-src 'self' data:; " +
		"connect-src 'self'; " +
		"object-src 'none'; " +
		"base-uri 'self'; " +
		"frame-ancestors 'none'; " +
		"form-action 'self'"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw [16]byte
		if _, err := rand.Read(raw[:]); err != nil {
			// Without randomness a nonce is a password everyone knows, so the
			// request is refused rather than served under a guessable policy.
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		nonce := base64.RawStdEncoding.EncodeToString(raw[:])
		csp := fmt.Sprintf(cspFmt, nonce)
		r = r.WithContext(context.WithValue(r.Context(), nonceKey{}, nonce))
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
