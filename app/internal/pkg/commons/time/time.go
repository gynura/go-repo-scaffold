// Package time ensures that time within the application is aligned with the business requirements
package time

import (
	"time"
)

const CustomTimeLayout = "2006-01-02T15:04:05.00Z"

// ServiceTime wraps time.Time to provide the custom requirements
type ServiceTime struct {
	time.Time
}

// Now creates a new ServiceTime with the current UTC time

func Now() ServiceTime {
	return ServiceTime{Time: time.Now().UTC()}
}

// String formats the time in the desired format: "2024-07-15T07:25:55.32Z"
func (ct ServiceTime) String() string {
	return ct.Time.Format(CustomTimeLayout)
}

// FromTime creates a ServiceTime from a time.Time one
func FromTime(t time.Time) (*ServiceTime, error) {
	return &ServiceTime{Time: t.UTC()}, nil
}

// Seconds returns the number of seconds since the ServiceTime instance was created
func (ct ServiceTime) Seconds() int64 {
	return int64(time.Since(ct.Time).Seconds())
}

// TimeSince returns the duration since the ServiceTime instance
func (ct ServiceTime) TimeSince() time.Duration {
	return time.Since(ct.Time)
}
