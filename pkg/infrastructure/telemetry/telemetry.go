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

// SeverityLevel defines the telemetry severity level
type SeverityLevel contracts.SeverityLevel

// Telemetry severity levels
const (
	Verbose     SeverityLevel = SeverityLevel(contracts.Verbose)
	Information SeverityLevel = SeverityLevel(contracts.Information)
	Warning     SeverityLevel = SeverityLevel(contracts.Warning)
	Error       SeverityLevel = SeverityLevel(contracts.Error)
	Critical    SeverityLevel = SeverityLevel(contracts.Critical)
)

// Initializes a new telemetry client
func Initialize(instrumentationKey string) (*Telemetry, error) {
	if instrumentationKey == "" {
		return nil, errors.New("app insights instrumentation key not initialized")
	}

	// Initialize telemetry client
	client := appinsights.NewTelemetryClient(instrumentationKey)

	return &Telemetry{client: client}, nil
}

// TrackTrace sends a trace telemetry event
func (t *Telemetry) TrackTrace(message string, severity SeverityLevel) {
	t.client.TrackTrace(message, contracts.SeverityLevel(severity))
	// Additional logic for tracking telemetry
}
