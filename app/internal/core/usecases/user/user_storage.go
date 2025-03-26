package user

import "userservice/internal/core/entities"

// Storage defines the port to retrieve User data from the plugged storage
type Storage interface {
	Add(id entities.ID) (*User, error)
	Get(id entities.ID) (*User, error)
	Remove(id entities.ID) error
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
