package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SHANRAQ_SERVER_ADDRESS", ":9090")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Server.Address != ":9090" {
		t.Fatalf("expected address :9090, got %s", cfg.Server.Address)
	}
	if cfg.Database.URL == "" {
		t.Fatalf("expected default database url")
	}
	if cfg.Auth.TokenSecret == "" {
		t.Fatalf("expected default auth secret")
	}
}

func TestMetricsTokenFromEnv(t *testing.T) {
	// The nested telemetry.metrics_token key must be readable from the
	// environment (it gates /metrics in production).
	t.Setenv("SHANRAQ_TELEMETRY_METRICS_TOKEN", "tok-123")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Telemetry.MetricsToken != "tok-123" {
		t.Fatalf("expected metrics_token from env, got %q", cfg.Telemetry.MetricsToken)
	}
}
