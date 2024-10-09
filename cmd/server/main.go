package main

import (
	"log"
	"net/http"

	"shanraq.org/config"
	"shanraq.org/internal/app"
	"shanraq.org/internal/utils"
)

func main() {
	// Loading configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Logger initialization
	logger := utils.NewLogger(cfg)

	// Router setup
	router := app.SetupRouter(cfg, logger)

	// Server startup
	addr := cfg.Server.Address
	logger.Infof("Server startup on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatalf("Server startup error: %v", err)
	}
}
