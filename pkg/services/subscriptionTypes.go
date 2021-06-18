package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionTypeService struct {
	Db *pg.DB
}

func (as *SubscriptionTypeService) New(body *models.NewSubscriptionTypeBody) *models.SubscriptionType {
	tm := new(models.SubscriptionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	as.Db.Model(tm).Insert()

	return tm
}

func (as *SubscriptionTypeService) GetAll(embed string) *[]models.SubscriptionType {
	wm := new([]models.SubscriptionType)

	query := as.Db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
