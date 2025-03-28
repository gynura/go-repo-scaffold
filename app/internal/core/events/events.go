// Package events defines the structure and behavior of this service's events
// It provides a standardized interface for event handling and enforces a controlled
// catalogue of event types to maintain consistency
package events

// ServiceEvent defines how all events from this service should behave
type ServiceEvent interface {
	// EventType returns the catalogue category of the event
	EventType() EventType
	// Payload returns the event's data
	Payload() interface{}
}

// EventType wraps string to enforce the definition of a catalogue of events and avoid passing random strings
type EventType string

func (e EventType) String() string {
	return string(e)
}
