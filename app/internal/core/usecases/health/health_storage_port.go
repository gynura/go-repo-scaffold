package health

// Storage defines the functionality that the healthqueue use case needs to do liveness checks on the plugged queue
type Storage interface {
	// GetStorageStatus returns the status of the plugged queue
	GetStorageStatus() Status
}
