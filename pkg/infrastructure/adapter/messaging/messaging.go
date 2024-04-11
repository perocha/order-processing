package messaging

import "context"

type Message struct {
	Topic string
	Body  []byte
}

type MessagingSystem interface {
	Init(ctx context.Context) error
	Publish(ctx context.Context, message Message) error
	//	Subscribe(ctx context.Context, topic string) (<-chan Message, error)
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Unsubscribe(ctx context.Context, topic string) error
	Close(ctx context.Context) error
}
