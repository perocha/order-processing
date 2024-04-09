package messaging

type MessagingSystem interface {
	Init() error
	StartListening() error
	StopListening() error
}
