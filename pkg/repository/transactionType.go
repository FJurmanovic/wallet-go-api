package repository

import (
	"context"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionTypeRepository struct {
	db *pg.DB
}

func NewTransactionTypeRepository(db *pg.DB) *TransactionTypeRepository {
	return &TransactionTypeRepository{
		db: db,
	}
}

/*
New

Inserts new row to transaction type table.

	   	Args:
			context.Context: Application context
			*model.NewTransactionTypeBody: object to create
		Returns:
			*model.TransactionType: Transaction Type object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionTypeRepository) New(ctx context.Context, tm *model.TransactionType) (*model.TransactionType, error) {
	db := as.db.WithContext(ctx)

	_, err := db.Model(tm).Insert()
	if err != nil {
		return nil, err
	}

	return tm, nil
}

/*
GetAll

Gets all rows from transaction type table.

	   	Args:
			context.Context: Application context
			string: Relations to embed
		Returns:
			*[]model.TransactionType: List of Transaction type objects from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionTypeRepository) GetAll(ctx context.Context, flt *filter.TransactionTypeFilter) (*[]model.TransactionType, error) {
	db := as.db.WithContext(ctx)

	wm := new([]model.TransactionType)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, flt.Embed).Select()
	if err != nil {
		return nil, err
	}

	return wm, nil
}
