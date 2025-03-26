package user

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/persistence"
)

// Storage defines the port to retrieve User data from the plugged storage
type Storage interface {
	Add(ctx context.Context, data User) (*User, error)
	Get(ctx context.Context, id entities.ResourceID) (*User, error)
	Remove(ctx context.Context, id entities.ResourceID) error
	List(filters UserFilters) (*UsersList, error)
}

// UserFilters defines filter options when listing Users
type UserFilters interface {
	// OrderByCreatedAt orders Users in storage by creation date
	OrderByCreatedAt(direction persistence.OrderDirection) UserFilters
	// OrderByUpdatedAt orders Users in storage by the date when it was last updated
	OrderByUpdatedAt(direction persistence.OrderDirection) UserFilters
	// FilterByCountry filters Users by their country
	FilterByCountry(userID string) UserFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list where the offset parameter defines
	Paged(limit int, offset int) UserFilters
}
