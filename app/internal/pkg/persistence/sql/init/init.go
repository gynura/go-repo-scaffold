package init

import (
	// This package must import all dialects inside sql and must be imported in order to have their init functions run
	_ "userservice/internal/pkg/persistence/sql/postgres"
)
