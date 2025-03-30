package userqueue

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/queue"
	"userservice/internal/pkg/errors"
)

// QueueAdapter defines the default implementation of the user.Queue interface
type QueueAdapter struct {
	publisher queue.Publisher
}

// ProvideQueueAdapterOptions defines the options to configure a UserUserQueueAdapter
type ProvideQueueAdapterOptions struct {
	Publisher *queue.Publisher
}

// ProvideQueueAdapter creates a UserUserQueueAdapter with the given options
func ProvideQueueAdapter(options ProvideQueueAdapterOptions) (*QueueAdapter, error) {
	if options.Publisher == nil {
		return nil, errors.NewInternalError("no publisher provided")
	}
	return &QueueAdapter{
		publisher: *options.Publisher,
	}, nil
}

var _ user.Queue = new(QueueAdapter)

func (adapter QueueAdapter) NotifyCreatedUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
	createdUserEvent := userEvent{
		eventType: createdUserEventType,
		userID:    userID,
	}

	err := adapter.publisher.Publish(createdUserEvent)
	if err != nil {
		// TODO: add error translator to be able to control error types from the infra layer
		// If we had the error translator we could do something like:
		// if translator.IsConnectionTimeout() then return errors.NewConnectionTimeout("queue connection timeout").FromError(err)...
		// For now we will wrap all of them as internal errors
		return errors.NewInternalError("error publishing created user event").FromError(err)
	}

	return nil
}

func (adapter QueueAdapter) NotifyRemoveUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
	removeUserEvent := userEvent{
		eventType: removedUserEventType,
		userID:    userID,
	}

	err := adapter.publisher.Publish(removeUserEvent)
	if err != nil {
		return errors.NewInternalError("error publishing remove user event").FromError(err)
	}

	return nil
}

func (adapter QueueAdapter) NotifyEditUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
	editUserEvent := userEvent{
		eventType: editedUserEventType,
		userID:    userID,
	}

	err := adapter.publisher.Publish(editUserEvent)
	if err != nil {
		return errors.NewInternalError("error publishing edit user event").FromError(err)
	}

	return nil
}
