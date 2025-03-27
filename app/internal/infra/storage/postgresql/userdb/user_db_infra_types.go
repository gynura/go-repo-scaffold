package userdb

import "userservice/internal/core/entities"

// UserDB is the data struct of the resource in the database
type UserDB struct {
	// ResourceID is the ID of the resource
	entities.ResourceID
}

// UserCreateDB is the data struct of the creation of a resource in the database
type UserCreateDB struct {
	// UserDB is the data struct of the resource in the database
	UserDB
}

// UserUpdateDB is the data struct of the update of a resource in the database
type UserUpdateDB struct {
	// UserDB is the data struct of the resource in the database
	UserDB
	// NewResourceVersion is the new resource version after the edition
	NewResourceVersion string `storage:"new_resource_version"`
}

// UserExistsDB is the data struct to check if a resource exists in the database
type UserExistsDB struct {
	// Exists is true if the resource exists
	Exists bool `storage:"exists_result" valid:"required"`
}
