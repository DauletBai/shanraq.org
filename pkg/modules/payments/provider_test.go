package payments

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// The default provider must fail closed: with no acquirer configured, no charge
// can be created and no webhook can settle anything. This is what guarantees no
// real money can move by accident before a provider is deliberately wired.
func TestDisabledProvider(t *testing.T) {
	var p Provider = disabled{}
	if p.Name() != "" {
		t.Errorf("disabled provider name = %q, want empty", p.Name())
	}
	if _, err := p.CreateCharge(context.Background(), Payment{Amount: 1000}); !errors.Is(err, ErrDisabled) {
		t.Errorf("CreateCharge err = %v, want ErrDisabled", err)
	}
	r := httptest.NewRequest("POST", "/pay/webhook/x", nil)
	res, err := p.HandleWebhook(r)
	if !errors.Is(err, ErrDisabled) {
		t.Errorf("HandleWebhook err = %v, want ErrDisabled", err)
	}
	if res.Paid || res.Failed {
		t.Error("disabled webhook must never report a settlement")
	}
}

func TestIsProvider(t *testing.T) {
	for _, code := range []string{ProviderKaspi, ProviderIoka, ProviderFreedom, ProviderHalyk} {
		if !IsProvider(code) {
			t.Errorf("IsProvider(%q) = false, want true", code)
		}
	}
	for _, bad := range []string{"", "stripe", "paypal", "kaspi2"} {
		if IsProvider(bad) {
			t.Errorf("IsProvider(%q) = true, want false", bad)
		}
	}
}

// The safety invariant: a provider without a working adapter must never yield a
// live provider, no matter what the admin selected — money can't move by
// accident. When an adapter lands (Implemented=true), that provider is exempt.
func TestUnimplementedProvidersStayDisabled(t *testing.T) {
	for _, p := range Catalog {
		if p.Implemented {
			continue
		}
		got := buildProvider(Settings{Enabled: true, Provider: p.Code}, nil)
		if _, ok := got.(disabled); !ok {
			t.Errorf("provider %q has no adapter but got %T, want disabled", p.Code, got)
		}
	}
}

func TestPaymentsOffYieldsDisabled(t *testing.T) {
	if _, ok := buildProvider(Settings{Enabled: false, Provider: ProviderKaspi}, nil).(disabled); !ok {
		t.Error("payments off must resolve to the disabled provider")
	}
}

// A module whose Init never ran (no settings store) must still fail closed
// rather than panic on the first webhook.
func TestProviderWithoutSettings(t *testing.T) {
	var m Module
	if _, ok := m.Provider().(disabled); !ok {
		t.Error("Provider() with no settings must be the disabled provider")
	}
}

// The callback endpoint must fail closed too: before any acquirer is connected,
// a POST to it settles nothing and says the service is not available, rather
// than accepting a body that could later be replayed.
func TestWebhookWithoutProviderIsUnavailable(t *testing.T) {
	r := chi.NewRouter()
	New().Routes(r)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("POST", "/pay/webhook/kaspi", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("webhook with no provider = %d, want 503", rec.Code)
	}
}
