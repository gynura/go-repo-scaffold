package user

import (
	"golang.org/x/crypto/bcrypt"
	"userservice/internal/core/entities"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// password type wraps a string to enforce hashing on the pass
type password string

// User defines the user model in the service.
type User struct {
	entities.MetaInformation
	SpecificData
}

// SpecificData defines the attributes that are not part of the meta information of a User resource.
type SpecificData struct {
	FirstName string
	LastName  string
	Email     string   `valid:"email"`
	password  password // private field to prevent direct assignment from outside the package
	Country   string
	Nickname  string
}

// SetPassword encapsulates the password assignment so it is not possible to set a password that is not hashed
// from outside this package where the User resource logic ends up becoming more diffuse
func (u *User) SetPassword(rawPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = password(hashed)
	return nil
}

// UsersList defines a paginated list of User resources.
type UsersList struct {
	Items []User
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
	User
}

// GetUserInput represents the input parameters to fetch a specific User by ID.
type GetUserInput struct {
	entities.ResourceID `valid:"required"`
}

// GetUserOutput represents the response structure containing the requested User data.
type GetUserOutput struct {
	User User
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
	User User
}

// DeleteUserInput represents the input parameters required to delete a User.
type DeleteUserInput struct {
	entities.ResourceID `valid:"required"`
}

// DeleteUserOutput represents the response structure after deleting a User.
type DeleteUserOutput struct {
	User User
}

// ListUsersInput represents the input parameters to retrieve a paginated list of Users.
type ListUsersInput struct {
	PageLimit      int    `valid:"optional,natural"`
	PageOffset     int    `valid:"optional,natural"`
	OrderBy        string `valid:"optional"`
	OrderDirection string `valid:"optional"`
}

// ListUsersOutput represents the response structure containing a paginated list of Users.
type ListUsersOutput struct {
	UsersList
}
