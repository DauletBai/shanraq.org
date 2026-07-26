package config

import (
	"strings"
	"testing"
)

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

// Secrets and operator identity must be injectable from the environment (.env)
// so they never live in a committed config file.
func TestSecretsAndOperatorFromEnv(t *testing.T) {
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_HOST", "smtp.example.com")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_USERNAME", "mailer")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_PASSWORD", "s3cret")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_FROM", "no-reply@shanraq.org")
	t.Setenv("SHANRAQ_OPERATOR_BIN", "123456789012")
	t.Setenv("SHANRAQ_OPERATOR_LEGAL_NAME", "Test Company")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Notifications.SMTP.Host != "smtp.example.com" || cfg.Notifications.SMTP.Username != "mailer" ||
		cfg.Notifications.SMTP.Password != "s3cret" || cfg.Notifications.SMTP.From != "no-reply@shanraq.org" {
		t.Fatalf("SMTP not bound from env: %+v", cfg.Notifications.SMTP)
	}
	if cfg.Operator.BIN != "123456789012" || cfg.Operator.LegalName != "Test Company" {
		t.Fatalf("operator not bound from env: %+v", cfg.Operator)
	}
}

// Production must refuse to start without the legal operator disclosure.
func TestOperatorRequiredInProduction(t *testing.T) {
	// Satisfy the other production guards so only the operator check can fail.
	t.Setenv("SHANRAQ_ENVIRONMENT", "production")
	t.Setenv("SHANRAQ_PUBLIC_BASE_URL", "https://shanraq.org")
	t.Setenv("SHANRAQ_AUTH_TOKEN_SECRET", "a-sufficiently-long-token-secret-1234567890")
	t.Setenv("SHANRAQ_TELEMETRY_ENABLE_METRICS", "false")

	if _, err := Load(""); err == nil || !strings.Contains(err.Error(), "operator") {
		t.Fatalf("expected an operator-required error in production, got %v", err)
	}

	// With the operator identity present, production config loads.
	t.Setenv("SHANRAQ_OPERATOR_LEGAL_NAME", "Test Company")
	t.Setenv("SHANRAQ_OPERATOR_BIN", "123456789012")
	t.Setenv("SHANRAQ_OPERATOR_ADDRESS", "Rudny, KZ")
	t.Setenv("SHANRAQ_OPERATOR_EMAIL", "support@shanraq.org")
	if _, err := Load(""); err != nil {
		t.Fatalf("production config with operator should load, got %v", err)
	}
}

func TestBootstrapAdminFromEnv(t *testing.T) {
	t.Setenv("SHANRAQ_BOOTSTRAP_ADMIN_EMAIL", "ceo@shanraq.org")
	t.Setenv("SHANRAQ_BOOTSTRAP_ADMIN_PASSWORD", "strong-enough-pass")
	t.Setenv("SHANRAQ_BOOTSTRAP_ADMIN_FIRST", "Daulet")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Bootstrap.AdminEmail != "ceo@shanraq.org" || cfg.Bootstrap.AdminPassword != "strong-enough-pass" ||
		cfg.Bootstrap.AdminFirst != "Daulet" {
		t.Fatalf("bootstrap not bound from env: %+v", cfg.Bootstrap)
	}
}
