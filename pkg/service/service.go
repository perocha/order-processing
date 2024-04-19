package service

import (
	"context"
	"os"

	"github.com/perocha/order-processing/pkg/appcontext"
	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/database"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging"
	"github.com/perocha/order-processing/pkg/infrastructure/telemetry"
)

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	messagingClient messaging.MessagingSystem
	orderRepo       database.OrderRepository
}

// Creates a new instance of ServiceImpl.
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

			if message.Error == nil {
				// New message received in channel. Process the event.
				properties := message.Event.ToMap()
				telemetryClient.TrackTrace(ctx, "services::Start::Received message", telemetry.Information, properties, true)
				s.processEvent(ctx, message.Event)
			} else {
				// Error received. In this case we'll discard message but report an exception
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
		// Create an order in the database
		err := s.orderRepo.CreateOrder(ctx, event.OrderPayload)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error creating order", err, telemetry.Error, properties, true)
			return err
		}

	case "delete_order":
		// Delete an order from the database
		err := s.orderRepo.DeleteOrder(ctx, event.OrderPayload.Id, event.OrderPayload.ProductCategory)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error deleting order", err, telemetry.Error, properties, true)
			return err
		}

	case "update_order":
		// Update an order in the database
		err := s.orderRepo.UpdateOrder(ctx, event.OrderPayload)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error updating order", err, telemetry.Error, properties, true)
			return err
		}

	default:
		// Handle unsupported event types or errors
		telemetryClient.TrackTrace(ctx, "services::processEvent::Unsupported event type", telemetry.Warning, nil, true)
	}

	return nil
}
