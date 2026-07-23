package articles

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// handlePaymentWebhook receives a provider's settlement callback. It never
// trusts the request on its face: the provider adapter authenticates it
// (signature, source) and only a verified success flips the payment to paid,
// which activates the order — idempotently, so a retried callback is safe.
func (m *Module) handlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	res, err := m.payProv.HandleWebhook(r)
	if err != nil {
		if errors.Is(err, ErrPaymentsDisabled) {
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
		if _, err := m.pay.MarkPaid(r.Context(), res.PaymentID); err != nil {
			m.rt.Logger.Error("mark paid", zap.Error(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	// Providers expect a 200 to stop retrying.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
