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

func TestUseCase_CreateUser_QueueError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockQueue.EXPECT().NotifyCreatedUser(gomock.Any(), gomock.Any()).Return(internalErr)

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

func TestUseCase_ListUsers_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	usersList := &user.UsersList{
		Items: []entities.User{exampleUser},
	}

	mocks.MockStorage.EXPECT().ProvideFilters().Return(mocks.MockStorageFilters)
	mocks.MockStorage.EXPECT().List(gomock.Any(), gomock.Any()).Return(usersList, nil)
	mocks.MockStorageFilters.EXPECT().FilterByCountry(gomock.Any()).Return(mocks.MockStorageFilters)
	mocks.MockStorageFilters.EXPECT().OrderByCreatedAt(gomock.Any()).Return(mocks.MockStorageFilters)
	mocks.MockStorageFilters.EXPECT().Paged(gomock.Any(), gomock.Any()).Return(mocks.MockStorageFilters)

	input := user.ListUsersInput{
		Country:        "ESP",
		OrderBy:        string(entities.OrderByCreatedAt),
		OrderDirection: string(entities.OrderDesc),
		PageLimit:      10,
		PageOffset:     0,
	}
	output, serviceErr := useCase.ListUsers(ctx, input)
	require.Nil(t, serviceErr)
	require.NotNil(t, output)
	require.Equal(t, *usersList, output.UsersList)
}

func TestUseCase_ListUsers_StorageError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().ProvideFilters().Return(mocks.MockStorageFilters)
	mocks.MockStorage.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, internalErr)
	mocks.MockStorageFilters.EXPECT().OrderByCreatedAt(gomock.Any()).Return(mocks.MockStorageFilters)

	input := user.ListUsersInput{}
	output, serviceErr := useCase.ListUsers(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}

func TestUseCase_GetUser_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)

	input := user.GetUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
	}
	output, serviceErr := useCase.GetUser(ctx, input)
	require.Nil(t, serviceErr)
	require.NotNil(t, output)
	require.Equal(t, exampleUser, output.User)
}

func TestUseCase_GetUser_NotFound_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)

	input := user.GetUserInput{
		ResourceID: entities.ResourceID{
			ID: "nonexistent",
		},
	}
	output, serviceErr := useCase.GetUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsNotFoundError())
	require.Nil(t, output)
}

func TestUseCase_GetUser_StorageError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, internalErr)

	input := user.GetUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
	}
	output, serviceErr := useCase.GetUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}

func TestUseCase_EditUser_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	userToEdit := exampleUser
	updatedUser := exampleUser
	updatedUser.FirstName = "updated"
	updatedUser.LastName = "person"
	updatedUser.Email = "updated@gmail.com"
	updatedUser.Country = "USA"
	updatedUser.Nickname = "updatedname"

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&userToEdit, nil)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(&updatedUser, nil)
	mocks.MockQueue.EXPECT().NotifyEditUser(gomock.Any(), gomock.Any()).Return(nil)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "updated@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "updatedname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.Nil(t, serviceErr)
	require.NotNil(t, output)
	require.Equal(t, updatedUser, output.User)
}

func TestUseCase_EditUser_NotFound_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "nonexistent",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "updated@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "updatedname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsNotFoundError())
	require.Nil(t, output)
}

func TestUseCase_EditUser_EmailAlreadyTaken_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	userToEdit := exampleUser
	otherUser := exampleUser
	otherUser.ID = "user2"
	otherUser.Email = "taken@gmail.com"

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&userToEdit, nil)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&otherUser, nil)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "taken@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "updatedname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInvalidArgumentError())
	require.Nil(t, output)
}

func TestUseCase_EditUser_NicknameAlreadyTaken_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "updated@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "takennickname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInvalidArgumentError())
	require.Nil(t, output)
}

func TestUseCase_EditUser_StorageError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	userToEdit := exampleUser

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&userToEdit, nil)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(nil, internalErr)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "updated@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "updatedname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}

func TestUseCase_EditUser_QueueError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	userToEdit := exampleUser
	updatedUser := exampleUser
	updatedUser.FirstName = "updated"

	mocks.MockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&userToEdit, nil)
	mocks.MockStorage.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().GetByNickname(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
	mocks.MockStorage.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(&updatedUser, nil)
	mocks.MockQueue.EXPECT().NotifyEditUser(gomock.Any(), gomock.Any()).Return(internalErr)

	input := user.EditUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
		FirstName: "updated",
		LastName:  "person",
		Email:     "updated@gmail.com",
		Password:  "newpassword123",
		Country:   "USA",
		Nickname:  "updatedname",
	}
	output, serviceErr := useCase.EditUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}

func TestUseCase_DeleteUser_Success(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)
	mocks.MockQueue.EXPECT().NotifyRemoveUser(gomock.Any(), gomock.Any()).Return(nil)

	input := user.DeleteUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
	}
	output, serviceErr := useCase.DeleteUser(ctx, input)
	require.Nil(t, serviceErr)
	require.NotNil(t, output)
	require.Equal(t, exampleUser, output.User)
}

func TestUseCase_DeleteUser_NotFound_Error(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)

	input := user.DeleteUserInput{
		ResourceID: entities.ResourceID{
			ID: "nonexistent",
		},
	}
	output, serviceErr := useCase.DeleteUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsNotFoundError())
	require.Nil(t, output)
}

func TestUseCase_DeleteUser_StorageError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil, internalErr)

	input := user.DeleteUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
	}
	output, serviceErr := useCase.DeleteUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}

func TestUseCase_DeleteUser_QueueError(t *testing.T) {
	mocks := usermocks.NewUserMocks(t)
	useCase, err := usermocks.ProvideTestingUseCase(*mocks)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mocks.MockStorage.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(&exampleUser, nil)
	mocks.MockQueue.EXPECT().NotifyRemoveUser(gomock.Any(), gomock.Any()).Return(internalErr)

	input := user.DeleteUserInput{
		ResourceID: entities.ResourceID{
			ID: "user1",
		},
	}
	output, serviceErr := useCase.DeleteUser(ctx, input)
	require.NotNil(t, serviceErr)
	require.True(t, serviceErr.IsInternalError())
	require.Nil(t, output)
}
