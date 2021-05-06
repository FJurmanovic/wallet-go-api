package services

import (
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
