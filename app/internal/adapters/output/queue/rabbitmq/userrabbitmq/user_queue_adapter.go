package userrabbitmq

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/queue/rabbitmq"
	"userservice/internal/pkg/errors"
)

// UserQueueAdapter defines the default implementation of the user.Queue interface
type UserQueueAdapter struct {
	publisher rabbitmq.Publisher
}

// ProvideUserQueueAdapterOptions defines the options to configure a UserUserQueueAdapter
type ProvideUserQueueAdapterOptions struct {
	Publisher *rabbitmq.Publisher
}

// ProvideUserQueueAdapter creates a UserUserQueueAdapter with the given options
func ProvideUserQueueAdapter(options ProvideUserQueueAdapterOptions) (*UserQueueAdapter, error) {
	if options.Publisher != nil {
		return nil, errors.NewInternalError("no publisher provided")
	}
	return &UserQueueAdapter{
		publisher: *options.Publisher,
	}, nil
}

var _ user.Queue = new(UserQueueAdapter)

func (adapter UserQueueAdapter) NotifyCreatedUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
	createdUserEvent := userEvent{
		eventType: createdUserEventType,
		userID:    userID,
	}

	err := adapter.publisher.Publish(createdUserEvent)
	if err != nil {
		// TODO: add error translator to be able to control error types from the infra layer
		// If we had the error translator we could do something like:
		// if translator.IsConnectionTimeout() then return errors.NewConnectionTimeout("rabbit mq connection timeout").FromError(err)...
		// For now we will wrap all of them as internal errors
		return errors.NewInternalError("error publishing created user event").FromError(err)
	}

	return nil
}

func (adapter UserQueueAdapter) NotifyRemoveUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
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

func (adapter UserQueueAdapter) NotifyEditUser(_ context.Context, userID entities.ResourceID) *errors.ServiceError {
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
