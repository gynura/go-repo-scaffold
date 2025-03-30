package healthqueue

import (
	"userservice/internal/core/usecases/health"
	"userservice/internal/infra/queue"
	"userservice/internal/pkg/errors"
)

// QueueAdapter defines the default implementation of the health.Queue interface
type QueueAdapter struct {
	publisher queue.Publisher
}

// ProvideQueueAdapterOptions defines the options to configure a UserUserQueueAdapter
type ProvideQueueAdapterOptions struct {
	Publisher *queue.Publisher
}

// ProvideQueueAdapter creates a UserUserQueueAdapter with the given options
func ProvideQueueAdapter(options ProvideQueueAdapterOptions) (*QueueAdapter, error) {
	if options.Publisher == nil {
		return nil, errors.NewInternalError("no publisher provided")
	}
	return &QueueAdapter{
		publisher: *options.Publisher,
	}, nil
}

var _ health.Queue = new(QueueAdapter)

func (q QueueAdapter) GetQueueStatus() health.Status {
	if q.publisher.IsAlive() {
		return health.StatusUp
	}
	return health.StatusDown
}
