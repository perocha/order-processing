package application

import (
	"log"

	"github.com/perocha/order-processing/domain/event"
	"github.com/perocha/order-processing/domain/order"
	"github.com/perocha/order-processing/pkg/domain/event"
	"github.com/perocha/order-processing/pkg/infrastructure/database"
	"github.com/perocha/order-processing/pkg/infrastructure/messaging"
)

type Orchestrator struct {
	messagingSystem messaging.MessagingSystem
	database        database.Database
	eventProcessor  event.EventProcessor
	orderProcessor  order.OrderProcessor
}

func NewOrchestrator(messagingSystem messaging.MessagingSystem, database database.Database, eventProcessor event.EventProcessor, orderProcessor order.OrderProcessor) *Orchestrator {
	return &Orchestrator{
		messagingSystem: messagingSystem,
		database:        database,
		eventProcessor:  eventProcessor,
		orderProcessor:  orderProcessor,
	}
}

func (o *Orchestrator) Init() error {
	// Initialize the messaging system and the database
	if err := o.messagingSystem.Init(); err != nil {
		return err
	}
	if err := o.database.Init(); err != nil {
		return err
	}
	return nil
}

func (o *Orchestrator) StartListening() error {
	// Start listening for events
	return o.messagingSystem.StartListening()
}

func (o *Orchestrator) ProcessEvent(event event.Event) error {
	// Process an event
	if err := o.eventProcessor.ProcessEvent(event); err != nil {
		return err
	}
	// Depending on the type of the event, process the order
	switch event.Type {
	case "create_order":
		return o.orderProcessor.CreateOrder(event.EventID)
	case "delete_order":
		return o.orderProcessor.DeleteOrder(event.EventID)
	case "update_order":
		return o.orderProcessor.UpdateOrder(event.EventID)
	default:
		log.Printf("Unknown event type: %s\n", event.Type)
		return nil
	}
}
