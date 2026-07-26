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

// Secrets and operator identity must be injectable from the environment (.env)
// so they never live in a committed config file.
func TestSecretsAndOperatorFromEnv(t *testing.T) {
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_HOST", "smtp.example.com")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_USERNAME", "mailer")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_PASSWORD", "s3cret")
	t.Setenv("SHANRAQ_NOTIFICATIONS_SMTP_FROM", "no-reply@shanraq.org")
	t.Setenv("SHANRAQ_OPERATOR_BIN", "251040024862")
	t.Setenv("SHANRAQ_OPERATOR_LEGAL_NAME", "Qazna Technologies")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Notifications.SMTP.Host != "smtp.example.com" || cfg.Notifications.SMTP.Username != "mailer" ||
		cfg.Notifications.SMTP.Password != "s3cret" || cfg.Notifications.SMTP.From != "no-reply@shanraq.org" {
		t.Fatalf("SMTP not bound from env: %+v", cfg.Notifications.SMTP)
	}
	if cfg.Operator.BIN != "251040024862" || cfg.Operator.LegalName != "Qazna Technologies" {
		t.Fatalf("operator not bound from env: %+v", cfg.Operator)
	}
}
