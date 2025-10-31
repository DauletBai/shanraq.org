package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/modules/health"
	"shanraq.org/pkg/modules/jobs"
	"shanraq.org/pkg/modules/migrations"
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

	const (
		jobWorkers     = 4
		jobPollSeconds = 2 * time.Second
	)

	jobModule := jobs.New(
		jobs.WithWorkerCount(jobWorkers),
		jobs.WithPollInterval(jobPollSeconds),
	)
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
	app.Register(auth.New())
	app.Register(jobModule)
	app.Register(webui.New(jobWorkers, jobPollSeconds))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}
