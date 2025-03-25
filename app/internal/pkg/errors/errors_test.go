package errors_test

import (
	"errors"
	"fmt"
	"testing"
	userserviceerrors "userservice/internal/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := errors.New("wrapped error")
	appErr := userserviceerrors.NewInternalError().WithMessage("Something went wrong").Wrap(err)

	expectedPrefix := fmt.Sprintf("[%s] (%s)", userserviceerrors.Internal, appErr.ErrorID())
	assert.Contains(t, appErr.Error(), expectedPrefix)
	assert.Contains(t, appErr.Error(), "Something went wrong")
	assert.Contains(t, appErr.Error(), "wrapped error")
}

func TestAppError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := userserviceerrors.NewNotFoundError().Wrap(originalErr)

	assert.Equal(t, originalErr, appErr.Unwrap())
}

func TestAppError_Wrap(t *testing.T) {
	originalErr := errors.New("database connection failed")
	appErr := userserviceerrors.NewInternalError().Wrap(originalErr)

	assert.Equal(t, originalErr, appErr.Unwrap())
	assert.Contains(t, appErr.Error(), "database connection failed")
}

func TestAppError_Message(t *testing.T) {
	appErr := userserviceerrors.NewInvalidArgumentError().WithMessage("Field %s is required", "username")

	assert.Contains(t, appErr.Error(), "Field username is required")
}

func TestIsAppError(t *testing.T) {
	appErr := userserviceerrors.NewAlreadyExistsError()
	assert.True(t, userserviceerrors.IsUserServiceError(appErr))
	assert.False(t, userserviceerrors.IsUserServiceError(errors.New("some random error")))
}

func TestAsAppError(t *testing.T) {
	appErr := userserviceerrors.NewNotFoundError()
	genericErr := errors.New("some error")

	extractedErr, ok := userserviceerrors.AsUserServiceError(appErr)
	assert.True(t, ok)
	assert.Equal(t, appErr, extractedErr)

	extractedErr, ok = userserviceerrors.AsUserServiceError(genericErr)
	assert.False(t, ok)
	assert.Nil(t, extractedErr)
}

func TestNewErrorFactories(t *testing.T) {
	internalErr := userserviceerrors.NewInternalError()
	notFoundErr := userserviceerrors.NewNotFoundError()
	alreadyExistsErr := userserviceerrors.NewAlreadyExistsError()
	invalidArgErr := userserviceerrors.NewInvalidArgumentError()

	assert.Equal(t, userserviceerrors.Internal, internalErr.Category())
	assert.Equal(t, userserviceerrors.NotFound, notFoundErr.Category())
	assert.Equal(t, userserviceerrors.AlreadyExists, alreadyExistsErr.Category())
	assert.Equal(t, userserviceerrors.InvalidArgument, invalidArgErr.Category())

	assert.NotEmpty(t, internalErr.ErrorID())
	assert.NotEmpty(t, notFoundErr.ErrorID())
	assert.NotEmpty(t, alreadyExistsErr.ErrorID())
	assert.NotEmpty(t, invalidArgErr.ErrorID())
}

func TestAppError_ErrorID(t *testing.T) {
	appErr := userserviceerrors.NewInternalError()
	parsedUUID, err := uuid.Parse(appErr.ErrorID())

	assert.NoError(t, err)
	assert.Equal(t, appErr.ErrorID(), parsedUUID.String())
}
