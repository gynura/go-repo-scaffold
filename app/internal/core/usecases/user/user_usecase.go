package user

import (
	"context"
	"log/slog"
	"os"
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
	logger  *slog.Logger
	queue   Queue
}

// ProvideUseCaseOptions defines the options to configure a UserUseCase implementation
type ProvideUseCaseOptions struct {
	Storage Storage
	Queue   Queue
}

// ProvideUseCase creates a UserUseCase with the given options.
func ProvideUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.NewInvalidArgumentError("no storage provided")
	}
	if options.Queue == nil {
		return nil, errors.NewInvalidArgumentError("no queue provided")
	}
	return &DefaultUseCase{
		storage: options.Storage,
		queue:   options.Queue,
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (useCase DefaultUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
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

	storedUser, addErr := useCase.storage.Add(ctx, user)
	if addErr != nil {
		// TODO: better this error log
		useCase.logger.Error(addErr.Error())
		if addErr.IsAlreadyExistsError() {
			return nil, errors.NewAlreadyExistsError("user %s already exists", user.ID).FromError(err).
				WithMessage("user %s already exists", user.ID)
		}
		return nil, errors.NewInternalError("failed to store user").FromError(err)
	}

	err = useCase.queue.NotifyCreatedUser(ctx, storedUser.ResourceID)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		User: *storedUser,
	}, nil
}

func (useCase DefaultUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	//TODO implement me
	panic("implement me")
}
