package usecase

import "github.com/perocha/order-processing/pkg/domain"

type EventInteractor struct {
	// Define your dependencies such as repositories
}

func (ei *EventInteractor) ProcessEvent(event domain.Event) {
	// Define your business logic for processing an event
}
