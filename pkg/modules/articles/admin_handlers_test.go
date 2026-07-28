package articles

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// adminApp returns a wired test app plus a logged-in leadership (admin) cookie,
// and resets the process-global editable state (tariffs/payments/pages tables +
// the tariff cache) after the test so cases don't leak into one another.
func adminApp(t *testing.T) (*testApp, *http.Cookie) {
	app := newTestApp(t)
	app.createUser("chief@t.test", "Parol12345")
	app.makeStaff("chief@t.test", "admin")
	cookie := app.login("chief@t.test", "Parol12345")
	t.Cleanup(func() {
		app.exec("DELETE FROM tariffs")
		app.exec("DELETE FROM payment_settings")
		app.exec("DELETE FROM content_pages")
		tariffCache.Store(nil)
	})
	return app, cookie
}

// ---- content pages ----

func TestAdminPageSaveAndServe(t *testing.T) {
	app, cookie := adminApp(t)
	form := url.Values{
		"title_kz": {"Шарттар"}, "body_kz": {"kz мәтіні"},
		"title_ru": {"Условия"}, "body_ru": {"новый русский текст условий использования"},
		"title_en": {"Terms"}, "body_en": {"new english terms body"},
	}
	w := app.do(http.MethodPost, "/admin/pages/terms", form, withCookie(cookie))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/pages/terms?ok=1" {
		t.Fatalf("save = %d loc=%q, want 303 ?ok=1", w.Code, w.Header().Get("Location"))
	}
	r := app.do(http.MethodGet, "/terms", url.Values{"lang": {"ru"}})
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), "новый русский текст") {
		t.Errorf("/terms did not serve the edited body (code %d)", r.Code)
	}
}

func TestAdminPageSaveRejectsBlankBody(t *testing.T) {
	app, cookie := adminApp(t)
	form := url.Values{
		"title_kz": {"Шарттар"}, "body_kz": {"kz"},
		"title_ru": {"Условия"}, "body_ru": {"   "}, // blank body must be rejected
		"title_en": {"Terms"}, "body_en": {"en"},
	}
	w := app.do(http.MethodPost, "/admin/pages/terms", form, withCookie(cookie))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "alert--error") {
		t.Fatalf("blank-body save = %d, want 200 re-render with error", w.Code)
	}
	// The live page must still render non-empty (built-in fallback), never blank.
	r := app.do(http.MethodGet, "/terms", url.Values{"lang": {"ru"}})
	if r.Code != http.StatusOK || strings.TrimSpace(r.Body.String()) == "" {
		t.Errorf("/terms should still render non-empty (code %d)", r.Code)
	}
}

func TestAdminPageSaveForbiddenForEditor(t *testing.T) {
	app := newTestApp(t)
	app.createUser("ed@t.test", "Parol12345")
	app.makeStaff("ed@t.test", "editor") // in the /admin group, but not leadership
	cookie := app.login("ed@t.test", "Parol12345")
	w := app.do(http.MethodPost, "/admin/pages/terms",
		url.Values{"title_ru": {"X"}, "body_ru": {"Y"}}, withCookie(cookie))
	if w.Code != http.StatusForbidden {
		t.Errorf("editor page save = %d, want 403", w.Code)
	}
}

// ---- payments ----

func TestPaymentSettingsSaveValidation(t *testing.T) {
	app, _ := adminApp(t)
	st := NewPaymentSettingsStore(app.pool, PaymentSettings{})
	ctx := context.Background()
	if err := st.Save(ctx, PaymentSettings{Enabled: true, Provider: "bogus"}, nil); err == nil {
		t.Error("enabled + unknown provider must be rejected")
	}
	if err := st.Save(ctx, PaymentSettings{Enabled: true, Provider: PayProviderKaspi}, nil); err != nil {
		t.Errorf("enabled + kaspi must save: %v", err)
	}
	if got := st.Get(); !got.Enabled || got.Provider != PayProviderKaspi {
		t.Errorf("cache not refreshed after save: %+v", got)
	}
	if err := st.Save(ctx, PaymentSettings{Enabled: false}, nil); err != nil {
		t.Errorf("disabled must save: %v", err)
	}
}

func TestAdminPaymentsSave(t *testing.T) {
	app, cookie := adminApp(t)
	// valid: enable kaspi
	w := app.do(http.MethodPost, "/admin/payments",
		url.Values{"enabled": {"on"}, "provider": {PayProviderKaspi}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin?ok=pay_set" {
		t.Fatalf("valid payments save = %d loc=%q, want 303 ?ok=pay_set", w.Code, w.Header().Get("Location"))
	}
	// invalid: enabled with an unknown provider -> rejected, payments not enabled
	w = app.do(http.MethodPost, "/admin/payments",
		url.Values{"enabled": {"on"}, "provider": {"nonsense"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin?ok=pay_bad" {
		t.Errorf("invalid payments save = %d loc=%q, want 303 ?ok=pay_bad", w.Code, w.Header().Get("Location"))
	}
}

// ---- tariffs ----

func TestTariffSaveManyPersistsAndClamps(t *testing.T) {
	app, _ := adminApp(t)
	st := NewTariffStore(app.pool)
	if err := st.Load(context.Background()); err != nil {
		t.Fatalf("load: %v", err)
	}
	err := st.SaveMany(context.Background(), map[string]int64{
		"ad.horizontal.30": 123456,      // kept
		"weight.high":      -5,          // minOne -> clamped to 1
		"banner.1":         99999999999, // over cap -> maxTariffValue
	}, nil)
	if err != nil {
		t.Fatalf("SaveMany: %v", err)
	}
	if got := tariffVal("ad.horizontal.30"); got != 123456 {
		t.Errorf("ad.horizontal.30 = %d, want 123456", got)
	}
	if got := tariffVal("weight.high"); got != 1 {
		t.Errorf("weight.high = %d, want 1 (min clamp)", got)
	}
	if got := tariffVal("banner.1"); got != maxTariffValue {
		t.Errorf("banner.1 = %d, want %d (max clamp)", got, maxTariffValue)
	}
}

func TestAdminTariffsSave(t *testing.T) {
	app, cookie := adminApp(t)
	// valid numeric change
	w := app.do(http.MethodPost, "/admin/tariffs",
		url.Values{"ad.horizontal.30": {"50000"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/tariffs?ok=1" {
		t.Fatalf("valid tariffs save = %d loc=%q, want 303 ?ok=1", w.Code, w.Header().Get("Location"))
	}
	if got := tariffVal("ad.horizontal.30"); got != 50000 {
		t.Errorf("after save ad.horizontal.30 = %d, want 50000", got)
	}
	// non-numeric -> rejected, no change, error surfaced (not a false success)
	w = app.do(http.MethodPost, "/admin/tariffs",
		url.Values{"ad.horizontal.30": {"abc"}}, withCookie(cookie))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/tariffs?err=1" {
		t.Errorf("non-numeric tariffs save = %d loc=%q, want 303 ?err=1", w.Code, w.Header().Get("Location"))
	}
	if got := tariffVal("ad.horizontal.30"); got != 50000 {
		t.Errorf("rejected save must not change value; got %d, want 50000", got)
	}
}
