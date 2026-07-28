package articles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
)

// Payment acquirer codes. Adding a real one is a new entry in the catalog plus
// a case in buildPaymentProvider — after that, connecting or disconnecting it is
// an admin toggle, not a code change on every switch.
const (
	PayProviderKaspi   = "kaspi"
	PayProviderIoka    = "ioka"
	PayProviderFreedom = "freedom"
	PayProviderHalyk   = "halyk"
)

// paymentProviderDef names a selectable acquirer for the admin panel.
type paymentProviderDef struct {
	Code  string
	Label string
	// Implemented is true once buildPaymentProvider can construct this adapter.
	// Until then it still lists in the panel but can't go live — the same
	// "not ready yet" signal the AI panel uses for a missing API key.
	Implemented bool
}

// paymentProviderCatalog is the ordered set of acquirers shown in the admin
// panel. Flip Implemented to true in the same change that adds the adapter.
var paymentProviderCatalog = []paymentProviderDef{
	{PayProviderKaspi, "Kaspi Pay", false},
	{PayProviderIoka, "ioka", false},
	{PayProviderFreedom, "Freedom Pay", false},
	{PayProviderHalyk, "Halyk ePay", false},
}

func isPaymentProvider(code string) bool {
	for _, p := range paymentProviderCatalog {
		if p.Code == code {
			return true
		}
	}
	return false
}

func paymentProviderImplemented(code string) bool {
	for _, p := range paymentProviderCatalog {
		if p.Code == code {
			return p.Implemented
		}
	}
	return false
}

// PaymentSettings is the runtime-switchable acquirer selection (no secrets — the
// keys live in the server config).
type PaymentSettings struct {
	Enabled  bool
	Provider string
}

// PaymentSettingsStore persists the selection in a single row and caches it,
// mirroring the AI settings pattern: the admin panel is the only writer, so the
// cache is refreshed on every write.
type PaymentSettingsStore struct {
	db  *pgxpool.Pool
	mu  sync.RWMutex
	cur PaymentSettings
}

// NewPaymentSettingsStore returns a store seeded with defaults until Load runs.
func NewPaymentSettingsStore(db *pgxpool.Pool, defaults PaymentSettings) *PaymentSettingsStore {
	return &PaymentSettingsStore{db: db, cur: defaults}
}

// Load reads the single settings row into the cache, seeding it from the
// config-derived defaults on first boot so the admin panel has a starting point.
func (s *PaymentSettingsStore) Load(ctx context.Context) (PaymentSettings, error) {
	var st PaymentSettings
	err := s.db.QueryRow(ctx,
		`SELECT enabled, provider FROM payment_settings WHERE id = 1`).
		Scan(&st.Enabled, &st.Provider)
	if errors.Is(err, pgx.ErrNoRows) {
		s.mu.RLock()
		seed := s.cur
		s.mu.RUnlock()
		if _, e := s.db.Exec(ctx,
			`INSERT INTO payment_settings (id, enabled, provider) VALUES (1, $1, $2)
			 ON CONFLICT (id) DO NOTHING`, seed.Enabled, seed.Provider); e != nil {
			return seed, fmt.Errorf("seed payment settings: %w", e)
		}
		return seed, nil
	}
	if err != nil {
		return st, fmt.Errorf("load payment settings: %w", err)
	}
	s.mu.Lock()
	s.cur = st
	s.mu.Unlock()
	return st, nil
}

// Get returns the cached settings.
func (s *PaymentSettingsStore) Get() PaymentSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cur
}

// Save validates and upserts the selection, then refreshes the cache. An enabled
// acquirer must be a known one; empty provider is only valid while payments off.
func (s *PaymentSettingsStore) Save(ctx context.Context, st PaymentSettings, by *uuid.UUID) error {
	st.Provider = strings.TrimSpace(st.Provider)
	if st.Enabled && !isPaymentProvider(st.Provider) {
		return fmt.Errorf("unknown payment provider %q", st.Provider)
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO payment_settings (id, enabled, provider, updated_at, updated_by)
		VALUES (1, $1, $2, NOW(), $3)
		ON CONFLICT (id) DO UPDATE SET
			enabled = EXCLUDED.enabled, provider = EXCLUDED.provider,
			updated_at = NOW(), updated_by = EXCLUDED.updated_by`,
		st.Enabled, st.Provider, by)
	if err != nil {
		return fmt.Errorf("save payment settings: %w", err)
	}
	s.mu.Lock()
	s.cur = st
	s.mu.Unlock()
	return nil
}

// buildPaymentProvider constructs the live provider for a selection. Secrets
// come from the server config; an off/unimplemented selection yields the safe
// no-op provider so money never moves by accident.
func (m *Module) buildPaymentProvider(st PaymentSettings) PaymentProvider {
	if st.Enabled {
		switch st.Provider {
		// One case per implemented adapter — e.g. once the Kaspi adapter lands:
		//   case PayProviderKaspi:
		//       return newKaspiProvider(m.rt.Config.Payments)
		}
	}
	return disabledProvider{}
}

// paymentProvider returns the effective provider for the current settings,
// resolved per request so an admin switch takes effect without a restart.
func (m *Module) paymentProvider() PaymentProvider {
	if m.paySettings == nil {
		return disabledProvider{}
	}
	return m.buildPaymentProvider(m.paySettings.Get())
}

// ---- admin panel ----

type paymentProviderStatus struct {
	Code        string
	Label       string
	Implemented bool
	IsActive    bool
}

// paymentsAdminView is the current acquirer configuration for the admin panel.
type paymentsAdminView struct {
	Enabled     bool
	Provider    string
	Providers   []paymentProviderStatus
	ActiveReady bool // the selected provider has a working adapter
}

func (m *Module) paymentsAdminView() paymentsAdminView {
	var st PaymentSettings
	if m.paySettings != nil {
		st = m.paySettings.Get()
	}
	v := paymentsAdminView{Enabled: st.Enabled, Provider: st.Provider}
	for _, p := range paymentProviderCatalog {
		ps := paymentProviderStatus{Code: p.Code, Label: p.Label, Implemented: p.Implemented, IsActive: p.Code == st.Provider}
		v.Providers = append(v.Providers, ps)
		if ps.IsActive {
			v.ActiveReady = ps.Implemented
		}
	}
	return v
}

// handleAdminPayments switches the active acquirer and payments on/off from the
// admin panel. Secret keys are never touched here — they live in the server
// config; this only records the choice, taking effect at once.
func (m *Module) handleAdminPayments(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if !canManageUsers(claims) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if m.paySettings == nil {
		http.Redirect(w, r, "/admin?ok=pay_bad", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	in := PaymentSettings{
		Enabled:  r.FormValue("enabled") == "on",
		Provider: strings.TrimSpace(r.FormValue("provider")),
	}
	var by *uuid.UUID
	if claims != nil {
		if id, err := uuid.Parse(claims.Subject); err == nil {
			by = &id
		}
	}
	if err := m.paySettings.Save(r.Context(), in, by); err != nil {
		m.rt.Logger.Warn("update payment settings", zap.Error(err))
		http.Redirect(w, r, "/admin?ok=pay_bad", http.StatusSeeOther)
		return
	}
	m.rt.Logger.Info("payment settings updated", zap.Bool("enabled", in.Enabled), zap.String("provider", in.Provider))
	http.Redirect(w, r, "/admin?ok=pay_set", http.StatusSeeOther)
}
