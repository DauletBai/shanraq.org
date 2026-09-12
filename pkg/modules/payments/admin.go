package payments

import (
	"context"

	"github.com/google/uuid"
)

// ProviderStatus is one acquirer as the admin panel shows it.
type ProviderStatus struct {
	Code        string
	Label       string
	Implemented bool
	IsActive    bool
}

// AdminView is the current acquirer configuration for the admin panel. The
// panel itself lives with the rest of the admin pages; this is the part only
// the payments module can answer.
type AdminView struct {
	Enabled     bool
	Provider    string
	Providers   []ProviderStatus
	ActiveReady bool // the selected provider has a working adapter
}

// Admin returns the acquirer configuration for the panel.
func (m *Module) Admin() AdminView {
	var st Settings
	if m.settings != nil {
		st = m.settings.Get()
	}
	v := AdminView{Enabled: st.Enabled, Provider: st.Provider}
	for _, p := range Catalog {
		ps := ProviderStatus{Code: p.Code, Label: p.Label, Implemented: p.Implemented, IsActive: p.Code == st.Provider}
		v.Providers = append(v.Providers, ps)
		if ps.IsActive {
			v.ActiveReady = ps.Implemented
		}
	}
	return v
}

// SaveAdmin records the acquirer choice made in the admin panel. Secret keys
// are never touched here — they live in the server config; this only records
// the choice, taking effect at once.
func (m *Module) SaveAdmin(ctx context.Context, st Settings, by *uuid.UUID) error {
	if m.settings == nil {
		return ErrDisabled
	}
	return m.settings.Save(ctx, st, by)
}
