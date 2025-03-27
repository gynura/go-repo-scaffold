package user

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/errors"
)

// Queue defines the functionality to send information through a queue.
// The only information we send is the unique identifier of the resource to keep services decoupled from each other's models
type Queue interface {
	NotifyCreatedUser(ctx context.Context, userID entities.ResourceID) *errors.ServiceError
	NotifyRemoveUser(ctx context.Context, userID entities.ResourceID) *errors.ServiceError
	NotifyEditUser(ctx context.Context, userID entities.ResourceID) *errors.ServiceError
}
