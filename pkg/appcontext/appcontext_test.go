package appcontext

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if OperationIDKeyContextKey != "operationID" {
		t.Errorf("Expected OperationIDKeyContextKey to be 'operationID', but got %s", OperationIDKeyContextKey)
	}

	if TelemetryContextKey != "telemetry" {
		t.Errorf("Expected TelemetryContextKey to be 'telemetry', but got %s", TelemetryContextKey)
	}
}
