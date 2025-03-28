//go:build wireinject

package graph

import (
	"github.com/google/wire"
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
)

func initializeUseCases(
	repositories *storageGraph,
) (*useCasesGraph, error) {
	wire.Build(useCasesSet,
		wire.FieldsOf(new(*storageGraph),
			"userStorage",
		),
	)
	return &useCasesGraph{}, nil
}
