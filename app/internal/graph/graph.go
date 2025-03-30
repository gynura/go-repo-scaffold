// Package graph defines the dependency management of the service
package graph

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userservice/internal/infra/http/generatedhttpinfra"
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
	graph.storageGraph, err = initializeStorage(graph.librariesGraph.persistenceFramework)
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

	server := &http.Server{
		Addr:         options.Host + ":" + options.Port,
		Handler:      infraServerHandler,
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

	go func() {
		log.Printf("Starting server on :%s", options.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}
