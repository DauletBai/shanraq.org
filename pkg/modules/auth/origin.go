package auth

import (
	"net/http"
	"net/url"
)

// SameOriginOnly rejects cross-origin requests to endpoints that authenticate
// by session cookie.
//
// A cookie-authed endpoint is a browser surface, and the browser surface is
// where CSRF lives: without this, a form on any other site could drive it with
// a logged-in staff member's session. The cookie is SameSite=Lax, so a
// cross-site POST does not carry it and the attack already fails — this is the
// second, explicit layer, the one that survives a change to the cookie policy.
//
// The media and articles modules each carry their own copy of this check,
// written before there was a shared one. Three copies of a defence is how a
// defence gets lost in a refactor; they should come here, but moving them is a
// change to working code that belongs in its own commit.
func SameOriginOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !sameOrigin(r) {
			http.Error(w, "cross-origin request blocked", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// sameOrigin prefers Origin and falls back to Referer for the older clients
// that send no Origin on same-site navigations. A request carrying neither is
// refused: this guards staff endpoints, where "probably fine" is not a policy.
func sameOrigin(r *http.Request) bool {
	if o := r.Header.Get("Origin"); o != "" && o != "null" {
		u, err := url.Parse(o)
		return err == nil && u.Host == r.Host
	}
	if ref := r.Header.Get("Referer"); ref != "" {
		u, err := url.Parse(ref)
		return err == nil && u.Host == r.Host
	}
	return false
}
