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

/*
New

Inserts new row to subscription type table.
   	Args:
		context.Context: Application context
		*models.NewSubscriptionTypeBody: Values to create new row
	Returns:
		*models.SubscriptionType: Created row from database.
*/
func (as *SubscriptionTypeService) New(ctx context.Context, body *models.NewSubscriptionTypeBody) *models.SubscriptionType {
	db := as.Db.WithContext(ctx)

	tm := new(models.SubscriptionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	db.Model(tm).Insert()

	return tm
}

/*
GetAll

Gets all rows from subscription type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.SubscriptionType: List of subscription type objects.
*/
func (as *SubscriptionTypeService) GetAll(ctx context.Context, embed string) *[]models.SubscriptionType {
	db := as.Db.WithContext(ctx)

	wm := new([]models.SubscriptionType)

	query := db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
