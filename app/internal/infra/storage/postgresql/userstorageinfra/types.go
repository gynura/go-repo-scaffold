package userstorageinfra

import "userservice/internal/core/entities"

// UserDB is the data struct of the resource in the database
type UserDB struct {
	entities.ResourceID `storage:"id"`
	FirstName           string `storage:"first_name"`
	LastName            string `storage:"last_name"`
	Email               string `storage:"email"`
	Password            string `storage:"password"`
	Country             string `storage:"country"`
	Nickname            string `storage:"nickname"`
	CreatedAt           string `storage:"created_at"`
	UpdatedAt           string `storage:"updated_at"`
}

// CreateInput is the data struct of the creation of a resource in the database
type CreateInput struct {
	UserDB
}

// EditInput is the data struct of needed to update the resource in the database
type EditInput struct {
	UserDB
}

// ExistsOutput is the data struct to check if a resource exists in the database
type ExistsOutput struct {
	Exists bool `storage:"exists_result" valid:"required"` // true if the resource exists
}
