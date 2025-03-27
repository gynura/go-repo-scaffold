package user

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"
)

// Storage defines the port to retrieve User data from the plugged storage
type Storage interface {
	Add(ctx context.Context, data User) (*User, *errors.ServiceError)
	Get(ctx context.Context, id entities.ResourceID) (*User, *errors.ServiceError)
	Remove(ctx context.Context, id entities.ResourceID) *errors.ServiceError
	List(filters StorageFilters) (*UsersList, *errors.ServiceError)
}

// StorageFilters defines filter options when listing Users
type StorageFilters interface {
	// OrderByCreatedAt orders Users in storage by creation date
	OrderByCreatedAt(direction persistence.OrderDirection) StorageFilters
	// OrderByUpdatedAt orders Users in storage by the date when it was last updated
	OrderByUpdatedAt(direction persistence.OrderDirection) StorageFilters
	// FilterByCountry filters Users by their country
	FilterByCountry(userID string) StorageFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list where the offset parameter defines
	Paged(limit int, offset int) StorageFilters
}
