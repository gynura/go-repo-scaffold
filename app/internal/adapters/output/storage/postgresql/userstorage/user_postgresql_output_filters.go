package userstorage

import (
	"userservice/internal/adapters/output/storage/postgresql"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/storage/postgresql/userstorageinfra"
	"userservice/internal/pkg/persistence"
)

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
