package config

import (
	"os"
)

// Config represents the application configuration.
type Config struct {
	EventHubConnectionString string
	CosmosDBConnectionString string
}

// NewConfig creates a new instance of Config with values loaded from environment variables.
func NewConfig() *Config {
	return &Config{
		EventHubConnectionString: getEnv("EVENT_HUB_CONNECTION_STRING", ""),
		CosmosDBConnectionString: getEnv("COSMOS_DB_CONNECTION_STRING", ""),
		// Initialize other configuration fields from environment variables
	}
}

// getEnv retrieves the value of the specified environment variable or returns a default value if not found.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
