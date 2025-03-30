package health

import "userservice/internal/pkg/commons/time"

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type Status string

// CheckOutput contains the health information of the service
type CheckOutput struct {
	Status    Status
	StartedAt time.ServiceTime
	Checks    ServiceChecks
}

// ServiceChecks defines the status of the services that are not the server hosting the business logic
type ServiceChecks struct {
	Queue    Status
	Database Status
}
