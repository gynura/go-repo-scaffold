// Package time ensures that time within the application is aligned with the business requirements
package time

import (
	"encoding/json"
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

// MarshalJSON defines the custom JSON serializer
func (ct ServiceTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}
