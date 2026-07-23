package articles

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service status values. A service is either fully available, in maintenance
// (visible but its paid action is disabled behind an apology), or off (hidden).
const (
	svcOn          = "on"
	svcMaintenance = "maintenance"
	svcOff         = "off"
)

// Known toggleable services. The status lives in the DB (flippable from the
// admin panel without a redeploy); this registry fixes which services exist and
// how they are labelled, so the admin UI and the enforcement points agree.
const (
	SvcAdOrders     = "ad_orders"     // paid advertising in the advertiser cabinet
	SvcListingPromo = "listing_promo" // paid promotion/feature/banner for listings

	// SvcSite is the global switch: when it is not 'on', the whole site serves a
	// maintenance page (503) to everyone but staff and the recovery routes. It
	// is not a per-feature paid service, so it is kept out of knownServices and
	// handled by its own middleware and admin control.
	SvcSite = "site"
)

// ServiceDef names a service and its i18n label key for the admin panel.
type ServiceDef struct {
	Code     string
	TitleKey string
}

// knownServices is the ordered set shown in the admin panel. Adding a service
// is one line here plus wiring its enforcement point.
var knownServices = []ServiceDef{
	{SvcAdOrders, "svc.ad_orders"},
	{SvcListingPromo, "svc.listing_promo"},
}

func isKnownService(code string) bool {
	for _, s := range knownServices {
		if s.Code == code {
			return true
		}
	}
	return false
}

// validServiceCode covers every togglable code, including the global site
// switch which is not a per-feature paid service.
func validServiceCode(code string) bool { return isKnownService(code) || code == SvcSite }

func isServiceStatus(s string) bool {
	return s == svcOn || s == svcMaintenance || s == svcOff
}

// ServiceFlag is one service's current state with its localized messages.
type ServiceFlag struct {
	Code      string
	TitleKey  string
	Status    string
	MessageKZ string
	MessageRU string
	MessageEN string
}

// Available reports whether the service's paid action may run.
func (f ServiceFlag) Available() bool { return f.Status == svcOn }

// Message returns the localized notice for the service's current state.
func (f ServiceFlag) Message(lang string) string {
	switch lang {
	case LangKZ:
		return f.MessageKZ
	case LangEN:
		return f.MessageEN
	default:
		return f.MessageRU
	}
}

// ServiceFlags caches the service switches in memory. The only writer is the
// admin panel in this process, so the cache is refreshed on every write and is
// always consistent with the DB for a single instance.
type ServiceFlags struct {
	db    *pgxpool.Pool
	mu    sync.RWMutex
	cache map[string]ServiceFlag
}

func NewServiceFlags(db *pgxpool.Pool) *ServiceFlags {
	return &ServiceFlags{db: db, cache: map[string]ServiceFlag{}}
}

// Load reads all flags into the cache. Unknown services in the DB are ignored;
// known services missing from the DB default to 'on'.
func (s *ServiceFlags) Load(ctx context.Context) error {
	rows, err := s.db.Query(ctx, `
		SELECT code, status, message_kz, message_ru, message_en FROM service_flags`)
	if err != nil {
		return fmt.Errorf("load service flags: %w", err)
	}
	defer rows.Close()

	loaded := map[string]ServiceFlag{}
	for rows.Next() {
		var f ServiceFlag
		if err := rows.Scan(&f.Code, &f.Status, &f.MessageKZ, &f.MessageRU, &f.MessageEN); err != nil {
			return err
		}
		loaded[f.Code] = f
	}
	if err := rows.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.cache = loaded
	s.mu.Unlock()
	return nil
}

// Flag returns a service's flag, defaulting to an available service if it has
// no row yet, with the registry's title key attached.
func (s *ServiceFlags) Flag(code string) ServiceFlag {
	s.mu.RLock()
	f, ok := s.cache[code]
	s.mu.RUnlock()
	if !ok {
		f = ServiceFlag{Code: code, Status: svcOn}
	}
	for _, d := range knownServices {
		if d.Code == code {
			f.TitleKey = d.TitleKey
			break
		}
	}
	if code == SvcSite {
		f.TitleKey = "svc.site"
	}
	return f
}

// SiteFlag returns the global site switch (defaults to available).
func (s *ServiceFlags) SiteFlag() ServiceFlag { return s.Flag(SvcSite) }

// SiteUp reports whether the site is serving normally (not in global
// maintenance). An unloaded/missing flag defaults to up, so a cold start or a
// DB hiccup never accidentally takes the whole site down.
func (s *ServiceFlags) SiteUp() bool { return s.Flag(SvcSite).Status == svcOn }

// Available reports whether a service's paid action may run right now.
func (s *ServiceFlags) Available(code string) bool { return s.Flag(code).Available() }

// All returns every known service in registry order, for the admin panel.
func (s *ServiceFlags) All() []ServiceFlag {
	out := make([]ServiceFlag, 0, len(knownServices))
	for _, d := range knownServices {
		out = append(out, s.Flag(d.Code))
	}
	return out
}

// Set upserts a service's status and localized messages, then refreshes the
// cache so the change takes effect immediately without a redeploy.
func (s *ServiceFlags) Set(ctx context.Context, code, status, msgKZ, msgRU, msgEN string, by *uuid.UUID) error {
	if !validServiceCode(code) || !isServiceStatus(status) {
		return fmt.Errorf("unknown service or status")
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO service_flags (code, status, message_kz, message_ru, message_en, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6)
		ON CONFLICT (code) DO UPDATE SET
			status = EXCLUDED.status,
			message_kz = EXCLUDED.message_kz,
			message_ru = EXCLUDED.message_ru,
			message_en = EXCLUDED.message_en,
			updated_at = NOW(),
			updated_by = EXCLUDED.updated_by`,
		code, status, msgKZ, msgRU, msgEN, by)
	if err != nil {
		return fmt.Errorf("set service flag: %w", err)
	}
	return s.Load(ctx)
}
