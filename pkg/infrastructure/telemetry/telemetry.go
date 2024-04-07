package telemetry

import (
	"errors"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

// Telemetry defines the telemetry client
type Telemetry struct {
	client appinsights.TelemetryClient
}

// Initializes a new telemetry client
func Initialize(instrumentationKey string) (*Telemetry, error) {
	if instrumentationKey == "" {
		return nil, errors.New("app insights instrumentation key not initialized")
	}

	client := appinsights.NewTelemetryClient(instrumentationKey)
	// Set additional configuration options if needed

	return &Telemetry{client: client}, nil
}

// TrackTrace sends a trace telemetry event
func (t *Telemetry) TrackTrace(message string, severity contracts.SeverityLevel) {
	t.client.TrackTrace(message, severity)
	// Additional logic for tracking telemetry
}
