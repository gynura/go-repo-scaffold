package user

import (
	"context"
	"errors"
)

// UseCase defines the functionality to manage User resources
type UseCase interface {
	// CreateUser creates a User. It returns the created User or an error if it fails
	CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
	// ListUsers returns all the Users or an error if it fails
	ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
	// GetUser returns the requested User or an error if it fails
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	// EditUser edits a User. It returns the edited User or an error if it fails
	EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, error)
	// DeleteUser deletes a User. It returns the deleted User or an error if it fails
	DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error)
}

// DefaultUseCase defines the default implementation of the user.UseCase interface
type DefaultUseCase struct {
	storage Storage
}

// ProvideUseCaseOptions defines the options to configure a UserUseCase implementation
type ProvideUseCaseOptions struct {
	Storage Storage
}

// ProvideUseCase creates a UserUseCase with the given options.
func ProvideUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.New("no storage provided")
	}
	return &DefaultUseCase{
		storage: options.Storage,
	}, nil
}

var _ UseCase = new(DefaultUseCase)
