package main

import (
	// The timezone database, embedded: the runtime image carries no system
	// zoneinfo, and without this the site's own clock silently falls back to UTC.
	_ "time/tzdata"

	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/ai"
	"shanraq.org/pkg/modules/apikeys"
	"shanraq.org/pkg/modules/articles"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/health"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/modules/media"
	"shanraq.org/pkg/modules/migrations"
	"shanraq.org/pkg/modules/notifier"
	"shanraq.org/pkg/modules/sms"
	"shanraq.org/pkg/modules/syndicate"
	"shanraq.org/pkg/modules/telemetry"
	"shanraq.org/pkg/modules/webui"
	"shanraq.org/pkg/shanraq"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to YAML/JSON/TOML configuration file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}
	if err := enforceProductionSecrets(cfg); err != nil {
		panic(err)
	}

	const (
		jobWorkers     = 4
		jobPollSeconds = 2 * time.Second
	)

	tenantResolver := jobs.AuthTenantResolver()

	notifierModule := notifier.New()
	authOpts := []auth.Option{auth.WithMailer(notifierModule)}
	if cfg.Auth.MFA.TOTP.Enabled {
		authOpts = append(authOpts, auth.WithTOTP(cfg.Auth.MFA.TOTP.Issuer))
	}
	// Wire the SMS gateway used for phone verification. With no provider set the
	// client is nil and codes are dev-logged (never sent); a named-but-misconfigured
	// provider fails fast so a broken production deploy is caught at boot.
	if smsClient, on, err := sms.New(cfg.SMS); err != nil {
		panic(fmt.Errorf("configure sms: %w", err))
	} else if on {
		authOpts = append(authOpts, auth.WithSMSSender(smsClient))
	}
	// The JSON signup endpoint must obey the same registration switch as the
	// browser form. The articles module owns the service flags and is built
	// further down, so the gate is a closure resolved at request time rather
	// than at construction: by the time anyone can POST, both modules exist.
	var articlesModule *articles.Module
	authOpts = append(authOpts, auth.WithSignupGate(func() error {
		if articlesModule == nil {
			return nil
		}
		return articlesModule.RegistrationGate()
	}))
	authModule := auth.New(authOpts...)
	apiKeyModule := apikeys.New(
		apikeys.WithHTTPMiddleware(authModule.RequireRoles("user", "operator", "admin")),
	)

	jobModule := jobs.New(
		jobs.WithWorkerCount(jobWorkers),
		jobs.WithPollInterval(jobPollSeconds),
		jobs.WithTenantResolver(tenantResolver),
		// Staff only. Enqueue takes an arbitrary job name and payload, and among
		// the registered handlers are ones that spend money and touch other
		// people's content: ai_translate rewrites the translations of any
		// article id it is given, syndicate_telegram re-posts to the channel.
		// A plain reader could mint an API key and drive both, so the role list
		// no longer includes "user".
		jobs.WithHTTPMiddleware(
			apiKeyModule.RequireAPIKey(),
			authModule.RequireRoles("operator", "admin"),
		),
		// The operator console reaches the same queue with the credential a
		// browser actually has. LoadSession turns the cookie into claims, which
		// RequireRoles then holds to the same staff roles and to the same check
		// that the account still backs them; SameOriginOnly is there because a
		// cookie-authed POST is a CSRF surface and retry/cancel/enqueue are
		// POSTs.
		jobs.WithConsoleMiddleware(
			authModule.LoadSession,
			auth.SameOriginOnly,
			authModule.RequireRoles("operator", "admin"),
		),
	)
	aiModule := ai.New()
	aiModule.RegisterJobs(jobModule)

	syndicateModule := syndicate.New(notifierModule)
	syndicateModule.RegisterJobs(jobModule)

	jobModule.Handle("send_welcome_email", func(ctx context.Context, rt *shanraq.Runtime, job jobs.Job) error {
		var payload struct {
			Email string `json:"email"`
		}
		if err := job.Decode(&payload); err != nil {
			return err
		}
		meta, _ := jobs.InfoFromContext(ctx)
		rt.Logger.Info("send_welcome_email",
			zap.String("email", payload.Email),
			zap.Int("attempt", meta.Attempts+1),
		)
		return nil
	})

	app := shanraq.New(cfg)
	app.Register(migrations.New())
	app.Register(telemetry.New())
	app.Register(health.New())
	app.Register(notifierModule)
	app.Register(authModule)
	app.Register(apiKeyModule)
	app.Register(jobModule)
	app.Register(aiModule)
	app.Register(syndicateModule)
	mediaModule := media.New(authModule)
	app.Register(mediaModule)
	articlesModule = articles.New(authModule, aiModule, syndicateModule, mediaModule, notifierModule)
	// Narration is synthesised off this server and posted back, so the writer is
	// a script with a key rather than a browser with a cookie.
	articlesModule.UseAPIKeyAuth(apiKeyModule.RequireAPIKey())
	articlesModule.RegisterJobs(jobModule)
	app.Register(articlesModule)
	app.Register(webui.New(jobWorkers, jobPollSeconds,
		webui.WithTenantResolver(func(r *http.Request) (uuid.UUID, bool) {
			return tenantResolver(r)
		}),
		// The operator console is staff-only (global job metrics + error text).
		webui.WithAuthGuard(authModule.RequireSession("/studio/login", "admin", "operator")),
	))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}

func enforceProductionSecrets(cfg config.Config) error {
	if !strings.EqualFold(cfg.Environment, "production") {
		return nil
	}
	secret := cfg.Auth.TokenSecret
	if secret == "" || secret == "replace-me-now" || len(secret) < 32 {
		return fmt.Errorf("insecure auth.token_secret for production environment")
	}
	return nil
}
