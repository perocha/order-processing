package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/perocha/order-processing/pkg/appcontext"

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
func (t *Telemetry) TrackTrace(ctx context.Context, message string, severity SeverityLevel, properties map[string]string, logToConsole bool) {
	if t.client == nil {
		panic("Telemetry client not initialized")
	}

	// Create the log message
	logMessage := fmt.Sprintf("Message: %s, Properties: %v, Severity: %v\n", message, properties, severity)

	// If logToConsole is true, print the log message
	if logToConsole {
		log.Println(logMessage)
	}

	// Create the new trace
	trace := appinsights.NewTraceTelemetry(message, contracts.SeverityLevel(severity))
	for k, v := range properties {
		trace.Properties[k] = v
	}

	// Get the operationID from the context
	if operationID, ok := ctx.Value(appcontext.OperationIDKeyContextKey).(string); ok {
		// Set parent id
		if operationID != "" {
			trace.Tags.Operation().SetParentId(operationID)
		}
	}

	// Send the trace to App Insights
	t.client.Track(trace)
}
