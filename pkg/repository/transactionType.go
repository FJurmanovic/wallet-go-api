package repository

import (
	"context"
	"fmt"
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
func (as *TransactionTypeRepository) New(ctx context.Context, body *model.NewTransactionTypeBody) (*model.TransactionType, *model.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(model.TransactionType)
	exceptionReturn := new(model.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	_, err := db.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400125"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
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
func (as *TransactionTypeRepository) GetAll(ctx context.Context, embed string) (*[]model.TransactionType, *model.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]model.TransactionType)
	exceptionReturn := new(model.Exception)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, embed).Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400133"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return wm, nil
}
