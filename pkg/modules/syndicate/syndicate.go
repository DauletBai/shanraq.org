// Package syndicate pushes published articles out to external, harder-to-block
// channels so the content survives even if the main domain is blocked:
//   - an always-on RSS feed at /feed.xml (read by aggregators and mirrors), and
//   - optional Telegram auto-posting on publish (activated by config).
//
// This is the resilience layer of the platform. It reads article data with raw
// SQL rather than importing the articles package, keeping the dependency graph
// acyclic (articles -> syndicate -> jobs).
package syndicate

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/shanraq"
)

// Module implements the RSS route and the Telegram publish job.
type Module struct {
	rt         *shanraq.Runtime
	db         *pgxpool.Pool
	log        *zap.Logger
	http       *http.Client
	baseURL    string
	tgEnabled  bool
	tgBotToken string
	tgChatID   string
}

// New returns an unconfigured module.
func New() *Module { return &Module{} }

func (m *Module) Name() string { return "syndicate" }

// Init reads config and prepares the HTTP client for Telegram.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	cfg := rt.Config.Syndicate
	m.rt = rt
	m.db = rt.DB
	m.log = rt.Logger
	m.http = &http.Client{Timeout: 10 * time.Second}
	m.baseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if m.baseURL == "" {
		m.baseURL = "http://localhost:8080"
	}
	m.tgBotToken = strings.TrimSpace(cfg.Telegram.BotToken)
	m.tgChatID = strings.TrimSpace(cfg.Telegram.ChatID)
	m.tgEnabled = cfg.Telegram.Enabled && m.tgBotToken != "" && m.tgChatID != ""

	if m.tgEnabled {
		m.log.Info("syndicate telegram enabled", zap.String("chat", m.tgChatID))
	} else {
		m.log.Info("syndicate telegram disabled (RSS still active at /feed.xml)")
	}
	return nil
}

// Routes exposes the RSS feed.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}
	r.Get("/feed.xml", m.handleRSS)
}

// TelegramEnabled reports whether Telegram auto-posting is configured.
func (m *Module) TelegramEnabled() bool { return m.tgEnabled }

// RegisterJobs attaches the Telegram publish handler to the queue.
func (m *Module) RegisterJobs(j *jobs.Module) {
	j.Handle(JobTelegram, m.handleTelegramJob)
}

// EnqueuePublish schedules a Telegram announcement for a newly published
// article. It is a no-op when Telegram is not configured.
func (m *Module) EnqueuePublish(ctx context.Context, store *jobs.Store, articleID uuid.UUID) error {
	if !m.tgEnabled || store == nil {
		return nil
	}
	payload, err := TelegramPayload(articleID)
	if err != nil {
		return err
	}
	return store.Enqueue(ctx, jobs.Job{
		ID:          uuid.New(),
		Name:        JobTelegram,
		Payload:     payload,
		RunAt:       time.Now(),
		MaxAttempts: 3,
	})
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
