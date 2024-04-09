package eventprocessor

import (
	"context"
	"log"

	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/adapter/messaging"
)

// EventProcessor is an interface for processing incoming events.
type EventProcessor interface {
	ProcessEvent(ctx context.Context, event event.Event) error
}

// EventProcessorImpl is a struct implementing the EventProcessor interface.
type EventProcessorImpl struct {
	messagingClient messaging.MessagingSystem
}

// NewEventProcessor creates a new instance of EventProcessorImpl.
func NewEventProcessor(messagingSystem messaging.MessagingSystem) *EventProcessorImpl {
	// Add any necessary initialization logic
	return &EventProcessorImpl{}
}

// ProcessEvent implements the ProcessEvent method of the EventProcessor interface.
func (ep *EventProcessorImpl) ProcessEvent(ctx context.Context, event event.Event) error {
	// Implement the logic to process the event
	// Determine the type of event and take appropriate actions
	log.Printf("Processing event: %v", event)
	return nil
}
