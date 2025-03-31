package user

import (
	"github.com/golang/mock/gomock"
	"testing"
	"userservice/internal/core/usecases/user"
	"userservice/internal/test/generatedmocks"
)

type UserMocks struct {
	*generatedmocks.MockQueue
	*generatedmocks.MockStorage
	*generatedmocks.MockStorageFilters

	controller *gomock.Controller
}

func NewUserMocks(t *testing.T) *UserMocks {
	c := gomock.NewController(t)
	mocks := UserMocks{
		MockQueue:          generatedmocks.NewMockQueue(c),
		MockStorage:        generatedmocks.NewMockStorage(c),
		MockStorageFilters: generatedmocks.NewMockStorageFilters(c),
		controller:         c,
	}

	return &mocks
}

func ProvideTestingUseCase(mocks UserMocks) (user.UseCase, error) {
	return user.ProvideDefaultUseCase(user.ProvideUseCaseOptions{
		Storage: mocks.MockStorage,
		Queue:   mocks.MockQueue,
	})
}
