package domain

import "time"

type Event struct {
	Type      string
	EventID   string
	Payload   map[string]interface{}
	Timestamp time.Time
}

// Define your event methods
