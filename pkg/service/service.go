package service

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/database"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging"
	"github.com/perocha/order-processing/pkg/service/eventprocessor"
	"github.com/perocha/order-processing/pkg/service/orderservice"
)

// Service aggregates all service interfaces.
type Service interface {
	// Add any necessary methods for the service
}

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	eventprocessor.EventProcessor
	orderservice.OrderService
}

// NewService creates a new instance of ServiceImpl.
// func NewService(messagingSystem messaging.MessagingSystem, orderRepository database.OrderRepository) *ServiceImpl {
func NewService(ctx context.Context, messagingSystem messaging.MessagingSystem, orderRepository database.OrderRepository) *ServiceImpl {
	//	eventProcessor := eventprocessor.NewEventProcessor(messagingSystem)
	orderService := orderservice.NewOrderService(orderRepository)

	return &ServiceImpl{
		//		EventProcessor: eventProcessor,
		OrderService: orderService,
	}
}

// ProcessEvent processes an incoming event.
func (s *ServiceImpl) ProcessEvent(ctx context.Context, event event.Event) error {
	// Perform any necessary preprocessing or validation
	// Based on the event type, determine the action to be taken
	switch event.Type {
	case "create_order":
		// Extract order information from the event
		// Call the OrderService to create the order
		// Publish a message indicating successful operation if needed
		err := s.OrderService.CreateOrder(ctx, event.OrderPayload)
		if err != nil {
			return err
		}
	case "delete_order":
		// Extract order ID from the event
		// Call the OrderService to delete the order
		// Publish a message indicating successful operation if needed
	case "update_order":
		// Extract order information from the event
		// Call the OrderService to update the order
		// Publish a message indicating successful operation if needed
	default:
		// Handle unsupported event types or errors
	}

	return nil
}
