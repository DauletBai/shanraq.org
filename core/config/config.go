// github.com/DauletBai/shanraq.org/config/config.go
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// DefaultConfigFileName is the standard name for the config file.
const DefaultConfigFileName = "config"

// Provider defines the interface for a configuration provider.
// This allows for different configuration loading strategies or sources in the future.
type Provider interface {
	// Load attempts to load configuration.
	// Implementations should define their own logic for finding config files
	// (e.g., from specific paths, environment variables, etc.).
	// 'configName' is typically the file name without extension (e.g., "config").
	// 'configPaths' are optional directory paths to search for the config file.
	Load(configName string, configPaths ...string) error

	// Get returns an interface{}.
	Get(key string) interface{}
	// GetString returns the value associated with the key as a string.
	GetString(key string) string
	// GetInt returns the value associated with the key as an int.
	GetInt(key string) int
	// GetBool returns the value associated with the key as a bool.
	GetBool(key string) bool
	// IsSet checks to see if the key has been set in any of the Data sources.
	IsSet(key string) bool
	// UnmarshalKey takes a single key and unmarshals it into a Struct.
	UnmarshalKey(key string, rawVal interface{}) error
	// Unmarshal unmarshals the config into a Struct.
	Unmarshal(rawVal interface{}) error
	// SetDefault sets the default value for a key.
	// Default values are not written to config files.
	SetDefault(key string, value interface{})
}

// ViperProvider implements the Provider interface using github.com/spf13/viper.
type ViperProvider struct {
	vp *viper.Viper
}

// NewViperProvider creates and initializes a new Viper-based configuration provider.
func NewViperProvider() *ViperProvider {
	v := viper.New()
	// Configure Viper
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // e.g., server.port becomes SERVER_PORT
	v.AutomaticEnv()                                            // Read environment variables
	v.AllowEmptyEnv(true)                                       // Allow empty env vars to not cause errors

	// You can set default config type, but it's often better to let Viper infer it
	// or set it based on the file found.
	// v.SetConfigType("yaml")

	return &ViperProvider{vp: v}
}

// Load loads configuration using Viper.
// It searches for a config file (e.g., "config.yaml", "config.json") in the provided paths
// or default locations if no paths are given.
// Environment variables can override file configurations.
func (p *ViperProvider) Load(configName string, configPaths ...string) error {
	if configName == "" {
		return errors.New("config name cannot be empty")
	}
	p.vp.SetConfigName(configName) // Name of config file (without extension)

	if len(configPaths) > 0 {
		for _, path := range configPaths {
			p.vp.AddConfigPath(path)
		}
	} else {
		// Default search paths for applications using the framework
		p.vp.AddConfigPath(".")
		p.vp.AddConfigPath("./configs") // Common practice
		p.vp.AddConfigPath("/etc/shanraq/") // For system-wide config on Linux
		// For user-specific config (e.g., $HOME/.config/shanraq/config.yaml)
		// Getting home dir can be tricky, usually handled by app initializing config
		if home, err := os.UserHomeDir(); err == nil {
			p.vp.AddConfigPath(home + "/.config/shanraq")
			p.vp.AddConfigPath(home + "/.shanraq")
		}
	}

	err := p.vp.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; this might be acceptable if all config is via ENV or defaults
			// The application/kernel can decide if this is a fatal error.
			// For the framework, we just report it.
			return fmt.Errorf("config file '%s' not found: %w", configName, err)
		}
		// Config file was found but another error was produced
		return fmt.Errorf("failed to read config file '%s': %w", p.vp.ConfigFileUsed(), err)
	}
	// fmt.Printf("Using config file: %s\n", p.vp.ConfigFileUsed()) // Useful for debugging
	return nil
}

func (p *ViperProvider) Get(key string) interface{}            { return p.vp.Get(key) }
func (p *ViperProvider) GetString(key string) string           { return p.vp.GetString(key) }
func (p *ViperProvider) GetInt(key string) int                 { return p.vp.GetInt(key) }
func (p *ViperProvider) GetBool(key string) bool               { return p.vp.GetBool(key) }
func (p *ViperProvider) IsSet(key string) bool                 { return p.vp.IsSet(key) }
func (p *ViperProvider) UnmarshalKey(key string, rawVal interface{}) error { return p.vp.UnmarshalKey(key, rawVal) }
func (p *ViperProvider) Unmarshal(rawVal interface{}) error    { return p.vp.Unmarshal(rawVal) }
func (p *ViperProvider) SetDefault(key string, value interface{}) { p.vp.SetDefault(key, value) }

// Ensure ViperProvider implements Provider
var _ Provider = (*ViperProvider)(nil)