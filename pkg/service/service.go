package service

import (
	"context"
	"os"

	"github.com/perocha/goadapters/database"
	"github.com/perocha/goadapters/messaging/message"
	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/perocha/order-processing/pkg/domain/order"
)

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	consumerInstance message.MessagingSystem
	producerInstance message.MessagingSystem
	orderRepo        database.DBRepository
}

// Creates a new instance of ServiceImpl.
func Initialize(ctx context.Context, consumer message.MessagingSystem, producer message.MessagingSystem, orderRepository database.DBRepository) *ServiceImpl {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "services::Initialize::Initializing service logic", telemetry.Information, nil, true)

	consumerInstance := consumer
	producerInstance := producer
	orderRepo := orderRepository

	return &ServiceImpl{
		consumerInstance: consumerInstance,
		producerInstance: producerInstance,
		orderRepo:        orderRepo,
	}
}

// Starts listening for incoming events.
func (s *ServiceImpl) Start(ctx context.Context, signals <-chan os.Signal) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	channel, cancelCtx, err := s.consumerInstance.Subscribe(ctx)
	if err != nil {
		telemetryClient.TrackException(ctx, "services::Start::Failed to subscribe to events", err, telemetry.Critical, nil, true)
		return err
	}

	telemetryClient.TrackTrace(ctx, "services::Start::Subscribed to events", telemetry.Information, nil, true)

	for {
		select {
		case message := <-channel:
			// Update the context with the operation ID
			ctx = context.WithValue(ctx, telemetry.OperationIDKeyContextKey, message.GetOperationID())

			if message.GetError() == nil {
				// New message received in channel. Process the event.
				telemetryClient.TrackTrace(ctx, "services::Start::Received message", telemetry.Information, nil, true)

				// Get the data from, that is expected to be an Order
				receivedOrder := order.NewOrder()
				err := receivedOrder.Deserialize(message.GetData())
				if err != nil {
					// Error deserializing message. Log and continue
					properties := map[string]string{
						"Error": err.Error(),
					}
					telemetryClient.TrackException(ctx, "services::Start::Error deserializing message", err, telemetry.Error, properties, true)
					continue
				}

				// If we reach this point, we have a valid event. Process it!!!
				err = s.processEvent(ctx, message.GetOperationID(), message.GetCommand(), *receivedOrder)
				if err != nil {
					// Error processing message. Log and continue
					properties := map[string]string{
						"Error": err.Error(),
					}
					telemetryClient.TrackException(ctx, "services::Start::Error processing message", err, telemetry.Error, properties, true)
				}
			} else {
				// Error received. In this case we'll discard message but report an exception
				properties := map[string]string{
					"Error": message.GetError().Error(),
				}
				telemetryClient.TrackException(ctx, "services::Start::Error processing message", message.GetError(), telemetry.Error, properties, true)
			}
		case <-ctx.Done():
			telemetryClient.TrackTrace(ctx, "services::Start::Context canceled. Stopping event listener.", telemetry.Information, nil, true)
			cancelCtx()
			s.consumerInstance.Close(ctx)
			return nil
		case <-signals:
			telemetryClient.TrackTrace(ctx, "services::Start::Received termination signal", telemetry.Information, nil, true)
			cancelCtx()
			s.consumerInstance.Close(ctx)
			return nil
		}
	}
}

// ProcessEvent processes an incoming event.
func (s *ServiceImpl) processEvent(ctx context.Context, operationID string, command string, orderInfo order.Order) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Based on the event type, determine the action to be taken
	switch command {

	case "create_order":
		// Create an order in the database
		err := s.orderRepo.CreateDocument(ctx, orderInfo.ProductCategory, orderInfo)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error creating order", err, telemetry.Error, properties, true)
			return err
		}

	case "delete_order":
		// Delete an order from the database
		err := s.orderRepo.DeleteDocument(ctx, orderInfo.ProductCategory, orderInfo.Id)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error deleting order", err, telemetry.Error, properties, true)
			return err
		}

	case "update_order":
		// Update an order in the database
		err := s.orderRepo.UpdateDocument(ctx, orderInfo.ProductCategory, orderInfo.Id, orderInfo)
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

	// Event processed successfully, we'll publish a message to event hub to confirm
	response := message.NewMessage(operationID, nil, "Processed", "", nil)
	s.producerInstance.Publish(ctx, response)

	return nil
}

// Stop the service
func (s *ServiceImpl) Stop(ctx context.Context) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "services::Stop::Stopping service", telemetry.Information, nil, true)

	s.consumerInstance.Close(ctx)
}
