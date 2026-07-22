// Command migrate applies the embedded database migrations to the DSN in
// DATABASE_URL. It is used by CI to prepare a scratch database for the
// integration tests, and is handy for operational one-off migration runs.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"shanraq.org/pkg/modules/migrations"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := migrations.Apply(ctx, dsn); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied")
}
