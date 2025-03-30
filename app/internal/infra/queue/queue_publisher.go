package queue

import "userservice/internal/core/events"

// Publisher defines the abstraction to plug any queue publisher service
type Publisher interface {
	// Publish publishes the informed service event through a queue
	Publish(event events.ServiceEvent) error
	// IsAlive returns true if the connection is alive or false if it's not
	IsAlive() bool
}
