package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/perocha/order-processing/pkg/domain"
)

type EventConsumer struct {
	createOrder CreateOrder
	deleteOrder DeleteOrder
	updateOrder UpdateOrder
}

// Creates a new event consumer
func NewEventConsumer(createOrder CreateOrder, deleteOrder DeleteOrder, updateOrder UpdateOrder) *EventConsumer {
	return &EventConsumer{
		createOrder: createOrder,
		deleteOrder: deleteOrder,
		updateOrder: updateOrder,
	}
}

// Consumes an event and processes it
func (mc *EventConsumer) ConsumeEvent(ctx context.Context, event domain.Event) error {
	switch event.Type {
	case "create_order":
		log.Println("create_order")
		order := mc.convertEventToOrder(event)
		return mc.createOrder.Execute(ctx, order)
	case "delete_order":
		log.Println("delete_order")
		order := mc.convertEventToOrder(event)
		return mc.deleteOrder.Execute(ctx, order.OrderID)
	case "update_order":
		log.Println("update_order")
		order := mc.convertEventToOrder(event)
		return mc.updateOrder.Execute(ctx, order)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return errors.New("unknown event type")
	}
}

// Convert event to order
func (mc *EventConsumer) convertEventToOrder(event domain.Event) domain.Order {
	order := domain.Order{
		OrderID: event.EventID,
	}
	return order
}
