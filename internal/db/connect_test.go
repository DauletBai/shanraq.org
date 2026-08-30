package db

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"shanraq.org/internal/config"
)

func poolCfg(url string) config.DatabaseConfig {
	return config.DatabaseConfig{
		URL: url, MaxConns: 4, MinConns: 0,
		MaxConnLifetime: time.Hour, MaxConnIdleTime: time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

// A malformed URL is the operator's typo, not a database that might come back,
// so it must fail at once rather than spend the retry budget on it.
func TestAMalformedURLFailsWithoutRetrying(t *testing.T) {
	start := time.Now()
	pool, err := Connect(context.Background(), poolCfg("://not a dsn"), zap.NewNop())
	if err == nil {
		pool.Close()
		t.Fatal("a malformed URL was accepted")
	}
	if !strings.Contains(err.Error(), "parse db url") {
		t.Errorf("error = %v; it should say the URL could not be parsed", err)
	}
	if d := time.Since(start); d > time.Second {
		t.Errorf("took %v; a URL that cannot be parsed should not be retried", d)
	}
}

// Connecting lazily and calling it healthy is the failure this guards against:
// the pool is pinged, so a database that is not there fails startup instead of
// looking fine until the first reader arrives.
func TestAnUnreachableDatabaseFailsStartup(t *testing.T) {
	if testing.Short() {
		t.Skip("this one waits out the connect retries")
	}
	// Port 1 is reserved and refuses immediately, so this exercises the retry
	// path without waiting on a network timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	pool, err := Connect(ctx, poolCfg("postgres://u:p@127.0.0.1:1/none?sslmode=disable"), zap.NewNop())
	if err == nil {
		pool.Close()
		t.Fatal("a pool was returned for a database that is not listening")
	}
	if !strings.Contains(err.Error(), "connect postgres") {
		t.Errorf("error = %v; it should name the failure to connect", err)
	}
}

// The settings a pool is given are the ones the operator wrote down.
func TestTheConfiguredLimitsReachThePool(t *testing.T) {
	dsn := os.Getenv("SHANRAQ_TEST_DB")
	if dsn == "" {
		t.Skip("set SHANRAQ_TEST_DB to connect to a real database")
	}
	cfg := poolCfg(dsn)
	cfg.MaxConns, cfg.MinConns = 7, 1
	pool, err := Connect(context.Background(), cfg, zap.NewNop())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if got := pool.Config().MaxConns; got != 7 {
		t.Errorf("MaxConns = %d, want 7", got)
	}
	if got := pool.Config().MinConns; got != 1 {
		t.Errorf("MinConns = %d, want 1", got)
	}
	if pool.Config().ConnConfig.Tracer == nil {
		t.Error("the pool was built without the tracer, so nothing would be logged or traced")
	}
	// It answers, which is what Ping was supposed to have established.
	if err := pool.Ping(context.Background()); err != nil {
		t.Errorf("ping: %v", err)
	}
}
