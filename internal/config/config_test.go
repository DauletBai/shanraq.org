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
