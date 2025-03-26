package user

import (
	"context"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/commons/time"
	"userservice/internal/pkg/errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
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
		return nil, errors.NewInvalidArgumentError("no storage provided")
	}
	return &DefaultUseCase{
		storage: options.Storage,
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (d DefaultUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.NewInvalidArgumentError("invalid input to create user").FromError(err).
			WithMessage("Invalid input to create user. Please review your input fields and try again.")
	}

	// As the API defines, if no ID is informed, the service will generate a UUID for the resource
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	now := time.Now()
	user := User{
		MetaInformation: entities.MetaInformation{
			ResourceID: entities.ResourceID{
				ID: input.ID,
			},
			Timestamps: entities.Timestamps{
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		SpecificData: SpecificData{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Country:   input.Country,
			Nickname:  input.Nickname,
		},
	}
	err = user.SetPassword(input.Password)
	if err != nil {
		// TODO: log errors
		return nil, errors.NewInternalError("failed to set user's password").FromError(err)
	}

	storedUser, err := d.storage.Add(ctx, user)
	if err != nil {
		// TODO: log errors
		return nil, errors.NewInternalError("failed to store user").FromError(err)
	}
	// TODO: publish a CreatedUser event here

	return &CreateUserOutput{
		User: *storedUser,
	}, nil
}

func (d DefaultUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultUseCase) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultUseCase) EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	//TODO implement me
	panic("implement me")
}
