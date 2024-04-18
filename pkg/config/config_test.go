package config

import (
	"os"
	"testing"
)

func TestConfig_WrongConnectionString(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "wrong_connection_string")

	// Initialize config
	_, err := InitializeConfig("config.yaml")

	// Check if config is nil
	if err == nil {
		t.Error("Expected error for wrong connection string")
	}
}

func TestConfig_GetVar(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "")

	// Initialize config
	cfg, _ := InitializeConfig("config.yaml")

	// Check if config is not nil
	if cfg == nil {
		t.Error("Expected config to not be nil")
	}

	// Test GetVar
	_, err := cfg.GetVar("NON_EXISTING_VAR")
	if err == nil {
		t.Error("Expected error for non-existing variable")
	}

	// Test GetVar with existing variable
	val, err := cfg.GetVar("TEST_VAR")
	if err != nil {
		t.Error("Expected no error for existing variable")
	}
	if val != "TEST_VALUE" {
		t.Error("GetVar returned unexpected value")
	}
}

func TestConfig_InitializeConfigFromYamlFile(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "")

	// Initialize config
	cfg, _ := InitializeConfig("config.yaml")

	// Check if config is not nil
	if cfg == nil {
		t.Error("Expected config to not be nil")
	}
}

func TestConfig_UseWrongFile(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "")

	// Initialize config
	cfg, _ := InitializeConfig("wrong_file.yaml")

	// Check if config is nil
	if cfg != nil {
		t.Error("Expected config to be nil")
	}
}

func TestConfig_UseConfigFileWithWrongVariable(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "")

	// Initialize config
	_, err := InitializeConfig("wrong_var_file.yaml")

	// Check if config is not nil
	if err == nil {
		t.Error("Expected error for wrong variable in config file")
	}
}

func TestConfig_GetVarWithNilClient(t *testing.T) {
	// Set environment variables
	os.Setenv("APPCONFIGURATION_CONNECTION_STRING", "")

	// Initialize config
	cfg, _ := InitializeConfig("config.yaml")

	// Check if config is not nil
	if cfg == nil {
		t.Error("Expected config to not be nil")
	}

	// Set client to nil
	cfg.client = nil

	// Test GetVar
	_, err := cfg.GetVar("TEST_VAR")
	if err == nil {
		t.Error("Expected error for nil client")
	}
}
