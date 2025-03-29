package persistence_test

import (
	"testing"
	"userservice/internal/pkg/persistence"

	"github.com/stretchr/testify/require"
	embedded "userservice/internal"
)

func TestStorageConfig_GetStatement_Success(t *testing.T) {
	storageConfigOptions := persistence.StorageConfigOptions{
		ReadDirAndFileFS: embedded.DatabaseMappers,
		MappersPath:      "",
		Driver:           "postgresql",
	}
	storageConfig, err := persistence.NewStorageConfig(storageConfigOptions)
	require.Nil(t, err)

	_, ok := storageConfig.GetStatement("service.user.insert")
	require.True(t, ok)

	_, ok = storageConfig.GetStatement("service.user.list")
	require.True(t, ok)
}

func TestStorageConfig_GetStatement_Errors(t *testing.T) {
	storageConfigOptions := persistence.StorageConfigOptions{
		ReadDirAndFileFS: embedded.DatabaseMappers,
		MappersPath:      "",
		Driver:           "postgresql",
	}
	storageConfig, err := persistence.NewStorageConfig(storageConfigOptions)
	require.Nil(t, err)

	_, ok := storageConfig.GetStatement("this.doesnt.exist")
	require.False(t, ok)

	storageConfigOptions = persistence.StorageConfigOptions{
		ReadDirAndFileFS: embedded.DatabaseMappers,
		MappersPath:      "foo",
		Driver:           "postgresql",
	}
	storageConfig, err = persistence.NewStorageConfig(storageConfigOptions)
	require.Error(t, err)
}
