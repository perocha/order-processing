package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/perocha/goadapters/database"
	"github.com/perocha/goadapters/messaging"
	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/domain/order"
)

// ServiceImpl is a struct implementing the Service interface.
type ServiceImpl struct {
	consumerInstance messaging.MessagingSystem
	producerInstance messaging.MessagingSystem
	//	orderRepo        database.OrderRepository
	orderRepo database.DBRepository
}

// Creates a new instance of ServiceImpl.
func Initialize(ctx context.Context, consumer messaging.MessagingSystem, producer messaging.MessagingSystem, orderRepository database.DBRepository) *ServiceImpl {
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
			ctx = context.WithValue(ctx, telemetry.OperationIDKeyContextKey, message.OperationID)

			if message.Error == nil {
				// New message received in channel. Process the event.
				telemetryClient.TrackTrace(ctx, "services::Start::Received message", telemetry.Information, nil, true)

				eventData, ok := message.Data.(map[string]interface{})
				if !ok {
					// Error received. In this case we'll discard message but report an exception
					properties := map[string]string{
						"Error": "Failed to cast message to map[string]interface{}",
					}
					telemetryClient.TrackException(ctx, "services::Start::Error processing message", nil, telemetry.Error, properties, true)
					continue
				}

				event, err := mapToEvent(eventData)
				if err != nil {
					// Error received. In this case we'll discard message but report an exception
					properties := map[string]string{
						"Error": err.Error(),
					}
					telemetryClient.TrackException(ctx, "services::Start::Error processing message", err, telemetry.Error, properties, true)
					continue
				}

				err = s.processEvent(ctx, event)
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
					"Error": message.Error.Error(),
				}
				telemetryClient.TrackException(ctx, "services::Start::Error processing message", message.Error, telemetry.Error, properties, true)
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
func (s *ServiceImpl) processEvent(ctx context.Context, event event.Event) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Based on the event type, determine the action to be taken
	switch event.Type {

	case "create_order":
		// Create an order in the database
		err := s.orderRepo.CreateDocument(ctx, event.OrderPayload.ProductCategory, event.OrderPayload)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error creating order", err, telemetry.Error, properties, true)
			return err
		}

	case "delete_order":
		// Delete an order from the database
		err := s.orderRepo.DeleteDocument(ctx, event.OrderPayload.ProductCategory, event.OrderPayload.Id)
		if err != nil {
			properties := map[string]string{
				"Error": err.Error(),
			}
			telemetryClient.TrackException(ctx, "services::processEvent::Error deleting order", err, telemetry.Error, properties, true)
			return err
		}

	case "update_order":
		// Update an order in the database
		err := s.orderRepo.UpdateDocument(ctx, event.OrderPayload.ProductCategory, event.OrderPayload.Id, event.OrderPayload)
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
	response := ResponseMessage{
		MessageID: uuid.New().String(),
		EventID:   event.EventID,
		Error:     nil,
		Status:    "Processed",
	}
	s.producerInstance.Publish(ctx, response)

	return nil
}

// Function to convert message to event type with order information
func mapToEvent(data map[string]interface{}) (event.Event, error) {
	event := event.Event{}

	var ok bool

	// Extract values from the map and convert to the Event struct
	if event.Type, ok = data["Type"].(string); !ok {
		return event, errors.New("type field not found or not a string")
	}
	if event.EventID, ok = data["EventID"].(string); !ok {
		return event, errors.New("eventID field not found or not a string")
	}

	// Parse Timestamp field
	timestampStr, ok := data["Timestamp"].(string)
	if !ok {
		return event, errors.New("timestamp field not found or not a string")
	}
	timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
	if err != nil {
		return event, errors.New("failed to parse Timestamp: " + err.Error())
	}
	event.Timestamp = timestamp

	// Extract OrderPayload from the map
	orderPayload, ok := data["OrderPayload"].(map[string]interface{})
	if !ok {
		return event, errors.New("orderPayload field not found or not a map")
	}

	// Convert OrderPayload to order.Order struct
	event.OrderPayload = order.Order{
		Id:              orderPayload["id"].(string),
		ProductCategory: orderPayload["ProductCategory"].(string),
		ProductID:       orderPayload["productId"].(string),
		CustomerID:      orderPayload["customerId"].(string),
		Status:          orderPayload["status"].(string),
	}

	return event, nil
}

// Stop the service
func (s *ServiceImpl) Stop(ctx context.Context) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "services::Stop::Stopping service", telemetry.Information, nil, true)

	s.consumerInstance.Close(ctx)
}
