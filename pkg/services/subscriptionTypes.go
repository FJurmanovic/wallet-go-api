package services

import (
	"context"
	"fmt"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionTypeService struct {
	db *pg.DB
}

func NewSubscriptionTypeService(db *pg.DB) *SubscriptionTypeService {
	return &SubscriptionTypeService{
		db: db,
	}
}

/*
New

Inserts new row to subscription type table.

	   	Args:
			context.Context: Application context
			*models.NewSubscriptionTypeBody: Values to create new row
		Returns:
			*models.SubscriptionType: Created row from database.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionTypeService) New(ctx context.Context, body *models.NewSubscriptionTypeBody) (*models.SubscriptionType, *models.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(models.SubscriptionType)
	exceptionReturn := new(models.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	_, err := db.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400114"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscriptionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return tm, nil
}

/*
GetAll

Gets all rows from subscription type table.

	   	Args:
			context.Context: Application context
			string: Relations to embed
		Returns:
			*[]models.SubscriptionType: List of subscription type objects.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionTypeService) GetAll(ctx context.Context, embed string) (*[]models.SubscriptionType, *models.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]models.SubscriptionType)
	exceptionReturn := new(models.Exception)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, embed).Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400135"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"subscriptionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return wm, nil
}
