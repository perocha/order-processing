package appcontext

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

// OperationIDKey represents the key type for the operation ID in context
type OperationIDKey string
type TelemetryObj string

const (
	// OperationIDKeyContextKey is the key used to store the operation ID in context
	OperationIDKeyContextKey OperationIDKey = "operationID"

	// TelemetryContextKey represents the key type for the telemetry object in context
	TelemetryContextKey TelemetryObj = "telemetry"
)

// Helper function to retrieve the telemetry client from the context
func GetTelemetryClient(ctx context.Context) *telemetry.Telemetry {
	telemetryClient, ok := ctx.Value(TelemetryContextKey).(*telemetry.Telemetry)
	if !ok {
		log.Panic("Telemetry client not found in context")
	}
	return telemetryClient
}
