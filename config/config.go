package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server struct {
		Address string
	}
	Database struct {
		Host string
		Port int
		User string
		Password string
		DBName string
		SSLMode string
	}
	JWT struct {
		SekretKey string
	}
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("config")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Overriding parameters from environment variables
	if v := os.Getenv("DATABASE PASSWORD"); v != "" {
		cfg.Database.Password = v
	}

	return &cfg, nil
}