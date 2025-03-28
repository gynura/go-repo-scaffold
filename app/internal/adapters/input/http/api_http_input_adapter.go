package http

import (
	"context"
	"errors"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/http/generatedhttpinfra"
)

// DefaultAPIAdapter implements the API Handler interface used by the server
type DefaultAPIAdapter struct {
	userUseCase user.UseCase
}

// DefaultAPIAdapterOptions defines the options to implement the DefaultAPIAdapter
type DefaultAPIAdapterOptions struct {
	UserUseCase user.UseCase
}

// ProvideDefaultAPIAdapter provides the default implementation of the REST API Handler
func ProvideDefaultAPIAdapter(options DefaultAPIAdapterOptions) (*DefaultAPIAdapter, error) {
	if options.UserUseCase == nil {
		return nil, errors.New("userUseCase is required")
	}

	return &DefaultAPIAdapter{
		userUseCase: options.UserUseCase,
	}, nil
}

var _ generatedhttpinfra.Handler = new(DefaultAPIAdapter)

func (d DefaultAPIAdapter) UserUsersCreate(ctx context.Context, req generatedhttpinfra.OptCreateUser) (generatedhttpinfra.UserUsersCreateRes, error) {
	input := user.CreateUserInput{
		ResourceID: entities.ResourceID{
			ID: req.Value.ID.Value,
		},
		FirstName: req.Value.SpecificData.FirstName,
		LastName:  req.Value.SpecificData.LastName,
		Email:     req.Value.SpecificData.Email,
		Password:  req.Value.SpecificData.Password,
		Country:   req.Value.SpecificData.Country,
		Nickname:  req.Value.SpecificData.Nickname,
	}
	createdUser, err := d.userUseCase.CreateUser(ctx, input)
	if err != nil {
		if err.IsInternalError() {
			response := generatedhttpinfra.UserUsersCreateInternalServerError(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.UserUsersCreateBadRequest(mapServiceErrorToHTTPError(*err))
		return &response, nil
	}

	return &generatedhttpinfra.DescribeUser{
		MetaInformation: generatedhttpinfra.MetaInformation{
			ID:        createdUser.ID,
			CreatedAt: createdUser.CreatedAt.Time,
			UpdatedAt: createdUser.UpdatedAt.Time,
		},
		SpecificData: generatedhttpinfra.UserSpecificData{
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Nickname:  createdUser.Nickname,
			Password:  createdUser.Password(),
			Email:     createdUser.Email,
			Country:   createdUser.Country,
		},
	}, nil
}

func (d DefaultAPIAdapter) UserUsersDescribe(ctx context.Context, params generatedhttpinfra.UserUsersDescribeParams) (generatedhttpinfra.UserUsersDescribeRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultAPIAdapter) UserUsersEdit(ctx context.Context, req generatedhttpinfra.OptUpdateUser, params generatedhttpinfra.UserUsersEditParams) (generatedhttpinfra.UserUsersEditRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultAPIAdapter) UserUsersList(ctx context.Context, params generatedhttpinfra.UserUsersListParams) (generatedhttpinfra.UserUsersListRes, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultAPIAdapter) UserUsersRemove(ctx context.Context, params generatedhttpinfra.UserUsersRemoveParams) (generatedhttpinfra.UserUsersRemoveRes, error) {
	//TODO implement me
	panic("implement me")
}
