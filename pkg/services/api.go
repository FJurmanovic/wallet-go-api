package services

import (
	"context"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type ApiService struct {
	Db *pg.DB
}

/*
GetFirst

Gets first row from API table.
   	Args:
   		context.Context: Application context
	Returns:
		models.ApiModel: Api object from database.
*/
func (as *ApiService) GetFirst(ctx context.Context) models.ApiModel {
	db := as.Db.WithContext(ctx)

	apiModel := models.ApiModel{Api: "Works"}
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
		*models.MessageResponse: Message response object.
		*models.Exception: Exception response object.
*/
func (as *ApiService) PostMigrate(ctx context.Context, version string) (*models.MessageResponse, *models.Exception) {
	db := as.Db.WithContext(ctx)

	mr := new(models.MessageResponse)
	er := new(models.Exception)

	migrate.Start(db, version)

	return mr, er
}
