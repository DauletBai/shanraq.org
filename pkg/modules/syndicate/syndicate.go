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

// Mailer sends an email. Satisfied by the notifier module.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

// Module implements the RSS route, Telegram publish job, and email digest.
type Module struct {
	rt           *shanraq.Runtime
	db           *pgxpool.Pool
	log          *zap.Logger
	http         *http.Client
	baseURL      string
	tgEnabled    bool
	tgBotToken   string
	tgChatID     string
	mailer       Mailer
	emailEnabled bool
}

// New returns a module. mailer (the notifier) powers the email digest; pass nil
// to disable email entirely.
func New(mailer Mailer) *Module { return &Module{mailer: mailer} }

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

	smtp := rt.Config.Notifications.SMTP
	m.emailEnabled = m.mailer != nil && strings.TrimSpace(smtp.Host) != "" && strings.TrimSpace(smtp.From) != ""

	if m.tgEnabled {
		m.log.Info("syndicate telegram enabled", zap.String("chat", m.tgChatID))
	} else {
		m.log.Info("syndicate telegram disabled (RSS still active at /feed.xml)")
	}
	if m.emailEnabled {
		m.log.Info("syndicate email digest enabled (weekly)")
	} else {
		m.log.Info("syndicate email digest disabled (configure SMTP to enable); subscriptions still stored")
	}
	return nil
}

// Routes exposes the RSS feed and subscription endpoints.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}
	r.Get("/feed.xml", m.handleRSS)
	r.Post("/subscribe", m.handleSubscribe)
	r.Get("/unsubscribe", m.handleUnsubscribe)
}

// Start runs the weekly digest scheduler. It checks a few times a day whether a
// week has elapsed since the last send; the send itself is a no-op without SMTP.
func (m *Module) Start(ctx context.Context, _ *shanraq.Runtime) error {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !m.emailEnabled || !m.digestDue(ctx) {
				continue
			}
			sent, err := m.SendDigest(ctx)
			if err != nil {
				m.log.Error("weekly digest", zap.Error(err))
				continue
			}
			m.markDigestSent(ctx)
			m.log.Info("weekly digest sent", zap.Int("recipients", sent))
		}
	}
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
	shanraq.StarterModule
} = (*Module)(nil)
