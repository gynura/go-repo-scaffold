//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"userservice/internal/adapters/input/http"
	"userservice/internal/infra/http/generatedhttpinfra"
)

type httpGraph struct {
	HTTPInputAdapter generatedhttpinfra.Handler
}

var httpSet = wire.NewSet(
	wire.Struct(new(httpGraph), "*"),

	// HTTP Input Adapter
	http.ProvideDefaultAPIAdapter,
	wire.Bind(new(generatedhttpinfra.Handler), new(*http.DefaultAPIAdapter)),
	wire.Struct(new(http.DefaultAPIAdapterOptions), "*"),
)

func initializeHTTPAdapters(
	useCases *useCasesGraph,
) (*httpGraph, error) {
	wire.Build(httpSet,
		wire.FieldsOf(new(*useCasesGraph),
			"userUseCase",
		),
	)
	return &httpGraph{}, nil
}
