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
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Telemetry     Telemetry           `mapstructure:"telemetry"`
	Logging       Logging             `mapstructure:"logging"`
	Auth          AuthConfig          `mapstructure:"auth"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
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
	EnableMetrics bool   `mapstructure:"enable_metrics"`
	MetricsPath   string `mapstructure:"metrics_path"`
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

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.mode", "production")

	v.SetDefault("auth.token_secret", "replace-me-now")
	v.SetDefault("auth.token_ttl", "15m")

	v.SetDefault("notifications.smtp.port", 587)
}
