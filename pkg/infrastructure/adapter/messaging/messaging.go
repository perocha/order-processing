package messaging

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain/event"
)

type EventWithContext struct {
	Ctx context.Context
	Msg event.Event
}

type MessagingSystem interface {
	Publish(ctx context.Context, event event.Event) error
	// TODO how to deal with "topic" concept?
	Subscribe(ctx context.Context) (<-chan EventWithContext, context.CancelFunc, error)
	Close(ctx context.Context) error
}
