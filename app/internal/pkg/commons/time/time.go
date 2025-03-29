// Package time ensures that time within the application is aligned with the business requirements
package time

import (
	"encoding/json"
	"fmt"
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

// FromString creates a ServiceTime from a string or an error if the format of the string is invalid
func FromString(newTime string) (*ServiceTime, error) {
	t, err := time.Parse(CustomTimeLayout, newTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time format %s, expected YYYY-MM-DDTHH:MM:SS.SSZ", newTime)
	}
	return &ServiceTime{Time: t}, nil
}

// FromTime creates a ServiceTime from a time.Time one
func FromTime(t time.Time) (*ServiceTime, error) {
	return &ServiceTime{Time: t.UTC()}, nil
}
