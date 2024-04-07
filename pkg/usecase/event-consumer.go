package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

// EventConsumer defines the interface for consuming events
type EventConsumer interface {
	ConsumeEvent(ctx context.Context, event domain.Event) error
}

// eventConsumerImpl implements the EventConsumer interface
type eventConsumerImpl struct {
	orderProcessor OrderProcessImpl
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(orderProcessor OrderProcessImpl) EventConsumer {
	return &eventConsumerImpl{
		orderProcessor: orderProcessor,
	}
}

// Consumes an event and processes it
func (e *eventConsumerImpl) ConsumeEvent(ctx context.Context, event domain.Event) error {
	log.Printf("event-consumer::Type: %s::EventID: %s::Timestamp: %s", event.Type, event.EventID, event.Timestamp)

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
func (e *eventConsumerImpl) convertEventToOrder(event domain.Event) domain.Order {
	order := domain.Order{
		OrderID: event.EventID,
	}
	return order
}
