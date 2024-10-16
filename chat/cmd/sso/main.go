package main

import (
	"log"
	"shanraq.org/internal/config"
	"shanraq.org/internal/db"
	"shanraq.org/internal/server"
)

func main() {
	// initialization of configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: %v", err)
	}

	// initialization of database
	dbConn, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection error: %v", err)
	}
	defer dbConn.Close()

	// Server startup
	srv := server.NewServer(cfg, dbConn)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server startup errorr: %v", err)
	}
}