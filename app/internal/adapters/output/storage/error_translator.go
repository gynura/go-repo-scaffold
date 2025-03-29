package storage

import (
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"
)

func PersistenceToServiceError(err error) *errors.ServiceError {
	if persistence.IsAlreadyExists(err) {
		return errors.NewAlreadyExistsError("resource already exists").FromError(err)
	}
	if persistence.IsNotFound(err) {
		return errors.NewNotFoundError("resource was not found").FromError(err)
	}
	return errors.NewInternalError("internal storage infra error").FromError(err)
}
