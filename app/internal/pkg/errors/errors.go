// Package errors defines custom app errors that the domain and application layers must use.
package errors

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// ErrCategory defines error categories for the app
type ErrCategory string

const (
	Internal        ErrCategory = "INTERNAL"
	InvalidArgument ErrCategory = "INVALID_ARGUMENT"
	NotFound        ErrCategory = "NOT_FOUND"
	AlreadyExists   ErrCategory = "ALREADY_EXISTS"
)

// UserServiceError represents an internal structured error type
type UserServiceError struct {
	category   ErrCategory // type of the error
	message    string      // its human-readable message
	wrappedErr error       // wrapped error, can be nil
	errorID    string      // traceable error identifier
}

// Error implements the error interface.
func (e UserServiceError) Error() string {
	if e.wrappedErr != nil {
		return fmt.Sprintf("[%s] (%s) %s: %v", e.category, e.errorID, e.message, e.wrappedErr)
	}
	return fmt.Sprintf("[%s] (%s) %s", e.category, e.errorID, e.message)
}

// Unwrap allows unwrapping the inner error using errors.Is() and errors.As().
func (e UserServiceError) Unwrap() error {
	return e.wrappedErr
}

// Wrap defines a DSL to wrap another error within the user service error
func (e UserServiceError) Wrap(originalErr error) UserServiceError {
	userServiceError := e
	userServiceError.wrappedErr = originalErr
	return userServiceError
}

// WithMessage defines a DSL to add a formatted message to a user service error
func (e UserServiceError) WithMessage(format string, args ...interface{}) UserServiceError {
	userServiceError := e
	userServiceError.message = fmt.Sprintf(format, args...)
	return userServiceError
}

func (e UserServiceError) Category() ErrCategory {
	return e.category
}

func (e UserServiceError) Message() string {
	return e.message
}

func (e UserServiceError) WrappedError() error {
	return e.wrappedErr
}

func (e UserServiceError) ErrorID() string {
	return e.errorID
}

// IsUserServiceError checks if an error is a UserServiceError
func IsUserServiceError(err error) bool {
	var e *UserServiceError
	return errors.As(err, &e)
}

// AsUserServiceError attempts to cast an error to a UserServiceError
func AsUserServiceError(err error) (*UserServiceError, bool) {
	var e *UserServiceError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Constructors for all supported types of UserServiceErrors. This way we enforce a sort of DSL.

func NewInternalError() *UserServiceError {
	return &UserServiceError{
		category: Internal,
		errorID:  uuid.New().String(),
	}
}

func NewNotFoundError() *UserServiceError {
	return &UserServiceError{
		category: NotFound,
		errorID:  uuid.New().String(),
	}
}

func NewAlreadyExistsError() *UserServiceError {
	return &UserServiceError{
		category: AlreadyExists,
		errorID:  uuid.New().String(),
	}
}

func NewInvalidArgumentError() *UserServiceError {
	return &UserServiceError{
		category: InvalidArgument,
		errorID:  uuid.New().String(),
	}
}
