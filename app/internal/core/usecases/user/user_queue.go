package user

import (
	"context"
	"userservice/internal/core/entities"
)

// Queue defines the functionality to send information through a queue.
// The only information we send is the unique identifier of the resource to keep services decoupled from each other's models
type Queue interface {
	NotifyNewUser(ctx context.Context, userID entities.ResourceID) error
	NotifyRemoveUser(ctx context.Context, userID entities.ResourceID) error
	NotifyEditUser(ctx context.Context, userID entities.ResourceID) error
}
