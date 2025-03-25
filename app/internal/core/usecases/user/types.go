package user

import "userservice/internal/core/entities"

const (
	defaultOrderDirection = entities.OrderDesc
)

// User defines the user model in the service
type User struct {
	// TODO: complete
}

// UsersList defines a list of User resources.
type UsersList struct {
	// Items User in collection.
	Items []User
	// TODO: Add pagination in entities for this
}
