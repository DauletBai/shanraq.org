// github.com/DauletBai/shanraq.org/core/kernel/kernel.go 
package kernel

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DauletBai/shanraq.org/core/config"
	"github.com/DauletBai/shanraq.org/core/logger"
	"github.com/DauletBai/shanraq.org/database"
)

// AppEnvironment defines the application environment (e.g., development, production).
type AppEnvironment string

const (
	EnvDevelopment AppEnvironment = "development"
	EnvStaging     AppEnvironment = "staging"
	EnvProduction  AppEnvironment = "production"
	EnvTesting     AppEnvironment = "testing"
)

// Kernel is the core of a Shanraq application.
// It holds and manages essential services.
type Kernel struct {
	configProvider config.Provider
	logger         logger.Interface
	db             *sql.DB
	basePath       string         
	environment    AppEnvironment 
	mu 			   sync.RWMutex 
}

// Option is a functional option for configuring the Kernel.
type Option func(*Kernel) error

// New creates a new Kernel instance with the given options.
// It initializes essential services like config and logger.
func New(opts ...Option) (*Kernel, error) {
	k := &Kernel{
		environment: EnvDevelopment, // Default to development
	}

	execPath, err := os.Executable()
	if err == nil {
		k.basePath = filepath.Dir(execPath)
	} else {
		cwd, _ := os.Getwd()
		k.basePath = cwd
	}

	// Apply functional options
	for _, opt := range opts {
		if err := opt(k); err != nil {
			return nil, fmt.Errorf("failed to apply kernel option: %w", err)
		}
	}

	// Initialize default config provider if not already set by an option
	if k.configProvider == nil {
		if err := WithDefaultConfigFile()(k); err != nil { 
			warningMsg := fmt.Sprintf("kernel warning: could not load default config: %v. Continuing with defaults/env vars.", err)
			if k.logger != nil {
				k.logger.Warn(warningMsg)
			} else {
				fmt.Fprintln(os.Stderr, warningMsg)
			}
		}
	}
	
	// Determine environment from config if set
	if k.configProvider != nil && k.configProvider.IsSet("app.environment") {
		envStr := k.configProvider.GetString("app.environment")
		k.SetEnvironment(AppEnvironment(strings.ToLower(envStr)))
	}

	// Initialize logger if not already set by an option
	if k.logger == nil {
		logLevelStr := "INFO"
		logJSON := false
		logAddSource := false
		if k.configProvider != nil {
			// Try to get logger settings from config
			logLevelStr = k.configProvider.GetString("logger.level")
			if logLevelStr == "" { // if GetString returns "" for unset, provide default
				logLevelStr = "INFO"
			}
			logJSON = k.configProvider.GetBool("logger.json_output")
			logAddSource = k.configProvider.GetBool("logger.add_source")
		}
		// Add source in development by default
		if k.Environment() == EnvDevelopment && !k.configProvider.IsSet("logger.add_source") {
		    logAddSource = true
		}

		defaultLogger := logger.NewSlogLogger(logger.LogLevel(strings.ToUpper(logLevelStr)), os.Stdout, logJSON, logAddSource)
		k.logger = defaultLogger
	}
	
	k.logger.Info("Shanraq Kernel initialized", "environment", k.environment, "basePath", k.basePath)
    if k.configProvider != nil {
        configFile := k.configProvider.ConfigFileUsed()
        if configFile != "" {
            k.logger.Info("Configuration loaded", "file", configFile)
        } else if k.configProvider.IsSet("app.name") { 
             k.logger.Info("Configuration loaded (no file, using environment variables or defaults)")
        } else {
            k.logger.Info("No configuration file loaded and no alternative sources found.")
        }
    }

	// INITIALIZING CONNECTION TO DB
	if k.configProvider != nil && (k.configProvider.IsSet("database.dsn") || k.configProvider.IsSet("database.host")) {
		dbConn, err := database.NewConnection(k.configProvider, k.logger) 
		if err != nil {
			k.logger.Error("Failed to connect to database during kernel initialization", "error", err)
		}
		k.db = dbConn
	} else {
		k.logger.Info("Database connection is not configured, skipping initialization.")
	}

	return k, nil
}

func (k *Kernel) Config() config.Provider {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.configProvider
}

func (k *Kernel) Logger() logger.Interface {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.logger
}

func (k *Kernel) DB() *sql.DB {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.db
}

func (k *Kernel) CloseDB() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.db != nil {
		k.logger.Info("Closing database connection...")
		return k.db.Close()
	}
	return nil
}

// BasePath returns the application's base path.
func (k *Kernel) BasePath() string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.basePath
}

// SetBasePath sets the application's base path.
func (k *Kernel) SetBasePath(path string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.basePath = path
}

// Environment returns the current application environment.
func (k *Kernel) Environment() AppEnvironment {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.environment
}

// SetEnvironment sets the application environment.
// It's good practice to set this early, before many services are initialized,
// as their behavior might depend on the environment.
func (k *Kernel) SetEnvironment(env AppEnvironment) {
	k.mu.Lock()
	defer k.mu.Unlock()
	switch env {
	case EnvDevelopment, EnvStaging, EnvProduction, EnvTesting:
		k.environment = env
	default:
		// Fallback or log warning if an unknown environment is set
		if k.logger != nil {
			k.logger.Warn(fmt.Sprintf("Unknown environment '%s' specified, defaulting to '%s'", env, EnvDevelopment))
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Unknown environment '%s' specified, defaulting to '%s'\n", env, EnvDevelopment)
		}
		k.environment = EnvDevelopment
	}
}

// func (k *Kernel) DB() *sql.DB { 
// 	k.mu.RLock()
// 	defer k.mu.RUnlock()
// 	return k.db
// }

// func (k *Kernel) CloseDB() error { 
// 	k.mu.Lock()
// 	defer k.mu.Unlock()
// 	if k.db != nil {
// 		k.logger.Info("Closing database connection...")
// 		return k.db.Close()
// 	}
// 	return nil
// }

func WithDB(db *sql.DB) Option {
	return func(k *Kernel) error {
		if db == nil {
			return errors.New("database connection (*sql.DB) cannot be nil if provided via option")
		}
		k.db = db
		return nil
	}
}

// --- Kernel Options ---
// WithConfigProvider sets a custom configuration provider.
func WithConfigProvider(provider config.Provider) Option {
	return func(k *Kernel) error {
		if provider == nil {
			return errors.New("config provider cannot be nil")
		}
		k.configProvider = provider
		return nil
	}
}

// WithConfigFile loads configuration from a specific file and paths using ViperProvider.
// If 'configName' is empty, config.DefaultConfigFileName ("config") is used.
// If 'paths' are empty, default search paths of ViperProvider will be used.
func WithConfigFile(configName string, paths ...string) Option {
	return func(k *Kernel) error {
		if configName == "" {
			configName = config.DefaultConfigFileName
		}
		provider := config.NewViperProvider()
		if err := provider.Load(configName, paths...); err != nil {
			// This error is not fatal for kernel creation itself, but will be logged by New().
			// It allows the application to start even if a config file is missing,
			// relying on environment variables or defaults.
			return fmt.Errorf("failed to load config '%s': %w", configName, err)
		}
		k.configProvider = provider
		return nil
	}
}

// WithDefaultConfigFile is a convenience option to load 'config.yaml' (or .json, .toml etc.)
// from standard application paths.
func WithDefaultConfigFile() Option {
	return WithConfigFile(config.DefaultConfigFileName) // Uses default paths in ViperProvider
}

// WithLogger sets a custom logger.
func WithLogger(l logger.Interface) Option {
	return func(k *Kernel) error {
		if l == nil {
			return errors.New("logger cannot be nil")
		}
		k.logger = l
		return nil
	}
}

// WithBasePath sets the application's base path.
func WithBasePath(path string) Option {
	return func(k *Kernel) error {
		if path == "" {
			return errors.New("base path cannot be empty")
		}
		// Optionally, check if path exists and is a directory
		// info, err := os.Stat(path)
		// if err != nil {
		// 	return fmt.Errorf("base path error: %w", err)
		// }
		// if !info.IsDir() {
		// 	return fmt.Errorf("base path '%s' is not a directory", path)
		// }
		k.basePath = path
		return nil
	}
}

// WithEnvironment sets the application environment directly.
func WithEnvironment(env AppEnvironment) Option {
	return func(k *Kernel) error {
		k.SetEnvironment(env) // Uses the setter which includes validation
		return nil
	}
}