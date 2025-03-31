package userqueue

import (
	"encoding/json"
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

func (event userEvent) Payload() ([]byte, error) {
	payload, err := json.Marshal(event.userID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
