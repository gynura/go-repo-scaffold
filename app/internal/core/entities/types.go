// Package entities defines custom data types ready to be used across the app
package entities

import "time"

// OrderBy defines the resource's attribute used to order when listing
type OrderBy string

// OrderDirection defines the order direction when listing
type OrderDirection string

const (
	OrderAsc  OrderDirection = "asc"
	OrderDesc OrderDirection = "desc"

	OrderByUpdatedAt OrderBy = "updatedAt"
	OrderByCreatedAt OrderBy = "createdAt"
)

// ID defines the standard format for unique identifiers
type ID struct {
	ID string `json:"id"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt"` // date time of when the resource was created
	UpdatedAt time.Time `json:"updatedAt"` // date time of when the resource was last edited using the app
}

// MetaInformation defines the meta information of all app's resources
type MetaInformation struct {
	ID         `json:",inline"`
	Timestamps `json:",inline"`
}

// Pagination defines the pagination of all app's resources lists
type Pagination struct {
	Limit        int  `json:"limit"`        // number of items per page in the collection
	Offset       int  `json:"offset"`       // starting point (index) of the collection
	HasMoreItems bool `json:"hasMoreItems"` // indicates if there are additional pages of data available
}
