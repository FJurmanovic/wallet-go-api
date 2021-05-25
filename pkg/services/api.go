package services

import (
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type ApiService struct {
	Db *pg.DB
}

func (as *ApiService) GetFirst() models.ApiModel {
	apiModel := models.ApiModel{Api: "Works"}
	as.Db.Model(&apiModel).First()
	return apiModel
}

func (as *ApiService) PostMigrate() (*models.MessageResponse, *models.Exception) {
	mr := new(models.MessageResponse)
	er := new(models.Exception)

	err := migrate.Start(as.Db)
	if err != nil {
		er.ErrorCode = "400999"
		er.StatusCode = 400
		er.Message = err.Error()
	}

	return mr, er
}
