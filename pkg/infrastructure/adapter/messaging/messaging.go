package messaging

import (
	"context"

	"github.com/perocha/order-processing/pkg/domain/event"
)

type MessagingSystem interface {
	Init(ctx context.Context) error
	Publish(ctx context.Context, event event.Event) error
	// TODO how to deal with "topic" concept?
	Subscribe(ctx context.Context) (<-chan event.Event, context.CancelFunc, error)
	Unsubscribe(ctx context.Context, topic string) error
	Close(ctx context.Context) error
}
