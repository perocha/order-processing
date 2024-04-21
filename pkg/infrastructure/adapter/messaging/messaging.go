package messaging

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain/event"
)

// EventWithOperationID struct that contains an operation ID and an Event
type Message struct {
	OperationID string
	Error       error
	Event       event.Event
}

type MessagingSystem interface {
	Publish(ctx context.Context, data interface{}) error
	// TODO how to deal with "topic" concept?
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}
