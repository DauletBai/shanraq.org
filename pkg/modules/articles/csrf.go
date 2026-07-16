package articles

import (
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

// verifyOrigin is a defense-in-depth CSRF guard for the browser (cookie-auth)
// surface. The session cookie is already SameSite=Lax — which stops the cookie
// from riding along on cross-site POSTs — so this is a second, explicit layer:
// every state-changing request must carry an Origin (or, failing that, a
// Referer) whose host matches the request host. Safe methods pass through.
//
// It is scoped to this module's router, which serves only browser forms; the
// token-authenticated API lives in other modules and is unaffected.
func (m *Module) verifyOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			next.ServeHTTP(w, r)
			return
		}
		if m.sameOrigin(r) {
			next.ServeHTTP(w, r)
			return
		}
		m.rt.Logger.Warn("blocked cross-origin state change",
			zap.String("path", r.URL.Path),
			zap.String("origin", r.Header.Get("Origin")),
			zap.String("referer", r.Header.Get("Referer")))
		http.Error(w, "cross-origin request blocked", http.StatusForbidden)
	})
}

// sameOrigin reports whether the request's Origin/Referer host matches the host
// the request was addressed to. A reverse proxy must forward the original Host
// header (the common default) for this to hold in production.
func (m *Module) sameOrigin(r *http.Request) bool {
	host := r.Host
	if o := r.Header.Get("Origin"); o != "" && o != "null" {
		if u, err := url.Parse(o); err == nil {
			return u.Host == host
		}
		return false
	}
	// No Origin (older clients / some same-origin navigations): fall back to Referer.
	if ref := r.Header.Get("Referer"); ref != "" {
		if u, err := url.Parse(ref); err == nil {
			return u.Host == host
		}
	}
	return false
}
