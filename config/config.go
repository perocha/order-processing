package config

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
)

type Config struct {
	AppConfigurationConnectionString string
	EventHubName                     string
	EventHubConnectionString         string
	CheckpointStoreContainerName     string
	CheckpointStoreConnectionString  string
	AppInsightsInstrumentationKey    string
	CosmosDBConnectionString         string
	client                           *azappconfig.Client
}

// InitializeConfig creates a new instance of Config with values loaded from environment variables
func InitializeConfig() *Config {
	cfg := &Config{}

	// Create a new App Configuration client
	connectionString := os.Getenv("APPCONFIGURATION_CONNECTION_STRING")
	if connectionString == "" {
		log.Println("Error: APPCONFIGURATION_CONNECTION_STRING environment variable is not set")
		return nil
	}

	var err error
	cfg.client, err = azappconfig.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Println("Error: Failed to create new App Configuration client")
		return nil
	}

	cfg.AppInsightsInstrumentationKey, _ = cfg.GetVar("APPINSIGHTS_INSTRUMENTATIONKEY")
	cfg.EventHubName, _ = cfg.GetVar("EVENTHUB_NAME")
	cfg.EventHubConnectionString, _ = cfg.GetVar("EVENTHUB_CONSUMERVNEXT_CONNECTION_STRING")
	cfg.CheckpointStoreContainerName, _ = cfg.GetVar("CHECKPOINTSTORE_CONTAINER_NAME")
	cfg.CheckpointStoreConnectionString, _ = cfg.GetVar("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING")
	cfg.CosmosDBConnectionString, _ = cfg.GetVar("COSMOSDB_CONNECTION_STRING")

	return cfg
}

// GetVar retrieves a configuration setting by key from App Configuration
func (cfg *Config) GetVar(key string) (string, error) {
	if cfg.client == nil {
		err := errors.New("app configuration client not initialized")
		log.Println("App configuration client not initialized")
		return "", err
	}

	// Get the setting value from App Configuration
	resp, err := cfg.client.GetSetting(context.TODO(), key, nil)
	if err != nil {
		log.Printf("Error: Failed to get configuration setting %s\n", key)
		return "", err
	}

	return *resp.Value, nil
}
