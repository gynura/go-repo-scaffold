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
	logger    *logger.Logger
	startedAt time.ServiceTime
}

// ProvideUseCaseOptions defines the options to configure a UseCase implementation
type ProvideUseCaseOptions struct {
}

// ProvideDefaultUseCase provides the default implementation of the UseCase logic
func ProvideDefaultUseCase(options ProvideUseCaseOptions) (*DefaultUseCase, error) {
	return &DefaultUseCase{
		logger: logger.NewLogger(),
	}, nil
}

var _ UseCase = new(DefaultUseCase)

func (d DefaultUseCase) Check(_ context.Context) (*CheckOutput, *errors.ServiceError) {
	return &CheckOutput{
		Status:    StatusUp,
		StartedAt: d.startedAt,
		Checks: ServiceChecks{
			Queue:    StatusUp,
			Database: StatusDown,
		},
	}, nil
}

func (d DefaultUseCase) SetStartedAt(_ context.Context) {
	d.startedAt = time.Now()
}
