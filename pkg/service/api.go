package service

import (
	"context"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
)

type ApiService struct {
	repository *repository.ApiRepository
}

func NewApiService(repository *repository.ApiRepository) *ApiService {
	return &ApiService{
		repository,
	}
}

/*
GetFirst

Gets first row from API table.

	   	Args:
	   		context.Context: Application context
		Returns:
			model.ApiModel: Api object from database.
*/
func (as ApiService) GetFirst(ctx context.Context) model.ApiModel {
	return as.repository.GetFirst(ctx)
}

/*
PostMigrate

Starts database migration.

	   	Args:
	   		context.Context: Application context
			string: Migration version
		Returns:
			*model.MessageResponse: Message response object.
			*model.Exception: Exception response object.
*/
func (as ApiService) PostMigrate(ctx context.Context, version string) (*model.MessageResponse, *model.Exception) {
	return as.repository.PostMigrate(ctx, version)
}
