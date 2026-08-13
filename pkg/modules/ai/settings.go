package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Provider codes. Adding a backend is a new code here plus a client in Init.
const (
	ProviderAnthropic = "anthropic"
	ProviderOpenAI    = "openai"
	ProviderKimi      = "kimi"
)

// ProviderDef names a selectable backend for the admin panel, with the model
// identifiers to offer when an administrator switches to it.
//
// Editor/Translate may be empty. A model id is a string the provider's API
// accepts and nothing else validates, so guessing one means the switch appears
// to work and then every call fails at run time. Where the right identifiers
// are not known here, the form clears the fields instead and asks for them —
// which is honest about who knows the answer.
type ProviderDef struct {
	Code      string
	Label     string
	Editor    string
	Translate string
}

// providerCatalog is the ordered set of backends shown in the admin panel.
var providerCatalog = []ProviderDef{
	{ProviderAnthropic, "Claude (Anthropic)", "claude-sonnet-5", "claude-haiku-4-5"},
	{ProviderOpenAI, "ChatGPT (OpenAI)", "", ""},
	{ProviderKimi, "Kimi K3 (Moonshot)", "", ""},
}

// isProvider reports whether code is a known provider.
func isProvider(code string) bool {
	for _, p := range providerCatalog {
		if p.Code == code {
			return true
		}
	}
	return false
}

// AISettings is the runtime-switchable model configuration (no secrets).
type AISettings struct {
	Enabled        bool
	Provider       string
	EditorModel    string
	TranslateModel string
	MaxTokens      int
}

// Settings persists the active model configuration in a single DB row and
// caches it in memory, mirroring the articles module's service-flag pattern:
// the admin panel is the only writer, so the cache is refreshed on every write.
type Settings struct {
	db  *pgxpool.Pool
	mu  sync.RWMutex
	cur AISettings
}

// NewSettings returns a store seeded with defaults until Load runs.
func NewSettings(db *pgxpool.Pool, defaults AISettings) *Settings {
	return &Settings{db: db, cur: defaults}
}

// Load reads the single settings row into the cache. If no row exists yet, it
// seeds one from the current (config-derived) defaults so the admin panel has a
// concrete starting point. Returns the loaded settings.
func (s *Settings) Load(ctx context.Context) (AISettings, error) {
	var st AISettings
	err := s.db.QueryRow(ctx, `
		SELECT enabled, provider, editor_model, translate_model, max_tokens
		FROM ai_settings WHERE id = 1`).
		Scan(&st.Enabled, &st.Provider, &st.EditorModel, &st.TranslateModel, &st.MaxTokens)
	if err == pgx.ErrNoRows {
		s.mu.RLock()
		seed := s.cur
		s.mu.RUnlock()
		if err := s.insert(ctx, seed); err != nil {
			return seed, err
		}
		return seed, nil
	}
	if err != nil {
		return st, fmt.Errorf("load ai settings: %w", err)
	}
	s.mu.Lock()
	s.cur = st
	s.mu.Unlock()
	return st, nil
}

// Get returns the cached settings.
func (s *Settings) Get() AISettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cur
}

func (s *Settings) insert(ctx context.Context, st AISettings) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO ai_settings (id, enabled, provider, editor_model, translate_model, max_tokens)
		VALUES (1, $1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING`,
		st.Enabled, st.Provider, st.EditorModel, st.TranslateModel, st.MaxTokens)
	if err != nil {
		return fmt.Errorf("seed ai settings: %w", err)
	}
	return nil
}

// normalizeSettings validates and clamps settings: the provider must be known,
// both models non-empty, and max_tokens within a sane range. Pure, so the rules
// can be tested without a database.
func normalizeSettings(st AISettings) (AISettings, error) {
	st.Provider = strings.TrimSpace(st.Provider)
	st.EditorModel = strings.TrimSpace(st.EditorModel)
	st.TranslateModel = strings.TrimSpace(st.TranslateModel)
	if !isProvider(st.Provider) {
		return st, fmt.Errorf("unknown provider %q", st.Provider)
	}
	if st.EditorModel == "" || st.TranslateModel == "" {
		return st, fmt.Errorf("editor and translate models are required")
	}
	if st.MaxTokens < 256 {
		st.MaxTokens = 256
	} else if st.MaxTokens > 32000 {
		st.MaxTokens = 32000
	}
	return st, nil
}

// Save validates and upserts the settings, then refreshes the cache.
func (s *Settings) Save(ctx context.Context, st AISettings, by *uuid.UUID) error {
	st, err := normalizeSettings(st)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(ctx, `
		INSERT INTO ai_settings (id, enabled, provider, editor_model, translate_model, max_tokens, updated_at, updated_by)
		VALUES (1, $1, $2, $3, $4, $5, NOW(), $6)
		ON CONFLICT (id) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			provider = EXCLUDED.provider,
			editor_model = EXCLUDED.editor_model,
			translate_model = EXCLUDED.translate_model,
			max_tokens = EXCLUDED.max_tokens,
			updated_at = NOW(),
			updated_by = EXCLUDED.updated_by`,
		st.Enabled, st.Provider, st.EditorModel, st.TranslateModel, st.MaxTokens, by)
	if err != nil {
		return fmt.Errorf("save ai settings: %w", err)
	}
	s.mu.Lock()
	s.cur = st
	s.mu.Unlock()
	return nil
}
