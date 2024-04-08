package event

import (
	"context"
	"errors"
	"log"

	"github.com/perocha/order-processing/pkg/domain/order"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

// EventConsumer defines the interface for consuming events
type EventConsumer interface {
	ConsumeEvent(ctx context.Context, event Event) error
}

// eventConsumerImpl implements the EventConsumer interface
type eventConsumerImpl struct {
	orderProcessor order.OrderProcessImpl
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(orderProcessor order.OrderProcessImpl) EventConsumer {
	return &eventConsumerImpl{
		orderProcessor: orderProcessor,
	}
}

// Consumes an event and processes it
func (e *eventConsumerImpl) ConsumeEvent(ctx context.Context, event Event) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// New event received, process it
	properties := event.ToMap()
	telemetryClient.TrackTrace(ctx, "ConsumeEvent::New event received, processing order", telemetry.Information, properties, true)

	switch event.Type {
	case "create_order":
		order := e.convertEventToOrder(event)
		return e.orderProcessor.ProcessCreateOrder(ctx, order)
	case "delete_order":
		order := e.convertEventToOrder(event)
		return e.orderProcessor.ProcessDeleteOrder(ctx, order.OrderID)
	case "update_order":
		order := e.convertEventToOrder(event)
		return e.orderProcessor.ProcessUpdateOrder(ctx, order)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return errors.New("unknown event type")
	}
}

// Convert event to order
func (e *eventConsumerImpl) convertEventToOrder(event Event) order.Order {
	order := order.Order{
		OrderID: event.EventID,
	}
	return order
}
