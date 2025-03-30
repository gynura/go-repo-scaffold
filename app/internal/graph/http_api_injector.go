//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"userservice/internal/adapters/input/http"
	"userservice/internal/infra/http/generatedhttpinfra"
)

type httpAPIGraph struct {
	HTTPHandler generatedhttpinfra.Handler
}

var httpAPISet = wire.NewSet(
	wire.Struct(new(httpAPIGraph), "*"),

	// HTTP Input Adapter
	http.ProvideDefaultAPIAdapter,
	wire.Bind(new(generatedhttpinfra.Handler), new(*http.DefaultAPIAdapter)),
	wire.Struct(new(http.DefaultAPIAdapterOptions), "*"),
)

func initializeHTTPAPIGraph(
	useCases *useCasesGraph,
) (*httpAPIGraph, error) {
	wire.Build(httpAPISet,
		wire.FieldsOf(new(*useCasesGraph),
			"userUseCase",
			"HealthUseCase",
		),
	)
	return &httpAPIGraph{}, nil
}
