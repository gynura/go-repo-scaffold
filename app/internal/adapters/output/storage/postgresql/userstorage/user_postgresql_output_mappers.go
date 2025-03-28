package userstorage

import (
	"userservice/internal/core/entities"
	"userservice/internal/infra/storage/postgresql/userstorageinfra"
	"userservice/internal/pkg/commons/time"
	"userservice/internal/pkg/errors"

	"github.com/asaskevich/govalidator"
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

func mapToUpdateDB(user entities.User) (*userstorageinfra.EditInput, error) {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return nil, errors.NewInternalError("invalid User format").FromError(err)
	}

	db := userstorageinfra.EditInput{
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

func mapFromDB(storedUser userstorageinfra.UserDB) (*entities.User, error) {
	_, err := govalidator.ValidateStruct(storedUser)
	if err != nil {
		return nil, errors.NewInternalError("invalid stored User format").FromError(err)
	}

	createdAt, err := time.FromString(storedUser.CreatedAt)
	if err != nil {
		return nil, errors.NewInternalError("invalid stored User createdAt").FromError(err)
	}

	updatedAt, err := time.FromString(storedUser.UpdatedAt)
	if err != nil {
		return nil, errors.NewInternalError("invalid stored User updatedAt").FromError(err)
	}

	user := entities.User{
		MetaInformation: entities.MetaInformation{
			ResourceID: entities.ResourceID{
				ID: storedUser.ID,
			},
			Timestamps: entities.Timestamps{
				CreatedAt: *createdAt,
				UpdatedAt: *updatedAt,
			},
		},
		UserSpecificData: entities.UserSpecificData{
			FirstName: storedUser.FirstName,
			LastName:  storedUser.LastName,
			Email:     storedUser.Email,
			Country:   storedUser.Country,
			Nickname:  storedUser.Nickname,
		},
	}

	// There's no need to encrypt the password as the password is ALWAYS stored already encrypted from the useCase logic
	err = user.SetPassword(storedUser.Password)
	if err != nil {
		return nil, errors.NewInternalError("invalid stored User password").FromError(err)
	}

	return &user, nil
}

func mapSliceFromDB(dbResources []userstorageinfra.UserDB) ([]entities.User, error) {
	userSlice := make([]entities.User, len(dbResources))
	for index := range dbResources {
		item, err := mapFromDB(dbResources[index])
		if err != nil {
			return nil, err
		}
		userSlice[index] = *item
	}

	return userSlice, nil
}
