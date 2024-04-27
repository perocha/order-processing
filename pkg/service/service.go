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
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "services::Initialize::Initializing service logic")

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
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	channel, cancelCtx, err := s.consumerInstance.Subscribe(ctx)
	if err != nil {
		xTelemetry.Error(ctx, "services::Start::Failed to subscribe to events", telemetry.String("Error", err.Error()))
		return err
	}

	xTelemetry.Info(ctx, "services::Start::Subscribed to events")

	for {
		select {
		case message := <-channel:
			// Update the context with the operation ID
			ctx = context.WithValue(ctx, telemetry.OperationIDKeyContextKey, message.GetOperationID())

			if message.GetError() == nil {
				// New message received in channel. Process the event.
				xTelemetry.Debug(ctx, "services::Start::Received message")

				// Get the data from, that is expected to be an Order
				receivedOrder := order.NewOrder()
				err := receivedOrder.Deserialize(message.GetData())
				if err != nil {
					// Error deserializing message. Log and continue
					xTelemetry.Error(ctx, "services::Start::Error deserializing message", telemetry.String("Error", err.Error()))
					continue
				}

				// If we reach this point, we have a valid event. Process it!!!
				err = s.processEvent(ctx, message.GetOperationID(), message.GetCommand(), *receivedOrder)
				if err != nil {
					// Error processing message. Log and continue
					xTelemetry.Error(ctx, "services::Start::Error processing message", telemetry.String("Error", err.Error()))
				}
			} else {
				// Error received. In this case we'll discard message but report an exception
				xTelemetry.Error(ctx, "services::Start::Error processing message", telemetry.String("Error", message.GetError().Error()))
			}
		case <-ctx.Done():
			xTelemetry.Info(ctx, "services::Start::Context canceled. Stopping event listener")
			cancelCtx()
			s.consumerInstance.Close(ctx)
			return nil
		case <-signals:
			xTelemetry.Info(ctx, "services::Start::Received termination signal")
			cancelCtx()
			s.consumerInstance.Close(ctx)
			return nil
		}
	}
}

// ProcessEvent processes an incoming event.
func (s *ServiceImpl) processEvent(ctx context.Context, operationID string, command string, orderInfo order.Order) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// Based on the event type, determine the action to be taken
	switch command {

	case "create_order":
		// Create an order in the database
		err := s.orderRepo.CreateDocument(ctx, orderInfo.ProductCategory, orderInfo)
		if err != nil {
			xTelemetry.Error(ctx, "services::processEvent::Error creating order", telemetry.String("Error", err.Error()))
			return err
		}

	case "delete_order":
		// Delete an order from the database
		err := s.orderRepo.DeleteDocument(ctx, orderInfo.ProductCategory, orderInfo.Id)
		if err != nil {
			xTelemetry.Error(ctx, "services::processEvent::Error deleting order", telemetry.String("Error", err.Error()))
			return err
		}

	case "update_order":
		// Update an order in the database
		err := s.orderRepo.UpdateDocument(ctx, orderInfo.ProductCategory, orderInfo.Id, orderInfo)
		if err != nil {
			xTelemetry.Error(ctx, "services::processEvent::Error updating order", telemetry.String("Error", err.Error()))
			return err
		}

	default:
		// Handle unsupported event types or errors
		xTelemetry.Error(ctx, "services::processEvent::Unsupported event type", telemetry.String("Command", command))
	}

	// Event processed successfully, we'll publish a message to event hub to confirm
	response := message.NewMessage(operationID, nil, "Processed", "", nil)
	s.producerInstance.Publish(ctx, response)

	return nil
}

// Stop the service
func (s *ServiceImpl) Stop(ctx context.Context) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Info(ctx, "services::Stop::Stopping service")

	s.consumerInstance.Close(ctx)
}
