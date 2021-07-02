package services

import (
	"context"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionTypeService struct {
	Db *pg.DB
}

func (as *SubscriptionTypeService) New(ctx context.Context, body *models.NewSubscriptionTypeBody) *models.SubscriptionType {
	db := as.Db.WithContext(ctx)

	tm := new(models.SubscriptionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	db.Model(tm).Insert()

	return tm
}

func (as *SubscriptionTypeService) GetAll(ctx context.Context, embed string) *[]models.SubscriptionType {
	db := as.Db.WithContext(ctx)

	wm := new([]models.SubscriptionType)

	query := db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
