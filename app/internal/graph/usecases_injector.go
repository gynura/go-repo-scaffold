//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"userservice/internal/adapters/output/queue/userqueue"
	"userservice/internal/core/usecases/user"
)

type useCasesGraph struct {
	UserUseCase user.UseCase
}

var useCasesSet = wire.NewSet(
	wire.Struct(new(useCasesGraph), "*"),

	// User Use Case
	user.ProvideDefaultUseCase,
	wire.Bind(new(user.UseCase), new(*user.DefaultUseCase)),
	wire.Struct(new(user.ProvideUseCaseOptions), "*"),

	// User Queue
	userqueue.ProvideQueueAdapter,
	wire.Bind(new(user.Queue), new(*userqueue.QueueAdapter)),
	wire.Struct(new(userqueue.ProvideQueueAdapterOptions), "*"),
)

func initializeUseCases(
	repositories *storageGraph,
	infra *infraGraph,
) (*useCasesGraph, error) {
	wire.Build(useCasesSet,
		wire.FieldsOf(new(*storageGraph),
			"userStorage",
		),
		wire.FieldsOf(new(*infraGraph),
			"publisher",
		),
	)
	return &useCasesGraph{}, nil
}
