package event

import "time"

type Event struct {
	Type      string
	EventID   string
	Timestamp time.Time
}

// Convert Event struct into a map[string]string
func (e *Event) ToMap() map[string]string {
	return map[string]string{
		"Type":      e.Type,
		"EventID":   e.EventID,
		"Timestamp": e.Timestamp.Format(time.RFC3339),
	}
}
