package repository

import (
	"context"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
)

type ApiRepository struct {
	db *pg.DB
}

func NewApiRepository(db *pg.DB) *ApiRepository {
	return &ApiRepository{
		db: db,
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
func (as ApiRepository) GetFirst(ctx context.Context) model.ApiModel {
	db := as.db.WithContext(ctx)
	apiModel := model.ApiModel{Api: "Works"}
	db.Model(&apiModel).First()
	return apiModel
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
func (as ApiRepository) PostMigrate(ctx context.Context, version string) []error {
	db := as.db.WithContext(ctx)
	return migrate.Start(db, version)
}
