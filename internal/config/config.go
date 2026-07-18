package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config is the top-level runtime configuration for the framework runtime.
type Config struct {
	Environment   string              `mapstructure:"environment"`
	PublicBaseURL string              `mapstructure:"public_base_url"`
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Telemetry     Telemetry           `mapstructure:"telemetry"`
	Logging       Logging             `mapstructure:"logging"`
	Auth          AuthConfig          `mapstructure:"auth"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
	AI            AIConfig            `mapstructure:"ai"`
	Syndicate     SyndicateConfig     `mapstructure:"syndicate"`
	Media         MediaConfig         `mapstructure:"media"`
	Social        SocialConfig        `mapstructure:"social"`
}

// SocialConfig holds public social-profile URLs shown in the top info bar.
// Empty entries are hidden.
type SocialConfig struct {
	Telegram  string `mapstructure:"telegram"`
	Instagram string `mapstructure:"instagram"`
	YouTube   string `mapstructure:"youtube"`
	Facebook  string `mapstructure:"facebook"`
}

// MediaConfig controls user media (images today; video via third-party embeds).
// Storage is pluggable behind a Store interface: "fs" keeps files on local disk
// under Dir and serves them from PublicPrefix (drop-in S3/MinIO backend is the
// media-2 step). Every uploaded image is re-encoded — which strips EXIF (incl.
// GPS) for author safety — resized to MaxDimension, and stamped with the brand
// watermark in the top-right corner.
type MediaConfig struct {
	Backend        string `mapstructure:"backend"`
	Dir            string `mapstructure:"dir"`
	PublicPrefix   string `mapstructure:"public_prefix"`
	MaxUploadBytes int64  `mapstructure:"max_upload_bytes"`
	MaxDimension   int    `mapstructure:"max_dimension"`
	Watermark      bool   `mapstructure:"watermark"`
}

// SyndicateConfig controls resilience channels: RSS is always on; Telegram
// auto-posting activates only when a bot token and chat are configured.
type SyndicateConfig struct {
	BaseURL  string         `mapstructure:"base_url"`
	Telegram TelegramConfig `mapstructure:"telegram"`
}

type TelegramConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	BotToken string `mapstructure:"bot_token"`
	ChatID   string `mapstructure:"chat_id"`
}

// AIConfig controls the optional Claude-backed writing assistant. It stays
// disabled until an API key is provided, so the app runs with zero AI spend.
type AIConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	APIKey         string `mapstructure:"api_key"`
	EditorModel    string `mapstructure:"editor_model"`
	TranslateModel string `mapstructure:"translate_model"`
	MaxTokens      int    `mapstructure:"max_tokens"`
}

// ServerConfig configures the embedded HTTP server.
type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig configures the PostgreSQL pool.
type DatabaseConfig struct {
	URL               string        `mapstructure:"url"`
	MaxConns          int32         `mapstructure:"max_conns"`
	MinConns          int32         `mapstructure:"min_conns"`
	MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime   time.Duration `mapstructure:"max_conn_idle_time"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
}

// Telemetry configures health/metrics endpoints.
type Telemetry struct {
	EnableMetrics bool          `mapstructure:"enable_metrics"`
	MetricsPath   string        `mapstructure:"metrics_path"`
	MetricsToken  string        `mapstructure:"metrics_token"`
	Tracing       TracingConfig `mapstructure:"tracing"`
}

type TracingConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Endpoint    string  `mapstructure:"endpoint"`
	Insecure    bool    `mapstructure:"insecure"`
	SampleRatio float64 `mapstructure:"sample_ratio"`
	ServiceName string  `mapstructure:"service_name"`
}

// Logging configures zap logger.
type Logging struct {
	Level string `mapstructure:"level"`
	Mode  string `mapstructure:"mode"`
}

// AuthConfig controls token generation and lifecycle.
type AuthConfig struct {
	TokenSecret string        `mapstructure:"token_secret"`
	TokenTTL    time.Duration `mapstructure:"token_ttl"`
	MFA         MFAConfig     `mapstructure:"mfa"`
}

type MFAConfig struct {
	TOTP TOTPConfig `mapstructure:"totp"`
}

type TOTPConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Issuer  string `mapstructure:"issuer"`
}

type NotificationsConfig struct {
	SMTP SMTPConfig `mapstructure:"smtp"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// Load reads configuration by merging defaults, optional file, and env vars.
func Load(configPath string) (Config, error) {
	_ = godotenv.Load() // ignore failure when .env missing

	v := viper.New()
	setDefaults(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("load config file: %w", err)
		}
	}

	v.SetEnvPrefix("SHANRAQ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")

	v.SetDefault("server.address", ":8080")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")

	v.SetDefault("database.url", "postgres://postgres:postgres@127.0.0.1:5432/shanraq?sslmode=disable")
	v.SetDefault("database.max_conns", 10)
	v.SetDefault("database.min_conns", 1)
	v.SetDefault("database.max_conn_lifetime", "30m")
	v.SetDefault("database.max_conn_idle_time", "5m")
	v.SetDefault("database.health_check_period", "30s")

	v.SetDefault("telemetry.enable_metrics", true)
	v.SetDefault("telemetry.metrics_path", "/metrics")
	// Registering the key is required for viper's AutomaticEnv to bind
	// SHANRAQ_TELEMETRY_METRICS_TOKEN during Unmarshal (nested keys are only
	// read from the environment once they are known).
	v.SetDefault("telemetry.metrics_token", "")
	v.SetDefault("telemetry.tracing.enabled", false)
	v.SetDefault("telemetry.tracing.endpoint", "")
	v.SetDefault("telemetry.tracing.sample_ratio", 0.1)
	v.SetDefault("telemetry.tracing.service_name", "shanraq-app")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.mode", "production")

	v.SetDefault("auth.token_secret", "replace-me-now")
	v.SetDefault("auth.token_ttl", "15m")
	v.SetDefault("auth.mfa.totp.enabled", false)
	v.SetDefault("auth.mfa.totp.issuer", "Shanraq")

	v.SetDefault("notifications.smtp.port", 587)

	v.SetDefault("ai.enabled", false)
	v.SetDefault("ai.editor_model", "claude-sonnet-5")
	v.SetDefault("ai.translate_model", "claude-haiku-4-5")
	v.SetDefault("ai.max_tokens", 4096)

	v.SetDefault("public_base_url", "http://localhost:8080")
	v.SetDefault("syndicate.base_url", "http://localhost:8080")
	v.SetDefault("syndicate.telegram.enabled", false)

	v.SetDefault("media.backend", "fs")
	v.SetDefault("media.dir", "./data/media")
	v.SetDefault("media.public_prefix", "/media")
	v.SetDefault("media.max_upload_bytes", 10<<20) // 10 MiB
	v.SetDefault("media.max_dimension", 2000)
	v.SetDefault("media.watermark", true)

	// Registered so viper's AutomaticEnv binds SHANRAQ_SOCIAL_* (nested keys are
	// only read from the environment once known). Empty = the icon is hidden.
	v.SetDefault("social.telegram", "")
	v.SetDefault("social.instagram", "")
	v.SetDefault("social.youtube", "")
	v.SetDefault("social.facebook", "")
}

// PublicBase returns the canonical public origin (no trailing slash), used for
// SEO canonical/hreflang links, RSS/Telegram links, and password-reset links.
// Syndicate.base_url is a fallback for backward compatibility.
func (c Config) PublicBase() string {
	v := strings.TrimRight(strings.TrimSpace(c.PublicBaseURL), "/")
	if v == "" {
		v = strings.TrimRight(strings.TrimSpace(c.Syndicate.BaseURL), "/")
	}
	return v
}

func validateConfig(cfg Config) error {
	var problems []string

	if strings.EqualFold(cfg.Environment, "production") && weakAuthSecret(cfg.Auth.TokenSecret) {
		problems = append(problems, "auth.token_secret must be at least 32 characters and not use default values in production")
	}

	// In production the public base URL must be a real HTTPS origin — it ends up
	// in canonical/hreflang tags, RSS/Telegram links, and password-reset emails.
	if strings.EqualFold(cfg.Environment, "production") {
		base := cfg.PublicBase()
		if !strings.HasPrefix(base, "https://") || strings.Contains(base, "localhost") || strings.Contains(base, "127.0.0.1") {
			problems = append(problems, "public_base_url must be set to a real https:// origin in production (not localhost)")
		}
	}

	// /metrics must not be world-readable in production: require a bearer token
	// (or disable metrics). See telemetry module for enforcement.
	if strings.EqualFold(cfg.Environment, "production") && cfg.Telemetry.EnableMetrics &&
		strings.TrimSpace(cfg.Telemetry.MetricsToken) == "" {
		problems = append(problems, "telemetry.metrics_token is required in production when telemetry.enable_metrics is true (protects /metrics)")
	}

	smtp := cfg.Notifications.SMTP
	if smtpConfigured(smtp) {
		if strings.TrimSpace(smtp.Host) == "" {
			problems = append(problems, "notifications.smtp.host is required when configuring SMTP")
		}
		if strings.TrimSpace(smtp.From) == "" {
			problems = append(problems, "notifications.smtp.from is required when configuring SMTP")
		}
		if smtp.Port <= 0 {
			problems = append(problems, "notifications.smtp.port must be greater than zero when configuring SMTP")
		}
	}

	tracing := cfg.Telemetry.Tracing
	if tracing.Enabled {
		if strings.TrimSpace(tracing.Endpoint) == "" {
			problems = append(problems, "telemetry.tracing.endpoint is required when tracing is enabled")
		}
		if tracing.SampleRatio < 0 || tracing.SampleRatio > 1 {
			problems = append(problems, "telemetry.tracing.sample_ratio must be between 0 and 1")
		}
	}

	totp := cfg.Auth.MFA.TOTP
	if totp.Enabled && strings.TrimSpace(totp.Issuer) == "" {
		problems = append(problems, "auth.mfa.totp.issuer is required when TOTP is enabled")
	}

	if len(problems) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(problems, "; "))
	}
	return nil
}

func smtpConfigured(cfg SMTPConfig) bool {
	return strings.TrimSpace(cfg.Host) != "" ||
		strings.TrimSpace(cfg.From) != "" ||
		strings.TrimSpace(cfg.Username) != "" ||
		strings.TrimSpace(cfg.Password) != ""
}

func weakAuthSecret(secret string) bool {
	normalized := strings.ToLower(strings.TrimSpace(secret))
	if normalized == "" {
		return true
	}
	if len(secret) < 32 {
		switch normalized {
		case "replace-me-now", "super-secret-key", "secret", "changeme":
			return true
		}
	}
	return len(secret) < 32
}
