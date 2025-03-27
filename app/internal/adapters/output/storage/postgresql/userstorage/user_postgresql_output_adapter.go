// Package userstorage defines the output database adapters for the User resource
package userstorage

import (
	"context"
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

func (s StorageAdapter) Add(ctx context.Context, data entities.User) (*entities.User, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (s StorageAdapter) Get(ctx context.Context, id entities.ResourceID) (*entities.User, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
}

func (s StorageAdapter) Remove(ctx context.Context, id entities.ResourceID) *errors.ServiceError {
	//TODO implement me
	panic("implement me")
}

func (s StorageAdapter) List(filters user.StorageFilters) (*user.UsersList, *errors.ServiceError) {
	//TODO implement me
	panic("implement me")
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
