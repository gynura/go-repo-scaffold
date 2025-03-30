//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"userservice/internal/adapters/output/storage/postgresql/userstorage"
	"userservice/internal/infra/storage/postgresql/userstorageinfra"
	"userservice/internal/pkg/persistence"

	"userservice/internal/core/usecases/user"
)

type storageGraph struct {
	userStorage user.Storage
}

var storageSet = wire.NewSet(
	wire.Struct(new(storageGraph), "*"),

	// User Database Infra
	userstorageinfra.ProvideStorageInfra,
	wire.Struct(new(userstorageinfra.StorageInfraOptions), "*"),

	// User Storage Output Adapter
	userstorage.ProvideStorageAdapter,
	wire.Bind(new(user.Storage), new(*userstorage.StorageAdapter)),
	wire.Struct(new(userstorage.StorageAdapterOptions), "*"),
)

func initializeStorage(
	persistenceFramework persistence.Storage,
) (*storageGraph, error) {
	wire.Build(storageSet)
	return &storageGraph{}, nil
}
