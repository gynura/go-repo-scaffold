package http

import (
	"context"
	"errors"
	"userservice/internal/core/entities"
	"userservice/internal/core/usecases/health"
	"userservice/internal/core/usecases/user"
	"userservice/internal/infra/http/generatedhttpinfra"
)

// DefaultAPIAdapter implements the API Handler interface used by the server
type DefaultAPIAdapter struct {
	userUseCase   user.UseCase
	healthUseCase health.UseCase
}

// DefaultAPIAdapterOptions defines the options to implement the DefaultAPIAdapter
type DefaultAPIAdapterOptions struct {
	UserUseCase   user.UseCase
	HealthUseCase health.UseCase
}

// ProvideDefaultAPIAdapter provides the default implementation of the REST API Handler
func ProvideDefaultAPIAdapter(options DefaultAPIAdapterOptions) (*DefaultAPIAdapter, error) {
	if options.UserUseCase == nil {
		return nil, errors.New("userUseCase is required")
	}
	if options.HealthUseCase == nil {
		return nil, errors.New("healthUseCase is required")
	}

	return &DefaultAPIAdapter{
		userUseCase:   options.UserUseCase,
		healthUseCase: options.HealthUseCase,
	}, nil
}

var _ generatedhttpinfra.Handler = new(DefaultAPIAdapter)

func (d DefaultAPIAdapter) ServiceUsersCreate(ctx context.Context, req generatedhttpinfra.OptCreateUser) (generatedhttpinfra.ServiceUsersCreateRes, error) {
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
		if err.IsInvalidArgumentError() {
			response := generatedhttpinfra.ServiceUsersCreateBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		if err.IsAlreadyExistsError() {
			response := generatedhttpinfra.ServiceUsersCreateBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.ServiceUsersCreateInternalServerError(mapServiceErrorToHTTPError(*err))
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

func (d DefaultAPIAdapter) ServiceUsersDescribe(ctx context.Context, params generatedhttpinfra.ServiceUsersDescribeParams) (generatedhttpinfra.ServiceUsersDescribeRes, error) {
	input := user.GetUserInput{
		ResourceID: entities.ResourceID{
			ID: params.UserId,
		},
	}
	describedUser, err := d.userUseCase.GetUser(ctx, input)
	if err != nil {
		if err.IsInvalidArgumentError() {
			response := generatedhttpinfra.ServiceUsersDescribeBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		if err.IsNotFoundError() {
			response := generatedhttpinfra.ServiceUsersDescribeNotFound(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.ServiceUsersDescribeInternalServerError(mapServiceErrorToHTTPError(*err))
		return &response, nil
	}

	return &generatedhttpinfra.DescribeUser{
		MetaInformation: generatedhttpinfra.MetaInformation{
			ID:        describedUser.ID,
			CreatedAt: describedUser.CreatedAt.Time,
			UpdatedAt: describedUser.UpdatedAt.Time,
		},
		SpecificData: generatedhttpinfra.UserSpecificData{
			FirstName: describedUser.FirstName,
			LastName:  describedUser.LastName,
			Nickname:  describedUser.Nickname,
			Password:  describedUser.Password(),
			Email:     describedUser.Email,
			Country:   describedUser.Country,
		},
	}, nil
}

func (d DefaultAPIAdapter) ServiceUsersEdit(ctx context.Context, req generatedhttpinfra.OptUpdateUser, params generatedhttpinfra.ServiceUsersEditParams) (generatedhttpinfra.ServiceUsersEditRes, error) {
	input := user.EditUserInput{
		ResourceID: entities.ResourceID{ID: params.UserId},
		FirstName:  req.Value.SpecificData.FirstName,
		LastName:   req.Value.SpecificData.LastName,
		Email:      req.Value.SpecificData.Email,
		Password:   req.Value.SpecificData.Password,
		Country:    req.Value.SpecificData.Country,
		Nickname:   req.Value.SpecificData.Nickname,
	}

	editedUser, err := d.userUseCase.EditUser(ctx, input)
	if err != nil {
		if err.IsInvalidArgumentError() {
			response := generatedhttpinfra.ServiceUsersEditBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		if err.IsNotFoundError() {
			response := generatedhttpinfra.ServiceUsersEditNotFound(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.ServiceUsersEditInternalServerError(mapServiceErrorToHTTPError(*err))
		return &response, nil
	}

	return &generatedhttpinfra.DescribeUser{
		MetaInformation: generatedhttpinfra.MetaInformation{
			ID:        editedUser.ID,
			CreatedAt: editedUser.CreatedAt.Time,
			UpdatedAt: editedUser.UpdatedAt.Time,
		},
		SpecificData: generatedhttpinfra.UserSpecificData{
			FirstName: editedUser.FirstName,
			LastName:  editedUser.LastName,
			Nickname:  editedUser.Nickname,
			Password:  editedUser.Password(),
			Email:     editedUser.Email,
			Country:   editedUser.Country,
		},
	}, nil
}

func (d DefaultAPIAdapter) ServiceUsersList(ctx context.Context, params generatedhttpinfra.ServiceUsersListParams) (generatedhttpinfra.ServiceUsersListRes, error) {
	input := user.ListUsersInput{
		PageLimit:      params.Limit.Value,
		PageOffset:     params.Offset.Value,
		OrderBy:        string(params.OrderBy.Value),
		OrderDirection: string(params.OrderDirection.Value),
		Country:        params.Country.Value,
	}

	usersList, err := d.userUseCase.ListUsers(ctx, input)
	if err != nil {
		if err.IsInvalidArgumentError() {
			response := generatedhttpinfra.ServiceUsersListBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.ServiceUsersListInternalServerError(mapServiceErrorToHTTPError(*err))
		return &response, nil
	}

	return &generatedhttpinfra.UserList{
		Items:        adaptUserList(usersList.Items),
		Limit:        usersList.Limit,
		Offset:       usersList.Offset,
		HasMoreItems: usersList.HasMoreItems,
	}, nil
}

func (d DefaultAPIAdapter) ServiceUsersRemove(ctx context.Context, params generatedhttpinfra.ServiceUsersRemoveParams) (generatedhttpinfra.ServiceUsersRemoveRes, error) {
	input := user.DeleteUserInput{
		ResourceID: entities.ResourceID{ID: params.UserId},
	}

	deletedUser, err := d.userUseCase.DeleteUser(ctx, input)
	if err != nil {
		if err.IsNotFoundError() {
			response := generatedhttpinfra.ServiceUsersRemoveNotFound(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		if err.IsInvalidArgumentError() {
			response := generatedhttpinfra.ServiceUsersRemoveBadRequest(mapServiceErrorToHTTPError(*err))
			return &response, nil
		}
		response := generatedhttpinfra.ServiceUsersRemoveInternalServerError(mapServiceErrorToHTTPError(*err))
		return &response, nil
	}

	return &generatedhttpinfra.DescribeUser{
		MetaInformation: generatedhttpinfra.MetaInformation{
			ID:        deletedUser.ID,
			CreatedAt: deletedUser.CreatedAt.Time,
			UpdatedAt: deletedUser.UpdatedAt.Time,
		},
		SpecificData: generatedhttpinfra.UserSpecificData{
			FirstName: deletedUser.FirstName,
			LastName:  deletedUser.LastName,
			Nickname:  deletedUser.Nickname,
			Password:  deletedUser.Password(),
			Email:     deletedUser.Email,
			Country:   deletedUser.Country,
		},
	}, nil
}

func (d DefaultAPIAdapter) ServiceHealthLive(ctx context.Context) (generatedhttpinfra.ServiceHealthLiveRes, error) {
	output, err := d.healthUseCase.Check(ctx)
	if err != nil {
		response := mapServiceErrorToHTTPError(*err)
		return &response, nil
	}

	return &generatedhttpinfra.Health{
		Status:    adaptStatus(output.Status),
		StartedAt: output.StartedAt.Time,
		Checks: generatedhttpinfra.HealthChecks{
			Queue:    adaptQueueStatus(output.Checks.Queue),
			Database: adaptDatabaseStatus(output.Checks.Database),
		},
	}, nil
}

func adaptStatus(status health.Status) generatedhttpinfra.HealthStatus {
	if status == health.StatusUp {
		return generatedhttpinfra.HealthStatusUP
	}
	return generatedhttpinfra.HealthStatusDOWN
}

func adaptQueueStatus(status health.Status) generatedhttpinfra.HealthChecksQueue {
	if status == health.StatusUp {
		return generatedhttpinfra.HealthChecksQueueUP
	}
	return generatedhttpinfra.HealthChecksQueueDOWN
}

func adaptDatabaseStatus(status health.Status) generatedhttpinfra.HealthChecksDatabase {
	if status == health.StatusUp {
		return generatedhttpinfra.HealthChecksDatabaseUP
	}
	return generatedhttpinfra.HealthChecksDatabaseDOWN
}

func adaptUserList(items []entities.User) []generatedhttpinfra.DescribeUser {
	var adaptedList []generatedhttpinfra.DescribeUser
	for _, item := range items {
		adaptedUser := generatedhttpinfra.DescribeUser{
			MetaInformation: generatedhttpinfra.MetaInformation{
				ID:        item.ID,
				CreatedAt: item.CreatedAt.Time,
				UpdatedAt: item.UpdatedAt.Time,
			},
			SpecificData: generatedhttpinfra.UserSpecificData{
				FirstName: item.FirstName,
				LastName:  item.LastName,
				Nickname:  item.Nickname,
				Password:  item.Password(),
				Email:     item.Email,
				Country:   item.Country,
			},
		}
		adaptedList = append(adaptedList, adaptedUser)
	}

	return adaptedList
}
