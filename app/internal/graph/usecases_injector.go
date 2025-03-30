//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"userservice/internal/adapters/output/queue/healthqueue"
	"userservice/internal/adapters/output/queue/userqueue"
	"userservice/internal/core/usecases/health"
	"userservice/internal/core/usecases/user"
)

type useCasesGraph struct {
	userUseCase   user.UseCase
	HealthUseCase health.UseCase
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

	// Health Use Case
	health.ProvideDefaultUseCase,
	wire.Bind(new(health.UseCase), new(*health.DefaultUseCase)),
	wire.Struct(new(health.ProvideUseCaseOptions), "*"),

	// Health Queue
	healthqueue.ProvideQueueAdapter,
	wire.Bind(new(health.Queue), new(*healthqueue.QueueAdapter)),
	wire.Struct(new(healthqueue.ProvideQueueAdapterOptions), "*"),
)

func initializeUseCases(
	repositories *storageGraph,
	infra *infraGraph,
) (*useCasesGraph, error) {
	wire.Build(useCasesSet,
		wire.FieldsOf(new(*storageGraph),
			"userStorage",
			"healthStorage",
		),
		wire.FieldsOf(new(*infraGraph),
			"publisher",
		),
	)
	return &useCasesGraph{}, nil
}
