package user

import (
	"context"
	"fmt"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/commons/logger"
	"userservice/internal/pkg/commons/time"
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	emailAlreadyExistsMessage    string = "Email %s is already registered"
	nicknameAlreadyExistsMessage string = "Nickname %s is not available"
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
	logger  *logger.Logger
	queue   Queue
}

// ProvideUseCaseOptions defines the options to configure a UserUseCase implementation
type ProvideUseCaseOptions struct {
	Storage Storage
	Queue   Queue
}

// ProvideDefaultUseCase provides the default implementation of the User use case logic
func ProvideDefaultUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.NewInvalidArgumentError("no storage provided")
	}
	if options.Queue == nil {
		return nil, errors.NewInvalidArgumentError("no queue provided")
	}

	return &DefaultUseCase{
		storage: options.Storage,
		queue:   options.Queue,
		logger:  logger.NewLogger(),
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (useCase DefaultUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		useCase.logger.Debug(ctx).Message(err.Error())
		return nil, errors.NewInvalidArgumentError("invalid input to create user").FromError(err).
			WithMessage("Invalid input to create user. Please review your input fields and try again.")
	}

	useCase.logger.Info(ctx).Attributes("trying to create new user", map[string]any{
		"userId": input.ID,
	})

	// Nicknames are unique alongside IDs and Emails, in order to avoid coupling this check to the plugged db, the business logic needs to implement this checks
	if useCase.nicknameAlreadyExists(ctx, input.Nickname) {
		useCase.logger.Debug(ctx).Message(fmt.Sprintf(nicknameAlreadyExistsMessage, input.Nickname))
		return nil, errors.NewInvalidArgumentError(nicknameAlreadyExistsMessage, input.Nickname).WithMessage(nicknameAlreadyExistsMessage, input.Nickname)
	}

	// Emails are also unique in our model so before creating a user we need to check if there's one
	// This way we can differentiate between an invalid argument and an already exists error
	if useCase.emailAlreadyExists(ctx, input.Email) {
		useCase.logger.Debug(ctx).Message(fmt.Sprintf(emailAlreadyExistsMessage, input.Email))
		return nil, errors.NewInvalidArgumentError(emailAlreadyExistsMessage, input.Email).WithMessage(emailAlreadyExistsMessage, input.Email)
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
		useCase.logger.Error(ctx).Attributes("error encrypting and setting user's password", map[string]any{
			"userId": input.ID,
			"info":   err.Error(),
		})
		return nil, errors.NewInternalError("error encrypting and setting user's password").FromError(err)
	}

	storedUser, addErr := useCase.storage.Add(ctx, user)
	if addErr != nil {
		useCase.logger.Error(ctx).UserServiceError(addErr, user.ID)
		if addErr.IsAlreadyExistsError() {
			return nil, errors.NewAlreadyExistsError("user %s already exists", user.ID).FromError(addErr).
				WithMessage("user %s already exists", user.ID)
		}
		return nil, errors.NewInternalError("failed to store user").FromError(addErr)
	}

	useCase.logger.Info(ctx).Attributes("created new user", map[string]any{
		"userId": input.ID,
	})

	queueErr := useCase.queue.NotifyCreatedUser(ctx, storedUser.ResourceID)
	if queueErr != nil {
		useCase.logger.Error(ctx).UserServiceError(queueErr, user.ID)
		return nil, errors.NewInternalError("error notifying created user with id: ", user.ID).FromError(queueErr)
	}

	return &CreateUserOutput{
		User: *storedUser,
	}, nil
}

func (useCase DefaultUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		useCase.logger.Debug(ctx).Message(err.Error())
		return nil, errors.NewInvalidArgumentError("invalid input to list users").FromError(err).
			WithMessage("Invalid input to list users. Please review your input fields and try again.")
	}

	filters := useCase.storage.ProvideFilters()
	direction := defaultOrderDirection
	if len(input.OrderDirection) > 0 {
		direction, err = entities.ToOrderDirection(input.OrderDirection)
		if err != nil {
			useCase.logger.Debug(ctx).Message("unsupported direction [%s] to order", input.OrderDirection)
			return nil, errors.NewInvalidArgumentError("unsupported direction [%s] to order", input.OrderDirection).
				WithMessage("unsupported direction [%s] to order", input.OrderDirection)
		}
	}

	// The way filtering works is that, if the attribute is defined in the input struct, that means we must query it
	if len(input.Country) > 0 {
		filters.FilterByCountry(input.Country)
	}

	filters.OrderByCreatedAt(persistence.OrderDirection(direction))
	if input.OrderBy == string(entities.OrderByUpdatedAt) {
		filters.OrderByUpdatedAt(persistence.OrderDirection(direction))
	}

	pageLimit := defaultLimit
	if input.PageLimit > 0 {
		pageLimit = input.PageLimit
		filters.Paged(pageLimit, input.PageOffset)
	}

	userList, listErr := useCase.storage.List(ctx, filters)
	if listErr != nil {
		useCase.logger.Error(ctx).ServiceError(listErr)
		return nil, errors.NewInternalError("failed to list users").FromError(listErr)
	}

	return &ListUsersOutput{
		UsersList: *userList,
	}, nil
}

func (useCase DefaultUseCase) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		useCase.logger.Debug(ctx).Message(err.Error())
		return nil, errors.NewInvalidArgumentError("invalid input to describe user").FromError(err).
			WithMessage("Invalid input to describe user. Please review your input fields and try again.")
	}

	storageUser, getErr := useCase.storage.Get(ctx, input.ResourceID)
	if getErr != nil {
		useCase.logger.Error(ctx).UserServiceError(getErr, input.ID)
		if getErr.IsNotFoundError() {
			return nil, errors.NewNotFoundError("user %s not found", input.ResourceID.ID).FromError(getErr).
				WithMessage("user %s not found", input.ResourceID.ID)
		}
		return nil, errors.NewInternalError("failed to get user").FromError(getErr)
	}

	return &GetUserOutput{
		User: *storageUser,
	}, nil
}

func (useCase DefaultUseCase) EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		useCase.logger.Debug(ctx).Message(err.Error())
		return nil, errors.NewInvalidArgumentError("invalid input to edit user").FromError(err).
			WithMessage("Invalid input to edit user. Please review your input fields and try again.")
	}

	useCase.logger.Info(ctx).Attributes("trying to update user information", map[string]any{
		"userId": input.ID,
	})

	if useCase.nicknameAlreadyExists(ctx, input.Nickname) {
		useCase.logger.Debug(ctx).Message(fmt.Sprintf(nicknameAlreadyExistsMessage, input.Nickname))
		return nil, errors.NewInvalidArgumentError(nicknameAlreadyExistsMessage, input.Nickname).WithMessage(nicknameAlreadyExistsMessage, input.Nickname)
	}

	if useCase.emailAlreadyExists(ctx, input.Email) {
		useCase.logger.Debug(ctx).Message(fmt.Sprintf(emailAlreadyExistsMessage, input.Email))
		return nil, errors.NewInvalidArgumentError(emailAlreadyExistsMessage, input.Email).WithMessage(emailAlreadyExistsMessage, input.Email)
	}

	// When editing a resource we must first fetch the state on the storage of that resource prior to the editing (source of truth)
	userToBeUpdated, getErr := useCase.storage.Get(ctx, input.ResourceID)
	if getErr != nil {
		useCase.logger.Error(ctx).UserServiceError(getErr, input.ID)
		if getErr.IsNotFoundError() {
			return nil, errors.NewNotFoundError("user %s not found", input.ResourceID.ID).FromError(getErr).
				WithMessage("user %s not found", input.ResourceID.ID)
		}
		return nil, errors.NewInternalError("failed to get user").FromError(getErr)
	}

	// Then we will edit that user with the editable fields
	userToBeUpdated.MetaInformation.Timestamps.UpdatedAt = time.Now()
	userToBeUpdated.Nickname = input.Nickname
	userToBeUpdated.Email = input.Email
	userToBeUpdated.FirstName = input.FirstName
	userToBeUpdated.LastName = input.LastName
	userToBeUpdated.Country = input.Country
	err = userToBeUpdated.EncryptAndSetPassword(input.Password)
	if err != nil {
		useCase.logger.Error(ctx).Attributes("error encrypting and setting user's password", map[string]any{
			"userId": input.ID,
			"info":   err.Error(),
		})
		return nil, errors.NewInternalError("failed to set user's password").FromError(err)
	}

	storedUser, editErr := useCase.storage.Edit(ctx, *userToBeUpdated)
	if editErr != nil {
		useCase.logger.Error(ctx).UserServiceError(editErr, input.ID)
		if editErr.IsNotFoundError() {
			return nil, errors.NewNotFoundError("user %s not found", input.ResourceID.ID).FromError(editErr).
				WithMessage("user %s not found", input.ResourceID.ID)
		}
		return nil, errors.NewInternalError("failed to edit user").FromError(editErr)
	}

	useCase.logger.Info(ctx).Attributes("successfully updated user information", map[string]any{
		"userId": input.ID,
	})

	queueErr := useCase.queue.NotifyEditUser(ctx, storedUser.ResourceID)
	if queueErr != nil {
		useCase.logger.Error(ctx).UserServiceError(queueErr, storedUser.ID)
		return nil, errors.NewInternalError("error notifying edited user with id: ", storedUser.ID).FromError(queueErr)
	}

	return &EditUserOutput{
		User: *storedUser,
	}, nil

}

func (useCase DefaultUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, *errors.ServiceError) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		useCase.logger.Debug(ctx).Message(err.Error())
		return nil, errors.NewInvalidArgumentError("invalid input to delete user").FromError(err).
			WithMessage("Invalid input to delete user. Please review your input fields and try again.")
	}

	useCase.logger.Info(ctx).Attributes("deleting user", map[string]any{
		"userId": input.ID,
	})

	removedUser, removeErr := useCase.storage.Remove(ctx, input.ResourceID)
	if removeErr != nil {
		useCase.logger.Error(ctx).UserServiceError(removeErr, input.ID)
		if removeErr.IsNotFoundError() {
			return nil, errors.NewNotFoundError("user %s not found", input.ResourceID.ID).FromError(err).
				WithMessage("user %s not found", input.ResourceID.ID)
		}
		return nil, errors.NewInternalError("failed to delete user").FromError(err)
	}

	useCase.logger.Info(ctx).Attributes("successfully deleted user", map[string]any{
		"userId": input.ID,
	})

	queueErr := useCase.queue.NotifyRemoveUser(ctx, removedUser.ResourceID)
	if queueErr != nil {
		useCase.logger.Error(ctx).UserServiceError(queueErr, removedUser.ID)
		return nil, errors.NewInternalError("error notifying removed user with id: ", removedUser.ID).FromError(queueErr)
	}

	return &DeleteUserOutput{
		User: *removedUser,
	}, nil
}

func (useCase DefaultUseCase) emailAlreadyExists(ctx context.Context, email string) bool {
	userWithEmail, err := useCase.storage.GetByEmail(ctx, email)
	if err != nil {
		if err.IsNotFoundError() {
			return false
		}
		useCase.logger.Error(ctx).ServiceError(err)
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
		useCase.logger.Error(ctx).ServiceError(err)
	}
	if userWithNickname != nil {
		return true
	}
	return false
}
