package articles

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"shanraq.org/pkg/modules/auth"
)

// TestExcluded verifies which requests are left out of analytics: the opt-out
// cookie, staff roles, and configured emails — while ordinary guests and
// non-listed users are counted.
func TestExcluded(t *testing.T) {
	m := &Module{excludeEmails: map[string]bool{"test@shanraq.org": true}}

	withClaims := func(r *http.Request, email string, roles ...string) *http.Request {
		c := &auth.Claims{Email: email, Roles: roles}
		return r.WithContext(auth.ContextWithClaims(r.Context(), c))
	}

	t.Run("plain guest is counted", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if m.excluded(httptest.NewRecorder(), r) {
			t.Fatal("guest must not be excluded")
		}
	})

	t.Run("opt-out cookie excludes", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: analyticsOptOutCookie, Value: "1"})
		if !m.excluded(nil, r) {
			t.Fatal("cookie must exclude")
		}
	})

	t.Run("admin excluded and cookie stamped", func(t *testing.T) {
		r := withClaims(httptest.NewRequest(http.MethodGet, "/", nil), "boss@shanraq.org", "admin")
		w := httptest.NewRecorder()
		if !m.excluded(w, r) {
			t.Fatal("admin must be excluded")
		}
		if got := w.Result().Cookies(); len(got) == 0 || got[0].Name != analyticsOptOutCookie || got[0].Value != "1" {
			t.Fatalf("admin should get the opt-out cookie stamped, got %v", got)
		}
	})

	t.Run("configured email excluded (case-insensitive)", func(t *testing.T) {
		r := withClaims(httptest.NewRequest(http.MethodGet, "/", nil), "TEST@Shanraq.org", "user")
		if !m.excluded(nil, r) {
			t.Fatal("listed email must be excluded")
		}
	})

	t.Run("ordinary signed-in user is counted", func(t *testing.T) {
		r := withClaims(httptest.NewRequest(http.MethodGet, "/", nil), "reader@example.com", "user")
		if m.excluded(httptest.NewRecorder(), r) {
			t.Fatal("ordinary registered reader must be counted")
		}
	})
}
