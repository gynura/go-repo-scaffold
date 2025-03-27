package userstorageinfra

import (
	"context"
	"fmt"
	"userservice/internal/core/entities"
	"userservice/internal/pkg/persistence"
)

const (
	addUserMapperID    = "service.user.insert"
	getUserMapperID    = "service.user.getById"
	editUserMapperID   = "service.user.update"
	removeUserMapperID = "service.user.delete"
	listUsersMapperID  = "service.user.list"
	existsUserMapperID = "service.user.exists"
)

func (repository *StorageInfra) Add(ctx context.Context, db CreateInput) (*UserDB, error) {
	err := repository.genericStorage.ExecuteStmt(ctx, addUserMapperID, db)
	if err != nil {
		return nil, err
	}

	result, err := repository.Get(ctx, db.ResourceID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, persistence.NewEntryNotAddedError()
	}

	return &result[0], nil
}

func (repository *StorageInfra) Get(ctx context.Context, id entities.ResourceID) ([]UserDB, error) {
	var userDBItems []UserDB
	db := UserDB{}
	db.ID = id.ID

	err := repository.genericStorage.QueryAll(ctx, getUserMapperID, db, &userDBItems)
	if err != nil {
		return nil, err
	}
	return userDBItems, nil
}

func (repository *StorageInfra) Edit(ctx context.Context, db EditInput) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, err := repository.Exists(ctx, db.ResourceID)
	if err != nil {
		return nil, err
	}

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editUserMapperID, db)
}

func (repository *StorageInfra) Remove(ctx context.Context, id entities.ResourceID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := UserDB{}
	db.ID = id.ID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeUserMapperID, db)
}

func (repository *StorageInfra) List(ctx context.Context, filters UserDBFilter) ([]UserDB, error) {
	userDBItems := make([]UserDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listUsersMapperID, &filters, &userDBItems)
	if err != nil {
		return nil, err
	}
	return userDBItems, nil
}

func (repository *StorageInfra) Exists(ctx context.Context, id entities.ResourceID) ([]ExistsOutput, error) {
	var existsResult []ExistsOutput

	db := UserDB{}
	db.ResourceID = id

	err := repository.genericStorage.QueryAll(ctx, existsUserMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type StorageInfraOptions struct {
	GenericStorage persistence.Storage
}

type StorageInfra struct {
	genericStorage persistence.Storage
}

func ProvideStorageInfra(options StorageInfraOptions) (*StorageInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &StorageInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
