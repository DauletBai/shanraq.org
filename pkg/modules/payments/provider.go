package payments

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"shanraq.org/pkg/shanraq"
)

// ErrDisabled is returned when no payment provider is configured — the default
// state, exactly like the AI and SMTP modules ship disabled.
var ErrDisabled = errors.New("payments are not configured")

// Charge is what a provider hands back for an amount: a QR to display and/or a
// checkout URL to open, plus the provider's own reference.
type Charge struct {
	QR  string // QR payload or data: image, for display
	URL string // hosted checkout / redirect URL
	Ref string // provider's transaction reference
}

// WebhookResult is the outcome a provider extracts from an authenticated
// callback: which of OUR payments it concerns and whether it succeeded.
type WebhookResult struct {
	PaymentID uuid.UUID
	Paid      bool
	Failed    bool
}

// Provider is the seam every acquirer/aggregator plugs into. A concrete adapter
// (ioka, Freedom Pay, Halyk ePay, Kaspi) implements it once the owner picks a
// provider and supplies sandbox keys; the order flow never changes.
type Provider interface {
	Name() string
	// CreateCharge asks the provider to open a charge for a payment we created,
	// returning a QR/URL to present to the buyer.
	CreateCharge(ctx context.Context, p Payment) (Charge, error)
	// HandleWebhook authenticates an incoming provider callback (signature,
	// source) and reports which payment it settles. It MUST reject an unsigned
	// or forged request — a payment is only ever marked paid from a verified
	// callback, never from the client.
	HandleWebhook(r *http.Request) (WebhookResult, error)
}

// disabled is the default: no money moves, nothing to misconfigure, no
// credentials in the repo. Every method fails closed.
type disabled struct{}

func (disabled) Name() string { return "" }
func (disabled) CreateCharge(context.Context, Payment) (Charge, error) {
	return Charge{}, ErrDisabled
}
func (disabled) HandleWebhook(*http.Request) (WebhookResult, error) {
	return WebhookResult{}, ErrDisabled
}

// Acquirer codes. Adding a real one is a new entry in the catalog plus a case
// in buildProvider — after that, connecting or disconnecting it is an admin
// toggle, not a code change on every switch.
const (
	ProviderKaspi   = "kaspi"
	ProviderIoka    = "ioka"
	ProviderFreedom = "freedom"
	ProviderHalyk   = "halyk"
)

// ProviderDef names a selectable acquirer for the admin panel.
type ProviderDef struct {
	Code  string
	Label string
	// Implemented is true once buildProvider can construct this adapter. Until
	// then it still lists in the panel but can't go live — the same "not ready
	// yet" signal the AI panel uses for a missing API key.
	Implemented bool
}

// Catalog is the ordered set of acquirers shown in the admin panel. Flip
// Implemented to true in the same change that adds the adapter.
var Catalog = []ProviderDef{
	{ProviderKaspi, "Kaspi Pay", false},
	{ProviderIoka, "ioka", false},
	{ProviderFreedom, "Freedom Pay", false},
	{ProviderHalyk, "Halyk ePay", false},
}

// IsProvider reports whether code names an acquirer we know.
func IsProvider(code string) bool {
	for _, p := range Catalog {
		if p.Code == code {
			return true
		}
	}
	return false
}

// ProviderImplemented reports whether an acquirer has a working adapter.
func ProviderImplemented(code string) bool {
	for _, p := range Catalog {
		if p.Code == code {
			return p.Implemented
		}
	}
	return false
}

// buildProvider constructs the live provider for a selection. Secrets come from
// the server config; an off or unimplemented selection yields the safe no-op
// provider so money never moves by accident.
func buildProvider(st Settings, rt *shanraq.Runtime) Provider {
	if st.Enabled && rt != nil {
		switch st.Provider {
		// One case per implemented adapter — e.g. once the Kaspi adapter lands:
		//   case ProviderKaspi:
		//       return newKaspi(rt.Config.Payments)
		}
	}
	return disabled{}
}
