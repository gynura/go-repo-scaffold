// Package healthstorage defines the output database adapters for the Health resource
package healthstorage

import (
	"database/sql"
	"errors"
	"userservice/internal/core/usecases/health"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// StorageAdapter implementation of health.Storage
type StorageAdapter struct {
	postgresqlConnectionString string
}

// StorageAdapterOptions configures a StorageAdapter
type StorageAdapterOptions struct {
	PostgresqlConnectionString string
}

// ProvideStorageAdapter creates a StorageAdapter with the given options
func ProvideStorageAdapter(options StorageAdapterOptions) (*StorageAdapter, error) {
	if len(options.PostgresqlConnectionString) == 0 {
		return nil, errors.New("postgresql connection string is required")
	}
	return &StorageAdapter{
		postgresqlConnectionString: options.PostgresqlConnectionString,
	}, nil
}

var _ health.Storage = new(StorageAdapter)

func (adapter StorageAdapter) GetStorageStatus() health.Status {
	if adapter.isDBAlive() {
		return health.StatusUp
	}
	return health.StatusDown
}

// TODO: dumbed down check, replace with something that actually tests the db (?)
func (adapter StorageAdapter) isDBAlive() bool {
	db, err := sql.Open("postgres", adapter.postgresqlConnectionString)
	if err != nil {
		return false
	}
	defer db.Close()

	return true
}
