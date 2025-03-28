// Package graph defines the dependency management of the service
package graph

import (
	"userservice/internal/pkg/commons/validators"
	"userservice/internal/pkg/persistence/sql"

	"github.com/asaskevich/govalidator"
)

const (
	defaultAPIPort = "8080"
)

type ServiceGraph struct {
	config Config

	useCasesGraph  *useCasesGraph
	storageGraph   *storageGraph
	librariesGraph *librariesGraph
	httpGraph      *httpGraph
}

func New(config Config) (*ServiceGraph, error) {
	valid, err := govalidator.ValidateStruct(config)
	if err != nil || !valid {
		return nil, err
	}

	return &ServiceGraph{
		config: config,
	}, nil
}

// Build builds the application graph
func (graph *ServiceGraph) Build() {
	var err error

	// Validators
	validators.SetValidators()

	// Common libs
	graph.librariesGraph, err = initializeLibraries(graph.config)
	if err != nil {
		panic(err)
	}

	// Storage layer
	graph.storageGraph, err = InitializeStorage(graph.librariesGraph.persistenceFramework)
	if err != nil {
		panic(err)
	}

	// Business logic layer
	graph.useCasesGraph, err = initializeUseCases(graph.storageGraph)
	if err != nil {
		panic(err)
	}

	// HTTP Adapters
	graph.httpGraph, err = initializeHTTPAdapters(graph.useCasesGraph)
	if err != nil {
		panic(err)
	}
}

func (graph *ServiceGraph) PersistenceFwConnection() sql.Connection {
	return graph.librariesGraph.persistenceConnection
}

type StartServerOptions struct {
	Port string
}

func (graph *ServiceGraph) StartServer(options StartServerOptions) {
	// TODO: complete
}
