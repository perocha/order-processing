package eventhub

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain"
)

type EventSubscriber struct {
	// Define your dependencies
}

func (es *EventSubscriber) Subscribe(ctx context.Context, handler func(domain.Event)) {
	// Implement your subscription logic
}
