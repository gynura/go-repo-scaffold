package userrabbitmq

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/queue/rabbitmq"
)

// UserQueueAdapter defines the default implementation of the user.Queue interface
type UserQueueAdapter struct {
	publisher rabbitmq.Publisher
}

// ProvideUserQueueAdapterOptions defines the options to configure a UserUserQueueAdapter
type ProvideUserQueueAdapterOptions struct {
	Publisher rabbitmq.Publisher
}

// ProvideUserQueueAdapter creates a UserUserQueueAdapter with the given options
func ProvideUserQueueAdapter(options ProvideUserQueueAdapterOptions) (*UserQueueAdapter, error) {
	return &UserQueueAdapter{}, nil
}

var _ user.Queue = new(UserQueueAdapter)

func (adapter UserQueueAdapter) NotifyNewUser(ctx context.Context, userID entities.ResourceID) error {
	newUserEvent := userEvent{
		eventType: createdUserEventType,
		userID:    userID,
	}

	err := adapter.publisher.Publish(newUserEvent)
	if err != nil {
		// TODO: add error translator to be able to control error types from the infra layer
		// If we had the error translator we could do something like:
		// if translator.IsConnectionTimeout() then return errors.NewConnectionTimeout("rabbit mq connection timeout").FromError(err)...
		return err
	}

	return nil
}

func (adapter UserQueueAdapter) NotifyRemoveUser(ctx context.Context, userID entities.ResourceID) error {
	//TODO implement me
	panic("implement me")
}

func (adapter UserQueueAdapter) NotifyEditUser(ctx context.Context, userID entities.ResourceID) error {
	//TODO implement me
	panic("implement me")
}
