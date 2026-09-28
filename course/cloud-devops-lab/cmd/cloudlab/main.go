package main

import (
	"log/slog"
	"os"

	"example.com/cloudlab/internal/app"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := app.ConfigFromEnv(version)
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(2)
	}
	if err := app.Run(cfg, logger); err != nil {
		logger.Error("service stopped", "error", err)
		os.Exit(1)
	}
}
