package config

import (
	"time"

	"go-test-api/internal/database"
)

const (
	EnvProduction  = "production"
	EnvDevelopment = "development"
)

// Config holds all application configuration
type Config struct {
	Environment string
	Port        string
	Database    database.Config
	JWTSecret   string
	JWTExpiry   time.Duration
}

// Load reads configuration from environment variables.
// In production, critical values (JWT_SECRET, DB_PASSWORD) must be set.
// In development, sensible defaults are provided.
func Load() Config {
	environment := getEnv("ENV", EnvDevelopment)
	isProduction := environment == EnvProduction

	// Critical secrets - required in production, optional in development
	var jwtSecret string
	var dbPassword string

	if isProduction {
		jwtSecret = mustGetEnv("JWT_SECRET")
		dbPassword = mustGetEnv("DB_PASSWORD")
	} else {
		jwtSecret = getEnv("JWT_SECRET", "dev-secret-key-not-for-production")
		dbPassword = getEnv("DB_PASSWORD", "gopassword")
	}

	return Config{
		Environment: environment,
		Port:        getEnv("PORT", "8080"),
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "gouser"),
			Password: dbPassword,
			DBName:   getEnv("DB_NAME", "gotestdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWTSecret: jwtSecret,
		JWTExpiry: 24 * time.Hour,
	}
}

// IsProduction returns true if running in production environment
func (c Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsDevelopment returns true if running in development environment
func (c Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}
