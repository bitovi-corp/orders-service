package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	ProductServiceURL string
	LoyaltyServiceURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", ""),
		LoyaltyServiceURL: getEnv("LOYALTY_SERVICE_URL", ""),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
