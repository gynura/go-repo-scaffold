package user

import (
	"userservice/internal/core/entities"
)

const (
	defaultOrderDirection = entities.OrderDesc
	defaultLimit          = 100
)

// UsersList defines a paginated list of User resources.
type UsersList struct {
	Items []entities.User
	entities.Pagination
}

// CreateUserInput represents the required input data to create a new User.
type CreateUserInput struct {
	entities.ResourceID `valid:"optional"`
	FirstName           string `valid:"required"`
	LastName            string `valid:"required"`
	Email               string `valid:"required,email"`
	Password            string `valid:"required"`
	Country             string `valid:"required"`
	Nickname            string `valid:"required"`
}

// CreateUserOutput represents the response structure after successfully creating a User.
type CreateUserOutput struct {
	entities.User
}

// GetUserInput represents the input parameters to fetch a specific User by ID.
type GetUserInput struct {
	entities.ResourceID `valid:"required"`
}

// GetUserOutput represents the response structure containing the requested User data.
type GetUserOutput struct {
	entities.User
}

// EditUserInput represents the input parameters required to update an existing User.
type EditUserInput struct {
	entities.ResourceID `valid:"required"`
	FirstName           string `valid:"required"`
	LastName            string `valid:"required"`
	Email               string `valid:"required,email"`
	Password            string `valid:"required"`
	Country             string `valid:"required"`
	Nickname            string `valid:"required"`
}

// EditUserOutput represents the response structure after updating a User.
type EditUserOutput struct {
	entities.User
}

// DeleteUserInput represents the input parameters required to delete a User.
type DeleteUserInput struct {
	entities.ResourceID `valid:"required"`
}

// DeleteUserOutput represents the response structure after deleting a User.
type DeleteUserOutput struct {
	entities.User
}

// ListUsersInput represents the input parameters to retrieve a paginated list of Users.
type ListUsersInput struct {
	PageLimit      int    `valid:"optional,natural"`
	PageOffset     int    `valid:"optional,natural"`
	OrderBy        string `valid:"optional"`
	OrderDirection string `valid:"optional"`
	Country        string `valid:"optional"`
}

// ListUsersOutput represents the response structure containing a paginated list of Users.
type ListUsersOutput struct {
	UsersList
}
