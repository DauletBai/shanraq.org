// Package ai provides the AI writing assistant: a co-editor that polishes
// drafts into a respectful, professional tone, an automatic translator that
// fills the KZ/RU/EN language versions of an article, plus moderation, support,
// and listing helpers.
//
// The LLM is reached through the provider-agnostic Completer interface, so the
// backend can be swapped at runtime — the active provider (Claude, ChatGPT, or
// Kimi) and models are chosen from the admin panel and stored in the DB, while
// the secret API keys stay in the server config, one per provider. The module
// stays disabled — and spends nothing — until a key is configured and AI is
// turned on.
package ai

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// ErrDisabled is returned by synchronous helpers when AI is not configured.
var ErrDisabled = errors.New("ai assistant is disabled")

// Request is one provider-agnostic completion call.
type Request struct {
	Model     string
	System    string
	User      string
	MaxTokens int
}

// Completer is the minimal LLM contract. Production uses Claude or an
// OpenAI-compatible backend; tests use a fake.
type Completer interface {
	Complete(ctx context.Context, req Request) (string, error)
}

// Module wires the assistant into the runtime and the job queue.
type Module struct {
	rt  *shanraq.Runtime
	db  *pgxpool.Pool
	log *zap.Logger

	// providers holds a client per backend that has a key configured; keyPresent
	// records which backends are usable (for the admin panel). Both are built
	// once at Init and never mutated afterwards, so they need no lock.
	providers  map[string]Completer
	keyPresent map[string]bool

	settings *Settings

	// Active snapshot, guarded by mu. applySettings is the only writer.
	mu             sync.RWMutex
	completer      Completer
	enabled        bool
	moderateModel  string
	translateModel string
	maxTokens      int
}

// New returns an unconfigured module; Init reads config and builds the clients.
func New() *Module {
	return &Module{providers: map[string]Completer{}, keyPresent: map[string]bool{}}
}

func (m *Module) Name() string { return "ai" }

// Init builds a client for every provider that has a key, loads the runtime
// model settings (seeding them from config on first run), and activates the
// selected provider.
func (m *Module) Init(ctx context.Context, rt *shanraq.Runtime) error {
	cfg := rt.Config.AI
	m.rt = rt
	m.db = rt.DB
	m.log = rt.Logger

	// Build one client per backend that has a key. Keys come from config first,
	// then the conventional env vars; they are never persisted anywhere.
	if key := firstNonEmpty(cfg.Providers.Anthropic.APIKey, envKey("ANTHROPIC_API_KEY"), cfg.APIKey); key != "" {
		m.providers[ProviderAnthropic] = newClaudeCompleter(key)
		m.keyPresent[ProviderAnthropic] = true
	}
	if key := firstNonEmpty(cfg.Providers.OpenAI.APIKey, envKey("OPENAI_API_KEY")); key != "" {
		m.providers[ProviderOpenAI] = newOpenAICompleter(key, cfg.Providers.OpenAI.BaseURL)
		m.keyPresent[ProviderOpenAI] = true
	}
	if key := firstNonEmpty(cfg.Providers.Kimi.APIKey, envKey("KIMI_API_KEY", "MOONSHOT_API_KEY")); key != "" {
		base := cfg.Providers.Kimi.BaseURL
		if strings.TrimSpace(base) == "" {
			base = "https://api.moonshot.ai/v1"
		}
		m.providers[ProviderKimi] = newOpenAICompleter(key, base)
		m.keyPresent[ProviderKimi] = true
	}

	defaults := AISettings{
		Enabled:        cfg.Enabled,
		Provider:       orDefault(cfg.Provider, ProviderAnthropic),
		TranslateModel: orDefault(cfg.TranslateModel, "claude-haiku-4-5"),
		MaxTokens:      cfg.MaxTokens,
	}
	if defaults.MaxTokens <= 0 {
		defaults.MaxTokens = 4096
	}
	m.settings = NewSettings(rt.DB, defaults)

	st, err := m.settings.Load(ctx)
	if err != nil {
		// Non-fatal: fall back to the config defaults so the app still boots.
		m.log.Warn("load ai settings", zap.Error(err))
		st = defaults
	}
	m.applySettings(st)

	if m.Enabled() {
		m.log.Info("ai assistant enabled",
			zap.String("provider", st.Provider),
			zap.String("translate_model", st.TranslateModel))
	} else {
		m.log.Info("ai assistant disabled (enable it and set a provider key to activate)")
	}
	return nil
}

// applySettings swaps the active provider/model snapshot atomically.
func (m *Module) applySettings(st AISettings) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.moderateModel = st.ModerateModel
	m.translateModel = st.TranslateModel
	m.maxTokens = st.MaxTokens
	if m.maxTokens <= 0 {
		m.maxTokens = 4096
	}
	c := m.providers[st.Provider]
	m.completer = c
	m.enabled = st.Enabled && c != nil
}

// Enabled reports whether the assistant can serve requests right now.
func (m *Module) Enabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled && m.completer != nil
}

// ReviewCheckEnabled reports whether an article is read by the model before it
// is published. This is the gate, and it is deliberately not the same question
// as whether AI is available at all: the translation and the writing assistant
// can run without anything standing between an author and the reader.
func (m *Module) ReviewCheckEnabled() bool {
	if m == nil || !m.Enabled() {
		return false
	}
	return m.settings.Get().ReviewCheck
}

// AutoTranslateEnabled reports whether the site translates articles for their
// authors. Off since the day the arithmetic was done: the rebuilt pipeline
// works — verified batches, every paragraph accounted for, every figure
// checked — and still takes eight minutes and twenty requests to do what an
// author's own model does in one, in a language they can read back.
func (m *Module) AutoTranslateEnabled() bool {
	if m == nil || !m.Enabled() {
		return false
	}
	return m.settings.Get().AutoTranslate
}

// translateClient returns the active completer and the translation model.
func (m *Module) translateClient() (Completer, string, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.enabled {
		return nil, "", 0
	}
	return m.completer, m.translateModel, m.maxTokens
}

// moderateClient returns the active completer and the model used for the cheap
// moderation tier — comment and listing screening, and the publication rules
// check.
//
// The operator can name a model for it. Unset, it stays on the translation
// model, which is the other cheap-tier job: all of these are classifications
// answering in a line of JSON, not pieces of writing, and none of them should
// be paying for the model that drafts a column.
func (m *Module) moderateClient() (Completer, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.enabled {
		return nil, ""
	}
	model := strings.TrimSpace(m.moderateModel)
	if model == "" {
		model = m.translateModel
	}
	if strings.TrimSpace(model) == "" {
		model = "claude-haiku-4-5"
	}
	return m.completer, model
}

// RegisterJobs attaches the async translation handler to the job queue.
func (m *Module) RegisterJobs(j *jobs.Module) {
	j.Handle(JobTranslate, m.handleTranslateJob)
}

// ---- admin panel API ----

// ProviderStatus is one backend's availability for the admin panel.
type ProviderStatus struct {
	Code     string
	Label    string
	HasKey   bool
	IsActive bool
	// Models this provider should be driven with. Empty when they are not known
	// here; the form then clears the fields rather than leaving another
	// provider's identifiers in place, which is what it used to do.
	Translate string
}

// AdminView is the current AI configuration for the admin panel.
type AdminView struct {
	Enabled        bool
	ReviewCheck    bool
	AutoTranslate  bool
	Provider       string
	ModerateModel  string
	TranslateModel string
	MaxTokens      int
	Providers      []ProviderStatus
	ActiveHasKey   bool
}

// AdminView returns the current settings plus per-provider key availability.
func (m *Module) AdminView() AdminView {
	st := AISettings{Provider: ProviderAnthropic}
	if m.settings != nil {
		st = m.settings.Get()
	}
	v := AdminView{
		Enabled:        st.Enabled,
		ReviewCheck:    st.ReviewCheck,
		AutoTranslate:  st.AutoTranslate,
		Provider:       st.Provider,
		ModerateModel:  st.ModerateModel,
		TranslateModel: st.TranslateModel,
		MaxTokens:      st.MaxTokens,
	}
	for _, p := range providerCatalog {
		ps := ProviderStatus{
			Code:      p.Code,
			Label:     p.Label,
			HasKey:    m.keyPresent[p.Code],
			IsActive:  p.Code == st.Provider,
			Translate: p.Translate,
		}
		v.Providers = append(v.Providers, ps)
		if ps.IsActive {
			v.ActiveHasKey = ps.HasKey
		}
	}
	return v
}

// UpdateSettings validates and stores new model settings, then activates them.
func (m *Module) UpdateSettings(ctx context.Context, in AISettings, by *uuid.UUID) error {
	if m.settings == nil {
		return errors.New("ai settings unavailable")
	}
	if err := m.settings.Save(ctx, in, by); err != nil {
		return err
	}
	m.applySettings(m.settings.Get())
	return nil
}

// setCompleter injects a completer (used by tests) and marks the module active.
func (m *Module) setCompleter(c Completer) {
	if m.providers == nil {
		m.providers = map[string]Completer{}
	}
	m.providers["test"] = c
	m.mu.Lock()
	m.completer = c
	m.enabled = true
	if m.translateModel == "" {
		m.translateModel = "test-translate"
	}
	if m.maxTokens == 0 {
		m.maxTokens = 1024
	}
	m.mu.Unlock()
	if m.log == nil {
		m.log = zap.NewNop()
	}
}

func orDefault(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

// firstNonEmpty returns the first non-blank string, or "".
func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func envKey(names ...string) string {
	for _, n := range names {
		if v := strings.TrimSpace(os.Getenv(n)); v != "" {
			return v
		}
	}
	return ""
}

var _ interface {
	shanraq.Module
	shanraq.InitializerModule
} = (*Module)(nil)
