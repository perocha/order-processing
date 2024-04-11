package eventprocessor

import (
	"context"
	"log"
	"time"

	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging"
)

// EventProcessor is an interface for processing incoming events.
type EventProcessor interface {
	ProcessEvent(ctx context.Context, event event.Event) error
	StartListening(ctx context.Context) error
}

// EventProcessorImpl is a struct implementing the EventProcessor interface.
type EventProcessorImpl struct {
	messagingClient messaging.MessagingSystem
}

// NewEventProcessor creates a new instance of EventProcessorImpl.
func NewEventProcessor(messagingSystem messaging.MessagingSystem) *EventProcessorImpl {
	return &EventProcessorImpl{
		messagingClient: messagingSystem,
	}
}

// ProcessEvent implements the ProcessEvent method of the EventProcessor interface.
func (ep *EventProcessorImpl) ProcessEvent(ctx context.Context, event event.Event) error {
	// Implement the logic to process the event
	// Determine the type of event and take appropriate actions
	log.Printf("Processing event: %v", event)
	return nil
}

// StartListening starts listening for incoming events.
func (ep *EventProcessorImpl) StartListening(ctx context.Context) error {
	// Implement the logic to start listening for incoming events
	log.Println("Listening for incoming events")

	channel, err := ep.messagingClient.Subscribe(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to events: %s", err.Error())
		return err
	}

	for {
		select {
		case message := <-channel:
			log.Printf("Received message: %s", string(message.Body))
			event := event.Event{}
			ep.ProcessEvent(ctx, event)
		case <-ctx.Done():
			log.Println("Context canceled. Stopping event listener.")
			return nil
		case <-time.After(1 * time.Minute):
			// Do nothing
			log.Println("Waiting...")
		}
	}

	return nil
}
