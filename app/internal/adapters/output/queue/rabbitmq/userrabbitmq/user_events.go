package userrabbitmq

import (
	"userservice/internal/core/entities"
)

const (
	createdUserEventType = "user.created"
	updatedUserEventType = "user.updated"
	removedUserEventType = "user.removed"
)

// userEvent defines the base to create User resource related events
type userEvent struct {
	eventType entities.EventType
	userID    entities.ResourceID
}

var _ entities.ServiceEvent = new(userEvent)

func (event userEvent) EventType() entities.EventType {
	return event.eventType
}

func (event userEvent) Payload() interface{} {
	return event.userID
}
