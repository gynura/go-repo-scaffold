package user

import "userservice/internal/core/entities"

const (
	defaultOrderDirection = entities.OrderDesc
)

// User defines the user model in the service
type User struct {
	entities.MetaInformation
	FirstName string
	LastName  string
	Email     string
	Password  string
	Country   string
	Nickname  string
}

// UsersList defines a list of User resources.
type UsersList struct {
	Items []User
	entities.Pagination
}
