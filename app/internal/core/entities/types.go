// Package entities defines custom data types ready to be used across the app
package entities

import "time"

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	OrderByUpdatedAt = "updatedAt"
	OrderByCreatedAt = "createdAt"
)

// ID defines the standard format for unique identifiers
type ID struct {
	ID string `json:"id"`
}

type Timestamps struct {
	// CreatedAt defines the date time of when the resource was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt defines the date time of when the resource was last edited using the app
	UpdatedAt time.Time `json:"updatedAt"`
}

// MetaInformation defines the meta information of all app's resources
type MetaInformation struct {
	ID         `json:",inline"`
	Timestamps `json:",inline"`
}
