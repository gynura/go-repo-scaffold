// Package time ensures that time within the application is aligned with the business requirements
package time

import (
	"encoding/json"
	"time"
)

const CustomTimeLayout = "2006-01-02T15:04:05.00Z"

// UserServiceTime wraps time.Time to provide the custom requirements
type UserServiceTime struct {
	time.Time
}

// Now creates a new UserServiceTime with the current UTC time
func Now() UserServiceTime {
	return UserServiceTime{Time: time.Now().UTC()}
}

// String formats the time in the desired format: "2024-07-15T07:25:55.32Z"
func (ct UserServiceTime) String() string {
	return ct.Time.Format(CustomTimeLayout)
}

// MarshalJSON defines the custom JSON serializer
func (ct UserServiceTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}
