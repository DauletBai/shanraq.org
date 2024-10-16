package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	GRPCAddress string
	ServerPort  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		GRPCAddress: viper.GetString("GRPC_ADDRESS"),
		ServerPort: viper.GetString("SERVER_PORT"),
	}

	return cfg, nil
}