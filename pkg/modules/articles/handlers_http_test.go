package articles

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// ---- public pages render ----

func TestPublicPagesRender(t *testing.T) {
	app := newTestApp(t)
	for _, path := range []string{"/", "/listings", "/studio/login", "/studio/register", "/about"} {
		w := app.do(http.MethodGet, path, nil)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, w.Code)
		}
	}
}

// ---- registration ----

func TestRegisterValidationAndSuccess(t *testing.T) {
	app := newTestApp(t)

	// A bad name re-renders the form with an error, no account created.
	w := app.do(http.MethodPost, "/studio/register", url.Values{
		"email": {"reg-bad@t.test"}, "password": {"Parol12345"},
		"first_name": {"иван"}, "last_name": {"Петров"}, "consent": {"on"},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("bad-name register = %d, want 200 (re-render)", w.Code)
	}
	if !strings.Contains(w.Body.String(), "alert--error") {
		t.Error("expected a validation error on the form")
	}

	// A valid registration creates the account and redirects.
	app.emails = append(app.emails, "reg-ok@t.test")
	w = app.do(http.MethodPost, "/studio/register", url.Values{
		"email": {"reg-ok@t.test"}, "password": {"Parol12345"},
		"first_name": {"Иван"}, "last_name": {"Петров"}, "consent": {"on"},
	})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("valid register = %d, want 303 (%s)", w.Code, w.Body.String())
	}
}

func TestRegisterRequiresConsent(t *testing.T) {
	app := newTestApp(t)
	w := app.do(http.MethodPost, "/studio/register", url.Values{
		"email": {"noconsent@t.test"}, "password": {"Parol12345"},
		"first_name": {"Иван"}, "last_name": {"Петров"}, // no consent
	})
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "alert--error") {
		t.Errorf("missing consent should re-render with error, got %d", w.Code)
	}
}

// ---- session gating ----

func TestSessionGating(t *testing.T) {
	app := newTestApp(t)
	app.createUser("member@t.test", "Parol12345")
	cookie := app.login("member@t.test", "Parol12345")

	// Authenticated: the studio dashboard renders.
	if w := app.do(http.MethodGet, "/studio", nil, withCookie(cookie)); w.Code != http.StatusOK {
		t.Errorf("authed /studio = %d, want 200", w.Code)
	}
	// Anonymous: redirected to login.
	if w := app.do(http.MethodGet, "/studio", nil); w.Code != http.StatusSeeOther {
		t.Errorf("anon /studio = %d, want 303", w.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	app := newTestApp(t)
	app.createUser("wrongpw@t.test", "Parol12345")
	w := app.do(http.MethodPost, "/studio/login", url.Values{
		"email": {"wrongpw@t.test"}, "password": {"not-the-password"},
	})
	// A failed login re-renders the form (200), does not set a session.
	if w.Code != http.StatusOK {
		t.Errorf("wrong password = %d, want 200 re-render", w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == "shanraq_session" && c.Value != "" {
			t.Error("failed login must not set a session cookie")
		}
	}
}

// ---- admin access control ----

func TestAdminAccessControl(t *testing.T) {
	app := newTestApp(t)
	app.createUser("plain@t.test", "Parol12345")
	plain := app.login("plain@t.test", "Parol12345")
	// A non-staff session cannot open the admin panel.
	if w := app.do(http.MethodGet, "/admin", nil, withCookie(plain)); w.Code != http.StatusSeeOther {
		t.Errorf("non-staff /admin = %d, want 303", w.Code)
	}

	app.createUser("boss@t.test", "Parol12345")
	app.makeStaff("boss@t.test", "admin")
	boss := app.login("boss@t.test", "Parol12345")
	if w := app.do(http.MethodGet, "/admin", nil, withCookie(boss)); w.Code != http.StatusOK {
		t.Errorf("admin /admin = %d, want 200", w.Code)
	}
}

// ---- listing submission requires a verified email ----

func TestListingCreateNeedsVerifiedEmail(t *testing.T) {
	app := newTestApp(t)
	app.createUser("seller@t.test", "Parol12345")
	cookie := app.login("seller@t.test", "Parol12345")
	w := app.do(http.MethodPost, "/listings/new", url.Values{
		"deal_type": {"sale"}, "property_type": {"apartment"},
		"price": {"1000000"}, "title": {"Тест квартира"}, "contact": {"+7 700"},
	}, withCookie(cookie))
	// Unverified email → the form re-renders with the verify-email error.
	if w.Code != http.StatusOK {
		t.Fatalf("listing create (unverified) = %d, want 200 re-render", w.Code)
	}
	if !strings.Contains(w.Body.String(), "alert--error") {
		t.Error("expected the verify-email error")
	}
}

// A lapsed session must not silently drop the filled form: publishing bounces
// to the login page with a reason the page can explain.
func TestListingCreateExpiredSessionRedirectsWithReason(t *testing.T) {
	app := newTestApp(t)
	w := app.do(http.MethodPost, "/listings/new", url.Values{
		"deal_type": {"sale"}, "property_type": {"apartment"},
		"title": {"Тест"}, "contact": {"+7 700"},
	})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("unauth listing create = %d, want 303", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/studio/login?reason=session_expired" {
		t.Errorf("redirect = %q, want /studio/login?reason=session_expired", loc)
	}
}

func TestLoginPageShowsSessionExpiredReason(t *testing.T) {
	app := newTestApp(t)
	w := app.do(http.MethodGet, "/studio/login?reason=session_expired", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("login page = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "alert--ok") {
		t.Error("expected a session-expired notice on the login page")
	}
}

// ---- payment webhook fails closed ----

func TestPaymentWebhookDisabled(t *testing.T) {
	app := newTestApp(t)
	w := app.do(http.MethodPost, "/pay/webhook/kaspi", nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("webhook with no provider = %d, want 503", w.Code)
	}
}

// ---- agent public page is verified-only ----

func TestAgentPublicNotFoundWhenUnverified(t *testing.T) {
	app := newTestApp(t)
	if w := app.do(http.MethodGet, "/agent/11111111-1111-1111-1111-111111111111", nil); w.Code != http.StatusNotFound {
		t.Errorf("unknown agent = %d, want 404", w.Code)
	}
}
