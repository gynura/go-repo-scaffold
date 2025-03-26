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

// ServiceError represents an internal structured error type
type ServiceError struct {
	category    ErrCategory // type of the error
	internalMsg string      // internal error message
	exposedMsg  string      // its human-readable message, the one that the API will typically expose
	innerErr    error       // wrapped error, can be nil
	errorID     string      // traceable error identifier
	// TODO: extend with a catalogueID that could be exposed as part of the doc
}

// Error implements the error interface.
func (e ServiceError) Error() string {
	if e.innerErr != nil {
		return fmt.Sprintf("[%s] (%s) %s: %v", e.category, e.errorID, e.exposedMsg, e.innerErr)
	}
	return fmt.Sprintf("[%s] (%s) %s", e.category, e.errorID, e.exposedMsg)
}

// Unwrap allows unwrapping the inner error using errors.Is() and errors.As().
func (e ServiceError) Unwrap() error {
	return e.innerErr
}

// FromError defines a DSL to wrap another error within the user service error
func (e ServiceError) FromError(originalErr error) ServiceError {
	userServiceError := e
	userServiceError.innerErr = originalErr
	return userServiceError
}

// WithMessage defines a DSL to add a formatted message that will be exposed to a user service error
func (e ServiceError) WithMessage(format string, args ...interface{}) ServiceError {
	userServiceError := e
	userServiceError.exposedMsg = fmt.Sprintf(format, args...)
	return userServiceError
}

func (e ServiceError) Category() ErrCategory {
	return e.category
}

func (e ServiceError) ExposedMessage() string {
	return e.exposedMsg
}

func (e ServiceError) InnerError() error {
	return e.innerErr
}

func (e ServiceError) ErrorID() string {
	return e.errorID
}

// IsUserServiceError checks if an error is a ServiceError
func IsUserServiceError(err error) bool {
	var e *ServiceError
	return errors.As(err, &e)
}

// AsServiceError attempts to cast an error to a ServiceError
func AsServiceError(err error) (*ServiceError, bool) {
	var e *ServiceError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Constructors for all supported types of UserServiceErrors. This way we enforce a sort of DSL.

func NewInternalError(format string, args ...interface{}) *ServiceError {
	return &ServiceError{
		category:    Internal,
		errorID:     uuid.New().String(),
		internalMsg: fmt.Sprintf(format, args...),
	}
}

func NewNotFoundError(format string, args ...interface{}) *ServiceError {
	return &ServiceError{
		category:    NotFound,
		errorID:     uuid.New().String(),
		internalMsg: fmt.Sprintf(format, args...),
	}
}

func NewAlreadyExistsError(format string, args ...interface{}) *ServiceError {
	return &ServiceError{
		category:    AlreadyExists,
		errorID:     uuid.New().String(),
		internalMsg: fmt.Sprintf(format, args...),
	}
}

func NewInvalidArgumentError(format string, args ...interface{}) *ServiceError {
	return &ServiceError{
		category:    InvalidArgument,
		errorID:     uuid.New().String(),
		internalMsg: fmt.Sprintf(format, args...),
	}
}
