package storage

import (
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"
)

func PersistenceToServiceError(err error) *errors.ServiceError {
	if persistence.IsAlreadyExists(err) {
		return errors.NewAlreadyExistsError(err.Error()).FromError(err)
	}
	if persistence.IsNotFound(err) {
		return errors.NewNotFoundError(err.Error()).FromError(err)
	}
	return errors.NewInternalError(err.Error()).FromError(err)
}
