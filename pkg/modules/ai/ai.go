// Package ai provides a Claude-backed writing assistant: a co-editor that
// polishes drafts into a respectful, professional tone, and an automatic
// translator that fills the KZ/RU/EN language versions of an article.
//
// The LLM is reached through the provider-agnostic Completer interface, so the
// backend (currently Anthropic Claude) can be swapped without touching the
// business logic. The module stays disabled — and spends nothing — until an
// API key is configured.
package ai

import (
	"context"
	"errors"
	"os"
	"strings"

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

// Completer is the minimal LLM contract. Production uses Claude; tests use a fake.
type Completer interface {
	Complete(ctx context.Context, req Request) (string, error)
}

// Module wires the assistant into the runtime and the job queue.
type Module struct {
	rt             *shanraq.Runtime
	db             *pgxpool.Pool
	log            *zap.Logger
	completer      Completer
	enabled        bool
	editorModel    string
	translateModel string
	maxTokens      int
}

// New returns an unconfigured module; Init reads config and (if a key is
// present) constructs the Claude client.
func New() *Module { return &Module{} }

func (m *Module) Name() string { return "ai" }

// Init reads AI config and, when enabled with an API key, builds the client.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	cfg := rt.Config.AI
	m.rt = rt
	m.db = rt.DB
	m.log = rt.Logger
	m.editorModel = orDefault(cfg.EditorModel, "claude-sonnet-5")
	m.translateModel = orDefault(cfg.TranslateModel, "claude-haiku-4-5")
	m.maxTokens = cfg.MaxTokens
	if m.maxTokens <= 0 {
		m.maxTokens = 4096
	}

	key := strings.TrimSpace(cfg.APIKey)
	if key == "" {
		key = strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
	}
	if cfg.Enabled && key != "" {
		m.completer = newClaudeCompleter(key)
		m.enabled = true
		m.log.Info("ai assistant enabled", zap.String("editor_model", m.editorModel), zap.String("translate_model", m.translateModel))
	} else {
		m.log.Info("ai assistant disabled (set ai.enabled and an API key to activate)")
	}
	return nil
}

// Enabled reports whether the assistant can serve requests.
func (m *Module) Enabled() bool { return m.enabled && m.completer != nil }

// RegisterJobs attaches the async translation handler to the job queue.
func (m *Module) RegisterJobs(j *jobs.Module) {
	j.Handle(JobTranslate, m.handleTranslateJob)
}

// setCompleter injects a completer (used by tests) and marks the module active.
func (m *Module) setCompleter(c Completer) {
	m.completer = c
	m.enabled = true
	if m.editorModel == "" {
		m.editorModel = "test-editor"
	}
	if m.translateModel == "" {
		m.translateModel = "test-translate"
	}
	if m.maxTokens == 0 {
		m.maxTokens = 1024
	}
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

var _ interface {
	shanraq.Module
	shanraq.InitializerModule
} = (*Module)(nil)
