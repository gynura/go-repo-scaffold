package health

import (
	"context"
	"userservice/internal/pkg/commons/logger"
	"userservice/internal/pkg/commons/time"
	"userservice/internal/pkg/errors"
)

// UseCase defines the functionality to manage User resources
type UseCase interface {
	// Check checks whether the service and its controlled dependencies are healthy or not
	Check(ctx context.Context) (*CheckOutput, *errors.ServiceError)
	// SetStartedAt sets the time at which the server started
	SetStartedAt(ctx context.Context)
}

// DefaultUseCase defines the default implementation of the health.UseCase interface
type DefaultUseCase struct {
	queue     Queue
	storage   Storage
	logger    *logger.Logger
	startedAt time.ServiceTime
}

// ProvideUseCaseOptions defines the options to configure a UseCase implementation
type ProvideUseCaseOptions struct {
	Queue   Queue
	Storage Storage
}

// ProvideDefaultUseCase provides the default implementation of the UseCase logic
func ProvideDefaultUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	if options.Queue == nil {
		return nil, errors.NewInternalError("queue is required")
	}

	return &DefaultUseCase{
		logger:  logger.NewLogger(),
		queue:   options.Queue,
		storage: options.Storage,
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (d *DefaultUseCase) Check(_ context.Context) (*CheckOutput, *errors.ServiceError) {
	output := &CheckOutput{
		Status:    StatusUp,
		StartedAt: d.startedAt,
		Checks: ServiceChecks{
			Queue:    d.queue.GetQueueStatus(),
			Database: d.storage.GetStorageStatus(),
		},
	}

	if output.Checks.Database == StatusDown || output.Checks.Queue == StatusDown {
		output.Status = StatusDown
	}

	return output, nil
}

func (d *DefaultUseCase) SetStartedAt(_ context.Context) {
	now := time.Now()
	d.startedAt = now
}
