// Package userstorage defines the output database adapters for the User resource
package userstorage

import (
	"context"
	"userservice/internal/adapters/output/storage"
	"userservice/internal/adapters/output/storage/postgresql"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/storage/postgresql/userstorageinfra"
	"userservice/internal/pkg/errors"
	"userservice/internal/pkg/persistence"
)

var _ user.Storage = new(StorageAdapter)

// StorageAdapter implementation of user.UserStorage
type StorageAdapter struct {
	infra *userstorageinfra.StorageInfra
}

// RepositoryOptions configures a StorageAdapter
type RepositoryOptions struct {
	Infra *userstorageinfra.StorageInfra
}

// NewRepository creates a StorageAdapter with the given options
func NewRepository(options RepositoryOptions) (*StorageAdapter, error) {
	return &StorageAdapter{
		infra: options.Infra,
	}, nil
}

func (adapter StorageAdapter) Add(ctx context.Context, data entities.User) (*entities.User, *errors.ServiceError) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.NewInternalError(err.Error()).FromError(err)
	}

	storageData, err := adapter.infra.Add(ctx, *db)
	if err != nil {
		return nil, storage.PersistenceToServiceError(err)
	}

	addedUser, err := mapFromDB(*storageData)
	if err != nil {
		return nil, errors.NewInternalError(err.Error()).FromError(err)
	}

	return addedUser, nil
}

func (adapter StorageAdapter) Get(ctx context.Context, id entities.ResourceID) (*entities.User, *errors.ServiceError) {
	storageData, err := adapter.infra.Get(ctx, id)
	if err != nil {
		return nil, storage.PersistenceToServiceError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NewNotFoundError("no User found with id: %s", id)
	}

	if len(storageData) > 1 {
		return nil, errors.NewInternalError("more than one User found with id: %s", id)
	}

	storedUser, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	return storedUser, nil
}

func (adapter StorageAdapter) Edit(ctx context.Context, data entities.User) (*entities.User, *errors.ServiceError) {
	db, err := mapToUpdateDB(data)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	result, err := adapter.infra.Edit(ctx, *db)
	if err != nil {
		return nil, storage.PersistenceToServiceError(err)
	}

	rowsAffected, errRowsAffected := result.Result.RowsAffected()
	if errRowsAffected != nil {
		return nil, errors.NewInternalError(errRowsAffected.Error())
	}

	if rowsAffected == 0 {
		return nil, errors.NewNotFoundError("no User found with id: %s", data.ResourceID)
	}

	if rowsAffected > 1 {
		return nil, errors.NewInternalError("more than one User found with id: %s", data.ResourceID)
	}

	storageData, err := adapter.Get(ctx, data.ResourceID)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}

	return storageData, nil
}

func (adapter StorageAdapter) Remove(ctx context.Context, id entities.ResourceID) *errors.ServiceError {
	_, err := adapter.infra.Exists(ctx, id)
	if err != nil {
		return storage.PersistenceToServiceError(err)
	}

	_, err = adapter.infra.Remove(ctx, id)
	if err != nil {
		return storage.PersistenceToServiceError(err)
	}

	return nil
}

func (adapter StorageAdapter) List(ctx context.Context, filters user.StorageFilters) (*user.UsersList, *errors.ServiceError) {
	f, ok := filters.(*storageFilter)
	if !ok {
		return nil, errors.NewInternalError("invalid filters")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := adapter.infra.List(ctx, *f.UserDBFilter)
	if err != nil {
		return nil, storage.PersistenceToServiceError(err)
	}

	collection := user.UsersList{}
	if f.Pagination != nil {
		collection.Offset = f.Pagination.Offset
		collection.Limit = f.Pagination.Limit - 1
		if len(storageData) == f.Pagination.Limit {
			collection.HasMoreItems = true
			storageData = storageData[:len(storageData)-1]
		}
		f.Pagination.Limit--
	} else {
		collection.Pagination = entities.Pagination{
			Limit:        len(storageData),
			Offset:       0,
			HasMoreItems: false,
		}
	}

	items, err := mapSliceFromDB(storageData)
	if err != nil {
		return nil, errors.NewInternalError(err.Error())
	}
	collection.Items = items

	return &collection, nil
}

type storageFilter struct {
	*userstorageinfra.UserDBFilter
}

var _ user.StorageFilters = (*storageFilter)(nil)

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
func (filter *storageFilter) Paged(limit int, offset int) user.StorageFilters {
	filter.UserDBFilter = filter.UserDBFilter.Paged(limit, offset)
	return filter
}

func (filter *storageFilter) OrderByCreatedAt(direction persistence.OrderDirection) user.StorageFilters {
	filter.UserDBFilter = filter.UserDBFilter.Sort("created_at", direction)
	return filter
}

func (filter *storageFilter) OrderByUpdatedAt(direction persistence.OrderDirection) user.StorageFilters {
	filter.UserDBFilter = filter.UserDBFilter.Sort("updated_at", direction)
	return filter
}

func (filter *storageFilter) FilterByCountry(country string) user.StorageFilters {
	filter.Country = country
	filter.AppendFilter(postgresql.NewEqualFilter("country"))
	return filter
}
