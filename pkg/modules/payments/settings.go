package payments

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Settings is the runtime-switchable acquirer selection (no secrets — the keys
// live in the server config).
type Settings struct {
	Enabled  bool
	Provider string
}

// SettingsStore persists the selection in a single row and caches it, mirroring
// the AI settings pattern: the admin panel is the only writer, so the cache is
// refreshed on every write.
type SettingsStore struct {
	db  *pgxpool.Pool
	mu  sync.RWMutex
	cur Settings
}

// NewSettingsStore returns a store seeded with defaults until Load runs.
func NewSettingsStore(db *pgxpool.Pool, defaults Settings) *SettingsStore {
	return &SettingsStore{db: db, cur: defaults}
}

// Load reads the single settings row into the cache, seeding it from the
// config-derived defaults on first boot so the admin panel has a starting point.
func (s *SettingsStore) Load(ctx context.Context) (Settings, error) {
	var st Settings
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
func (s *SettingsStore) Get() Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cur
}

// Save validates and upserts the selection, then refreshes the cache. An enabled
// acquirer must be a known one; empty provider is only valid while payments off.
func (s *SettingsStore) Save(ctx context.Context, st Settings, by *uuid.UUID) error {
	st.Provider = strings.TrimSpace(st.Provider)
	if st.Enabled && !IsProvider(st.Provider) {
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
