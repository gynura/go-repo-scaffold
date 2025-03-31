package user_test

import (
	"github.com/stretchr/testify/require"
	"testing"
	"userservice/internal/core/usecases/user"
	usermocks "userservice/internal/test/helpers/user"
)

func TestUseCase_ProvideDefaultUseCase_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	options := user.ProvideUseCaseOptions{
		Storage: mocks.MockStorage,
		Queue:   mocks.MockQueue,
	}
	useCase, err := user.ProvideDefaultUseCase(options)
	require.NoError(t, err)
	require.NotNil(t, useCase)
}

//func TestUseCase_CreateUser_Success(t *testing.T) {
//	mocks := usermocks.NewUserMocks(t)
//
//}
