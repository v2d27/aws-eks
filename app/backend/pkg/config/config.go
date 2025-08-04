package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Host          string
	Port          string
	AllowedOrigin string
	Environment   string
}

// Load loads configuration from environment variables and .env file
func Load() *Config {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	return &Config{
		Host:          getEnvOrDefault("HOST", "0.0.0.0"),
		Port:          getEnvOrDefault("PORT", "8080"),
		AllowedOrigin: getEnvOrDefault("ALLOWED_ORIGIN", "http://localhost:3000"),
		Environment:   getEnvOrDefault("ENV", "development"),
	}
}

// Address returns the full address string
func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
