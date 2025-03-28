package userstorageinfra

import "userservice/internal/core/entities"

// UserDB is the data struct of the resource in the database
type UserDB struct {
	entities.ResourceID `storage:"id"`
	FirstName           string `storage:"first_name" valid:"required"`
	LastName            string `storage:"last_name" valid:"required"`
	Email               string `storage:"email" valid:"required, email"`
	Password            string `storage:"password" valid:"required"`
	Country             string `storage:"country" valid:"required"`
	Nickname            string `storage:"nickname" valid:"required"`
	CreatedAt           string `storage:"created_at" valid:"required"`
	UpdatedAt           string `storage:"updated_at" valid:"required"`
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
