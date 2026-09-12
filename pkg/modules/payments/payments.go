// Package payments is the money seam of the site: it records an intent to
// charge, settles it only from a provider callback it has authenticated, and
// releases a hold that was never paid.
//
// It knows nothing about what is being bought. A module that sells something --
// an advertising slot today, a book tomorrow -- registers a Settlement for its
// own kind of payment, and payments calls it inside the same transaction that
// flips the payment to paid. That is why an order can never be active without a
// payment, or paid without its order.
package payments

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"shanraq.org/pkg/shanraq"
)

// Module wires the payment store, the runtime acquirer selection and the
// provider callback endpoint.
type Module struct {
	rt       *shanraq.Runtime
	store    *Store
	settings *SettingsStore
}

// New returns the module. Payments ship disabled: no provider is configured,
// and every provider call fails closed until an admin selects one that has an
// adapter.
func New() *Module { return &Module{} }

func (m *Module) Name() string { return "payments" }

// Init opens the store and loads the acquirer selection.
func (m *Module) Init(ctx context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	// Which acquirer is live (and whether payments are on) is a runtime choice
	// in the admin panel — secrets stay in config. Defaults come from config so
	// a fresh database has a starting point.
	m.settings = NewSettingsStore(rt.DB, Settings{Provider: rt.Config.Payments.Provider})
	if _, err := m.settings.Load(ctx); err != nil {
		rt.Logger.Warn("payment settings unavailable, using defaults", zap.Error(err))
	}
	return nil
}

// Routes registers the provider callback. It is server-to-server and carries no
// Origin header, so it must stay outside any CSRF-guarded browser group.
func (m *Module) Routes(r chi.Router) {
	r.Post("/pay/webhook/{provider}", m.handleWebhook)
}

// expiryInterval is how often abandoned holds are released. A slot held by a
// checkout nobody finished is inventory nobody can buy, so the sweep runs often
// enough that a buyer who comes back in an hour finds the slot free.
const expiryInterval = 5 * time.Minute

// Start releases holds that were never paid.
func (m *Module) Start(ctx context.Context, rt *shanraq.Runtime) error {
	t := time.NewTicker(expiryInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			n, err := m.store.ExpirePending(ctx)
			if err != nil {
				rt.Logger.Warn("expire pending payments", zap.Error(err))
				continue
			}
			if n > 0 {
				rt.Logger.Info("released unpaid holds", zap.Int("count", n))
			}
		}
	}
}

// Store exposes the payment store to the modules that sell something.
func (m *Module) Store() *Store { return m.store }

// Handle registers what happens to a seller's own record when a payment of this
// kind settles or expires. See Settlement.
func (m *Module) Handle(kind string, s Settlement) { m.store.Handle(kind, s) }

// Provider returns the effective provider for the current settings, resolved
// per call so an admin switch takes effect without a restart.
func (m *Module) Provider() Provider {
	if m.settings == nil {
		return disabled{}
	}
	return buildProvider(m.settings.Get(), m.rt)
}

// handleWebhook receives a provider's settlement callback. It never trusts the
// request on its face: the provider adapter authenticates it (signature,
// source) and only a verified success flips the payment to paid, which
// activates the order — idempotently, so a retried callback is safe.
func (m *Module) handleWebhook(w http.ResponseWriter, r *http.Request) {
	res, err := m.Provider().HandleWebhook(r)
	if err != nil {
		if errors.Is(err, ErrDisabled) {
			// No provider configured — nothing legitimately calls this yet.
			http.Error(w, "payments disabled", http.StatusServiceUnavailable)
			return
		}
		// A rejected/forged callback is a client error, not ours.
		m.rt.Logger.Warn("payment webhook rejected", zap.Error(err))
		http.Error(w, "invalid webhook", http.StatusBadRequest)
		return
	}
	if res.Paid {
		if _, err := m.store.MarkPaid(r.Context(), res.PaymentID); err != nil {
			m.rt.Logger.Error("mark paid", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	// Providers expect a 200 to stop retrying.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
	shanraq.StarterModule
} = (*Module)(nil)
