package http

import (
	"net/http"
	"userservice/internal/infra/http/generatedhttpinfra"
	"userservice/internal/pkg/errors"
)

const (
	httpCodeBadRequest    = 400
	httpCodeNotFound      = 404
	httpCodeInternalError = 500
	httpCodeConflict      = 409

	defaultInternalErrorMessage = "internal server error"
)

// errorCategoryToErrorStatus maps service error categories to HTTP error statuses
var errorCategoryToErrorStatus = map[errors.ErrCategory]generatedhttpinfra.ErrorErrorStatus{
	errors.ErrCategoryNotFound:        generatedhttpinfra.ErrorErrorStatusNOTFOUND,
	errors.ErrCategoryAlreadyExists:   generatedhttpinfra.ErrorErrorStatusALREADYEXISTS,
	errors.ErrCategoryInvalidArgument: generatedhttpinfra.ErrorErrorStatusINVALIDARGUMENT,
	errors.ErrCategoryInternal:        generatedhttpinfra.ErrorErrorStatusINTERNAL,
}

// errorStatusToHTTPCode maps service error categories to their HTTP response codes
var errorCategoryToHTTPCode = map[errors.ErrCategory]int{
	errors.ErrCategoryNotFound:        httpCodeNotFound,
	errors.ErrCategoryAlreadyExists:   httpCodeConflict,
	errors.ErrCategoryInvalidArgument: httpCodeBadRequest,
	errors.ErrCategoryInternal:        httpCodeInternalError,
}

// mapServiceErrorToHTTPError converts a service error to a standardized HTTP error response
func mapServiceErrorToHTTPError(err errors.ServiceError) generatedhttpinfra.Error {
	httpCode, ok := errorCategoryToHTTPCode[err.Category()]
	if !ok {
		httpCode = http.StatusInternalServerError
	}

	status, ok := errorCategoryToErrorStatus[err.Category()]
	if !ok {
		status = generatedhttpinfra.ErrorErrorStatusINTERNAL
	}

	// Construct error response
	return generatedhttpinfra.Error{
		Error: generatedhttpinfra.ErrorError{
			Code:        httpCode,
			Status:      status,
			Description: selectErrorDescription(err, httpCode),
			ErrorId:     err.ErrorID(),
		},
	}
}

// selectErrorDescription chooses appropriate error message based on error type
func selectErrorDescription(err errors.ServiceError, httpCode int) string {
	// We expose specific error message for client-side errors
	if httpCode < httpCodeInternalError {
		return err.ExposedMessage()
	}

	// But details must be hidden for internal errors always
	return defaultInternalErrorMessage
}
