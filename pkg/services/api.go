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

func (as *ApiService) GetFirst(ctx context.Context) models.ApiModel {
	db := as.Db.WithContext(ctx)

	apiModel := models.ApiModel{Api: "Works"}
	db.Model(&apiModel).First()
	return apiModel
}

func (as *ApiService) PostMigrate(ctx context.Context, version string) (*models.MessageResponse, *models.Exception) {
	db := as.Db.WithContext(ctx)

	mr := new(models.MessageResponse)
	er := new(models.Exception)

	migrate.Start(db, version)

	return mr, er
}
