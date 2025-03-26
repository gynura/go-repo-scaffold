package errors_test

import (
	"errors"
	"fmt"
	"testing"
	userserviceerrors "userservice/internal/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceError_Error(t *testing.T) {
	err := errors.New("wrapped error")
	appErr := userserviceerrors.NewInternalError("foo").
		FromError(err).
		WithMessage("Something went wrong")

	expectedPrefix := fmt.Sprintf("[%s] (%s)", userserviceerrors.Internal, appErr.ErrorID())
	assert.Contains(t, appErr.Error(), expectedPrefix)
	assert.Contains(t, appErr.Error(), "Something went wrong")
	assert.Contains(t, appErr.Error(), "wrapped error")
}

func TestUserServiceError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := userserviceerrors.NewNotFoundError("foo").FromError(originalErr)

	assert.Equal(t, originalErr, appErr.Unwrap())
}

func TestUserServiceError_Wrap(t *testing.T) {
	originalErr := errors.New("database connection failed")
	appErr := userserviceerrors.NewInternalError("foo").FromError(originalErr)

	assert.Equal(t, originalErr, appErr.Unwrap())
	assert.Contains(t, appErr.Error(), "database connection failed")
}

func TestUserServiceError_Message(t *testing.T) {
	appErr := userserviceerrors.NewInvalidArgumentError("foo").WithMessage("Field %s is required", "username")

	assert.Contains(t, appErr.Error(), "Field username is required")
}

func TestIsUserServiceError(t *testing.T) {
	appErr := userserviceerrors.NewAlreadyExistsError("foo")
	assert.True(t, userserviceerrors.IsUserServiceError(appErr))
	assert.False(t, userserviceerrors.IsUserServiceError(errors.New("some random error")))
}

func TestAsUserServiceError(t *testing.T) {
	appErr := userserviceerrors.NewNotFoundError("foo")
	genericErr := errors.New("some error")

	extractedErr, ok := userserviceerrors.AsUserServiceError(appErr)
	assert.True(t, ok)
	assert.Equal(t, appErr, extractedErr)

	extractedErr, ok = userserviceerrors.AsUserServiceError(genericErr)
	assert.False(t, ok)
	assert.Nil(t, extractedErr)
}

func TestUserServiceError_ErrorCategoriesConstructors(t *testing.T) {
	internalErr := userserviceerrors.NewInternalError("foo")
	notFoundErr := userserviceerrors.NewNotFoundError("foo")
	alreadyExistsErr := userserviceerrors.NewAlreadyExistsError("foo")
	invalidArgErr := userserviceerrors.NewInvalidArgumentError("foo")

	assert.Equal(t, userserviceerrors.Internal, internalErr.Category())
	assert.Equal(t, userserviceerrors.NotFound, notFoundErr.Category())
	assert.Equal(t, userserviceerrors.AlreadyExists, alreadyExistsErr.Category())
	assert.Equal(t, userserviceerrors.InvalidArgument, invalidArgErr.Category())

	assert.NotEmpty(t, internalErr.ErrorID())
	assert.NotEmpty(t, notFoundErr.ErrorID())
	assert.NotEmpty(t, alreadyExistsErr.ErrorID())
	assert.NotEmpty(t, invalidArgErr.ErrorID())
}

func TestUserServiceError_ErrorID(t *testing.T) {
	appErr := userserviceerrors.NewInternalError("foo")
	parsedUUID, err := uuid.Parse(appErr.ErrorID())

	assert.NoError(t, err)
	assert.Equal(t, appErr.ErrorID(), parsedUUID.String())
}
