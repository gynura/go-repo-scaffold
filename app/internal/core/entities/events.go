package entities

// ServiceEvent defines how all events from this service should behave
type ServiceEvent interface {
	// EventType returns the catalogue category of the event
	EventType() EventType
	// Payload returns the event's data
	Payload() interface{}
}

// EventType wraps string to enforce the definition of a catalogue of events and not passing random strings
type EventType string

func (e EventType) String() string {
	return string(e)
}
