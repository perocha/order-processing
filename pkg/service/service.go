package service

import (
	"context"
	"os"

	"github.com/perocha/order-processing/pkg/appcontext"
	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/domain/order"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/database"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	messagingClient messaging.MessagingSystem
	orderRepo       database.OrderRepository
}

// NewService creates a new instance of ServiceImpl.
func Initialize(ctx context.Context, messagingSystem messaging.MessagingSystem, orderRepository database.OrderRepository) *ServiceImpl {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "services::Initialize::Initializing service logic", telemetry.Information, nil, true)

	messagingClient := messagingSystem
	orderRepo := orderRepository

	return &ServiceImpl{
		messagingClient: messagingClient,
		orderRepo:       orderRepo,
	}
}

// Starts listening for incoming events.
func (s *ServiceImpl) Start(ctx context.Context, signals <-chan os.Signal) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	channel, cancelCtx, err := s.messagingClient.Subscribe(ctx)
	if err != nil {
		telemetryClient.TrackException(ctx, "services::Start::Failed to subscribe to events", err, telemetry.Critical, nil, true)
		return err
	}

	telemetryClient.TrackTrace(ctx, "services::Start::Subscribed to events", telemetry.Information, nil, true)

	for {
		select {
		case message := <-channel:
			// Update the context with the operation ID
			ctx = context.WithValue(ctx, appcontext.OperationIDKeyContextKey, message.OperationID)

			if message.Error != nil {
				// New message received in channel. Process the event.
				properties := message.Event.ToMap()
				telemetryClient.TrackTrace(ctx, "services::Start::Received message", telemetry.Information, properties, true)
				s.processEvent(ctx, message.Event)
			} else {
				// Discard message but report exception
				properties := map[string]string{
					"Error": message.Error.Error(),
				}
				telemetryClient.TrackException(ctx, "services::Start::Error processing message", message.Error, telemetry.Error, properties, true)
			}
		case <-ctx.Done():
			telemetryClient.TrackTrace(ctx, "services::Start::Context canceled. Stopping event listener.", telemetry.Information, nil, true)
			cancelCtx()
			s.messagingClient.Close(ctx)
			return nil
		case <-signals:
			telemetryClient.TrackTrace(ctx, "services::Start::Received termination signal", telemetry.Information, nil, true)
			cancelCtx()
			s.messagingClient.Close(ctx)
			return nil
		}
	}
}

// Stop the service
func (s *ServiceImpl) Stop(ctx context.Context) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "services::Stop::Stopping service", telemetry.Information, nil, true)

	s.messagingClient.Close(ctx)
}

// ProcessEvent processes an incoming event.
func (s *ServiceImpl) processEvent(ctx context.Context, event event.Event) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Based on the event type, determine the action to be taken
	switch event.Type {
	case "create_order":
		// Call the OrderService to create the order
		err := s.createOrder(ctx, event.OrderPayload)
		if err != nil {
			telemetryClient.TrackException(ctx, "services::processEvent::Error creating order", err, telemetry.Error, nil, true)
			return err
		}

		// Publish a message indicating successful operation if needed

		telemetryClient.TrackTrace(ctx, "services::processEvent::Order created", telemetry.Information, nil, true)
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

// CreateOrder implements the CreateOrder method of the OrderService interface.
func (s *ServiceImpl) createOrder(ctx context.Context, order order.Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Log the order creation
	properties := order.ToMap()
	telemetryClient.TrackTrace(ctx, "services::createOrder::Creating order", telemetry.Information, properties, true)

	s.orderRepo.CreateOrder(ctx, order)
	return nil
}
