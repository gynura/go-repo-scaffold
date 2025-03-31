package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/pkg/commons/time"
	"userservice/internal/pkg/errors"
	usermocks "userservice/internal/test/helpers/user"
)

var (
	ctx         = context.Background()
	exampleUser = entities.User{
		MetaInformation: entities.MetaInformation{
			ResourceID: entities.ResourceID{
				ID: "user1",
			},
			Timestamps: entities.Timestamps{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		UserSpecificData: entities.UserSpecificData{
			FirstName: "fooman",
			LastName:  "fooman",
			Email:     "fooman@gmail.com",
			Country:   "fooman",
			Nickname:  "fooman",
		},
	}
	notFoundErr = errors.NewNotFoundError("test error")
	internalErr = errors.NewInternalError("test error")
)

func TestUseCase_ProvideDefaultUseCase_All(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	options := user.ProvideUseCaseOptions{
		Storage: mocks.MockStorage,
		Queue:   mocks.MockQueue,
	}
	useCase, err := user.ProvideDefaultUseCase(options)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	options = user.ProvideUseCaseOptions{
		Storage: mocks.MockStorage,
	}
	useCase, err = user.ProvideDefaultUseCase(options)
	require.Error(t, err)
	require.Nil(t, useCase)

	options = user.ProvideUseCaseOptions{
		Queue: mocks.MockQueue,
	}
	useCase, err = user.ProvideDefaultUseCase(options)
	require.Error(t, err)
	require.Nil(t, useCase)
}

func TestUseCase_CreateUser_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockQueue.EXPECT().NotifyCreatedUser(gomock.Any(), gomock.Any()).Return(nil)

	input := user.CreateUserInput{
		ResourceID: entities.ResourceID{
			ID: "theOne",
		},
		FirstName: "bar",
		LastName:  "foo",
		Email:     "bar@gmail.com",
		Password:  "supersecretpassword",
		Country:   "ESP",
		Nickname:  "supernickname",
	}
	output, serviceErr := useCase.CreateUser(ctx, input)
	require.Nil(t, serviceErr)
	require.NotNil(t, output)
}

func TestUseCase_CreateUser_EmailAlreadyTaken_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)

	input := user.CreateUserInput{
		ResourceID: entities.ResourceID{
			ID: "theOne",
		},
		FirstName: "bar",
		LastName:  "foo",
		Email:     "bar@gmail.com",
		Password:  "supersecretpassword",
		Country:   "ESP",
		Nickname:  "supernickname",
	}
	output, serviceErr := useCase.CreateUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInvalidArgumentError())
	require.Nil(t, output)
}

func TestUseCase_CreateUser_NicknameAlreadyTaken_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)

	input := user.CreateUserInput{
		ResourceID: entities.ResourceID{
			ID: "theOne",
		},
		FirstName: "bar",
		LastName:  "foo",
		Email:     "bar@gmail.com",
		Password:  "supersecretpassword",
		Country:   "ESP",
		Nickname:  "supernickname",
	}
	output, serviceErr := useCase.CreateUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInvalidArgumentError())
	require.Nil(t, output)
}

func TestUseCase_CreateUser_StorageError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil, internalErr)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)

	input := user.CreateUserInput{
		ResourceID: entities.ResourceID{
			ID: "theOne",
		},
		FirstName: "bar",
		LastName:  "foo",
		Email:     "bar@gmail.com",
		Password:  "supersecretpassword",
		Country:   "ESP",
		Nickname:  "supernickname",
	}
	output, serviceErr := useCase.CreateUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}
