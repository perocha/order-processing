package domain

import "time"

type Event struct {
	Type      string
	EventID   string
	Timestamp time.Time
}

// Define your event methods
