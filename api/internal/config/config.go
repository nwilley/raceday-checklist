package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Port: envOrDefault("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     envOrDefault("DB_HOST", "127.0.0.1"),
			Port:     envOrDefault("DB_PORT", "3306"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}

	if err := cfg.Database.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg DatabaseConfig) Validate() error {
	if cfg.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return fmt.Errorf("DB_PORT must be a valid port: %w", err)
	}
	return nil
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
