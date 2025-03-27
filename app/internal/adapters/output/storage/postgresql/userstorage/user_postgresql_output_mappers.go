package userstorage

import (
	"github.com/asaskevich/govalidator"
	"os/user"
	"userservice/internal/core/entities"
	"userservice/internal/infra/storage/postgresql/userstorageinfra"
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"
)

func mapToCreateDB(user entities.User) (*userstorageinfra.CreateInput, error) {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return nil, errors.NewInternalError("invalid User format").FromError(err)
	}

	db := userstorageinfra.CreateInput{
		UserDB: userstorageinfra.UserDB{
			ResourceID: user.ResourceID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Email:      user.Email,
			Password:   user.Password(),
			Country:    user.Country,
			Nickname:   user.Nickname,
			CreatedAt:  user.CreatedAt.String(),
			UpdatedAt:  user.UpdatedAt.String(),
		},
	}
	return &db, nil
}

func mapToUpdateDB(user user.User) (*userstorageinfra.EditInput, error) {
	// TODO: complete
}

func mapFromDB(db userstorageinfra.UserDB) (*user.User, error) {
	// TODO: complete
}

func mapSliceFromDB(dbResources []userstorageinfra.UserDB) ([]user.User, error) {
	userSlice := make([]user.User, len(dbResources))
	for index := range dbResources {
		item, err := mapFromDB(dbResources[index])
		if err != nil {
			return nil, err
		}
		userSlice[index] = *item
	}

	return userSlice, nil
}

func persistenceToServiceError(err error) error {
	if persistence.IsAlreadyExists(err) {
		return errors.NewAlreadyExistsError(err.Error()).FromError(err)
	}
	if persistence.IsNotFound(err) {
		return errors.NewNotFoundError(err.Error()).FromError(err)
	}
	return errors.NewInternalError(err.Error()).FromError(err)
}
