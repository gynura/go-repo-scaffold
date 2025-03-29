package user

import (
	"context"
	"fmt"
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
	CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, *errors.ServiceError)
	// ListUsers returns all the Users or an error if it fails
	ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, *errors.ServiceError)
	// GetUser returns the requested User or an error if it fails
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, *errors.ServiceError)
	// EditUser edits a User. It returns the edited User or an error if it fails
	EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, *errors.ServiceError)
	// DeleteUser deletes a User. It returns the deleted User or an error if it fails
	DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, *errors.ServiceError)
}

// DefaultUseCase defines the default implementation of the user.UseCase interface
type DefaultUseCase struct {
	storage Storage
	logger  *slog.Logger
	//queue   Queue
}

// ProvideUseCaseOptions defines the options to configure a UserUseCase implementation
type ProvideUseCaseOptions struct {
	Storage Storage
	//Queue   Queue
}

// ProvideDefaultUseCase provides the default implementation of the User use case logic
func ProvideDefaultUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.NewInvalidArgumentError("no storage provided")
	}
	//if options.Queue == nil {
	//	return nil, errors.NewInvalidArgumentError("no queue provided")
	//}
	return &DefaultUseCase{
		storage: options.Storage,
		//queue:   options.Queue,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (useCase DefaultUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.NewInvalidArgumentError("invalid input to create user").FromError(err).
			WithMessage("Invalid input to create user. Please review your input fields and try again.")
	}

	// Nicknames are unique alongside IDs and Emails, in order to avoid coupling this check to the plugged db, the business logic needs to implement this checks
	if useCase.nicknameAlreadyExists(ctx, input.Nickname) {
		return nil, errors.NewInvalidArgumentError("nickname %s already exists", input.Nickname).WithMessage("nickname %s already exists", input.Nickname)
	}

	// Emails are also unique in our model so before creating a user we need to check if there's one
	// This way we can differentiate between an invalid argument and an already exists error
	if useCase.emailAlreadyExists(ctx, input.Email) {
		return nil, errors.NewInvalidArgumentError("email %s already exists", input.Email).WithMessage("email %s already exists", input.Email)
	}

	// As the API defines, if no ID is informed, the service will generate a UUID for the resource
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	now := time.Now()
	user := entities.User{
		MetaInformation: entities.MetaInformation{
			ResourceID: entities.ResourceID{
				ID: input.ID,
			},
			Timestamps: entities.Timestamps{
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		UserSpecificData: entities.UserSpecificData{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Country:   input.Country,
			Nickname:  input.Nickname,
		},
	}
	err = user.EncryptAndSetPassword(input.Password)
	if err != nil {
		// TODO: log errors
		return nil, errors.NewInternalError("failed to set user's password").FromError(err)
	}

	storedUser, addErr := useCase.storage.Add(ctx, user)
	if addErr != nil {
		// TODO: improve this error log
		useCase.logger.Error(addErr.Error())
		if addErr.IsAlreadyExistsError() {
			return nil, errors.NewAlreadyExistsError("user %s already exists", user.ID).FromError(err).
				WithMessage("user %s already exists", user.ID)
		}
		return nil, errors.NewInternalError("failed to store user").FromError(err)
	}

	// TODO: uncomment once rabbitmq is injected
	//queueErr := useCase.queue.NotifyCreatedUser(ctx, storedUser.ResourceID)
	//if queueErr != nil {
	//	return nil, queueErr
	//}

	return &CreateUserOutput{
		User: *storedUser,
	}, nil
}

func (useCase DefaultUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.NewInvalidArgumentError("invalid input to describe user").FromError(err).
			WithMessage("Invalid input to describe user. Please review your input fields and try again.")
	}

	storageUser, getErr := useCase.storage.Get(ctx, input.ResourceID)
	if getErr != nil {
		if getErr.IsNotFoundError() {
			return nil, errors.NewNotFoundError("user %s not found", input.ResourceID.ID).FromError(err).
				WithMessage("user %s not found", input.ResourceID.ID)
		}
		return nil, errors.NewInternalError("failed to get user").FromError(err)
	}

	return &GetUserOutput{
		User: *storageUser,
	}, nil
}

func (useCase DefaultUseCase) EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (useCase DefaultUseCase) emailAlreadyExists(ctx context.Context, email string) bool {
	userWithEmail, err := useCase.storage.GetByEmail(ctx, email)
	if err != nil {
		if err.IsNotFoundError() {
			return false
		}
		useCase.logger.Error(fmt.Sprintf("failed to check if email %s already exists", email), err.Error())
	}
	if userWithEmail != nil {
		return true
	}
	return false
}

func (useCase DefaultUseCase) nicknameAlreadyExists(ctx context.Context, nickname string) bool {
	userWithNickname, err := useCase.storage.GetByNickname(ctx, nickname)
	if err != nil {
		if err.IsNotFoundError() {
			return false
		}
		useCase.logger.Error(fmt.Sprintf("failed to check if nickname %s already exists", nickname), err.Error())
	}
	if userWithNickname != nil {
		return true
	}
	return false
}
