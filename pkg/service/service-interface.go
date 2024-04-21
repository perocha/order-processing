package service

// Message to send once the events are processed
type ResponseMessage struct {
	MessageID string `json:"message_id"`
	EventID   string `json:"event_id"`
	Error     error  `json:"error"`
	Status    string `json:"status"`
}
