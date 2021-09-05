package services

import (
	"context"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionStatusService struct {
	Db *pg.DB
}

/*
New

Inserts new row to transaction status table.
   	Args:
		context.Context: Application context
		*models.NewTransactionStatusBody: object to create
	Returns:
		*models.TransactionType: Transaction Type object from database.
*/
func (as *TransactionStatusService) New(ctx context.Context, body *models.NewTransactionStatusBody) *models.TransactionStatus {
	db := as.Db.WithContext(ctx)

	tm := new(models.TransactionStatus)

	tm.Init()
	tm.Name = body.Name
	tm.Status = body.Status

	db.Model(tm).Insert()

	return tm
}

/*
GetAll

Gets all rows from transaction status table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.TransactionStatus: List of Transaction status objects from database.
*/
func (as *TransactionStatusService) GetAll(ctx context.Context, embed string) *[]models.TransactionStatus {
	db := as.Db.WithContext(ctx)

	wm := new([]models.TransactionStatus)

	query := db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
