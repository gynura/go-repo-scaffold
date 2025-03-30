// Package graph defines the dependency management of the service
package graph

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userservice/internal/infra/http/generatedhttpinfra"
	"userservice/internal/infra/http/middlewares"
	"userservice/internal/pkg/commons/logger"
	"userservice/internal/pkg/commons/validators"
	_ "userservice/internal/pkg/persistence/sql/init" // Imported to run the supported sql dialects registry

	"github.com/asaskevich/govalidator"
)

type ServiceGraph struct {
	config Config

	useCasesGraph  *useCasesGraph
	storageGraph   *storageGraph
	librariesGraph *librariesGraph
	httpAPIGraph   *httpAPIGraph
	infraGraph     *infraGraph
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

	// Infra layer
	graph.infraGraph, err = initializeInfraGraph(graph.config)
	if err != nil {
		panic(err)
	}

	// Storage layer
	graph.storageGraph, err = initializeStorage(graph.librariesGraph.persistenceFramework, graph.config.GetPostgresqlConnectionString())
	if err != nil {
		panic(err)
	}

	// Business logic layer
	graph.useCasesGraph, err = initializeUseCases(graph.storageGraph, graph.infraGraph)
	if err != nil {
		panic(err)
	}

	// HTTP API Adapter
	graph.httpAPIGraph, err = initializeHTTPAPIGraph(graph.useCasesGraph)
	if err != nil {
		panic(err)
	}
}

// StartServerOptions defines configuration for starting the HTTP server
type StartServerOptions struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}

func (graph *ServiceGraph) StartServer(options StartServerOptions) error {
	_, err := govalidator.ValidateStruct(options)
	if err != nil {
		return fmt.Errorf("error validating start server options: %w", err)
	}

	infraServerHandler, err := generatedhttpinfra.NewServer(graph.httpAPIGraph.HTTPHandler)
	if err != nil {
		return fmt.Errorf("failed to create http infra handler server: %w", err)
	}

	loggingMiddleware := middlewares.LoggingMiddleware(infraServerHandler)

	server := &http.Server{
		Addr:         options.Host + ":" + options.Port,
		Handler:      loggingMiddleware,
		ReadTimeout:  15 * time.Second, // TODO: make them configurable (?)
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	serviceLogger := logger.NewLogger()

	// Before starting the server we need to set the startAt time for the healthqueue use case
	graph.useCasesGraph.HealthUseCase.SetStartedAt(ctx)

	go func() {
		serviceLogger.Info(ctx).Message("Starting server on %s:%s", options.Host, options.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serviceLogger.Info(ctx).Message("HTTP server ListenAndServe error: %v", err)
		}
	}()

	// Very simple graceful shutdown approach, once the server is terminated (context cancellation) we
	// try to give a 5 seconds timeout allowing ongoing connections to finish within that time frame
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: make graceful timeout configurable
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	serviceLogger.Info(ctx).Message("Server gracefully stopped")
	return nil
}
