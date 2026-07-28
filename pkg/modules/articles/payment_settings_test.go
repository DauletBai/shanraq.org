package articles

import "testing"

func TestIsPaymentProvider(t *testing.T) {
	for _, code := range []string{PayProviderKaspi, PayProviderIoka, PayProviderFreedom, PayProviderHalyk} {
		if !isPaymentProvider(code) {
			t.Errorf("isPaymentProvider(%q) = false, want true", code)
		}
	}
	for _, bad := range []string{"", "stripe", "paypal", "kaspi2"} {
		if isPaymentProvider(bad) {
			t.Errorf("isPaymentProvider(%q) = true, want false", bad)
		}
	}
}

// The safety invariant: a provider without a working adapter must never yield a
// live provider, no matter what the admin selected — money can't move by
// accident. When an adapter lands (Implemented=true), that provider is exempt.
func TestUnimplementedProvidersStayDisabled(t *testing.T) {
	var m Module
	for _, p := range paymentProviderCatalog {
		if p.Implemented {
			continue
		}
		got := m.buildPaymentProvider(PaymentSettings{Enabled: true, Provider: p.Code})
		if _, ok := got.(disabledProvider); !ok {
			t.Errorf("provider %q has no adapter but got %T, want disabledProvider", p.Code, got)
		}
	}
}

func TestPaymentsOffYieldsDisabled(t *testing.T) {
	var m Module
	got := m.buildPaymentProvider(PaymentSettings{Enabled: false, Provider: PayProviderKaspi})
	if _, ok := got.(disabledProvider); !ok {
		t.Errorf("payments off must resolve to disabledProvider, got %T", got)
	}
}

func TestPaymentProviderNilSettings(t *testing.T) {
	var m Module // paySettings is nil
	if _, ok := m.paymentProvider().(disabledProvider); !ok {
		t.Error("paymentProvider() with nil settings must be disabledProvider")
	}
}
