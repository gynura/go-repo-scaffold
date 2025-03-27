package userrabbitmq

import (
	"userservice/internal/core/entities"
	"userservice/internal/core/events"
)

const (
	createdUserEventType = "user.created"
	editedUserEventType  = "user.edited"
	removedUserEventType = "user.removed"
)

// userEvent defines the base to create User resource related events
type userEvent struct {
	eventType events.EventType
	userID    entities.ResourceID
}

var _ events.ServiceEvent = new(userEvent)

func (event userEvent) EventType() events.EventType {
	return event.eventType
}

func (event userEvent) Payload() interface{} {
	return event.userID
}
