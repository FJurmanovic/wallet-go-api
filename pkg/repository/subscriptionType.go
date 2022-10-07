package repository

import (
	"context"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionTypeRepository struct {
	db *pg.DB
}

func NewSubscriptionTypeRepository(db *pg.DB) *SubscriptionTypeRepository {
	return &SubscriptionTypeRepository{
		db: db,
	}
}

/*
New

Inserts new row to subscription type table.

	   	Args:
			context.Context: Application context
			*model.NewSubscriptionTypeBody: Values to create new row
		Returns:
			*model.SubscriptionType: Created row from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionTypeRepository) New(ctx context.Context, tm *model.SubscriptionType) (*model.SubscriptionType, error) {
	db := as.db.WithContext(ctx)

	_, err := db.Model(tm).Insert()
	if err != nil {
		return nil, err
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
			*[]model.SubscriptionType: List of subscription type objects.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionTypeRepository) GetAll(ctx context.Context, flt *filter.SubscriptionTypeFilter) (*[]model.SubscriptionType, error) {
	wm := new([]model.SubscriptionType)
	db := as.db.WithContext(ctx)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, flt.Embed).Select()
	if err != nil {
		return nil, err
	}

	return wm, nil
}
