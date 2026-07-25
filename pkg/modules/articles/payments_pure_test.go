package articles

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
)

// The default provider must fail closed: with no acquirer configured, no charge
// can be created and no webhook can settle anything. This is what guarantees no
// real money can move by accident before a provider is deliberately wired.
func TestDisabledPaymentProvider(t *testing.T) {
	var p PaymentProvider = disabledProvider{}
	if p.Name() != "" {
		t.Errorf("disabled provider name = %q, want empty", p.Name())
	}
	if _, err := p.CreateCharge(context.Background(), Payment{Amount: 1000}); !errors.Is(err, ErrPaymentsDisabled) {
		t.Errorf("CreateCharge err = %v, want ErrPaymentsDisabled", err)
	}
	r := httptest.NewRequest("POST", "/pay/webhook/x", nil)
	res, err := p.HandleWebhook(r)
	if !errors.Is(err, ErrPaymentsDisabled) {
		t.Errorf("HandleWebhook err = %v, want ErrPaymentsDisabled", err)
	}
	if res.Paid || res.Failed {
		t.Error("disabled webhook must never report a settlement")
	}
}
