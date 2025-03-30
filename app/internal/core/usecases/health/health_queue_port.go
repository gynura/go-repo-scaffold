package health

// Queue defines functionality that the healthqueue use case needs to do liveness checks on the plugged queue
type Queue interface {
	// GetQueueStatus returns the status of the plugged queue
	GetQueueStatus() Status
}
