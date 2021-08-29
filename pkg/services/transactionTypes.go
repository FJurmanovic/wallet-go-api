package services

import (
	"context"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionTypeService struct {
	Db *pg.DB
}

/*
New

Inserts new row to transaction type table.
   	Args:
		context.Context: Application context
		*models.NewTransactionTypeBody: object to create
	Returns:
		*models.TransactionType: Transaction Type object from database.
*/
func (as *TransactionTypeService) New(ctx context.Context, body *models.NewTransactionTypeBody) *models.TransactionType {
	db := as.Db.WithContext(ctx)

	tm := new(models.TransactionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	db.Model(tm).Insert()

	return tm
}

/*
GetAll

Gets all rows from transaction type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.TransactionType: List of Transaction type objects from database.
*/
func (as *TransactionTypeService) GetAll(ctx context.Context, embed string) *[]models.TransactionType {
	db := as.Db.WithContext(ctx)

	wm := new([]models.TransactionType)

	query := db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
