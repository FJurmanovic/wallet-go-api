package services

import (
	"context"
	"fmt"
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
		*models.Exception: Exception payload.
*/
func (as *TransactionStatusService) New(ctx context.Context, body *models.NewTransactionStatusBody) (*models.TransactionStatus, *models.Exception) {
	db := as.Db.WithContext(ctx)

	tm := new(models.TransactionStatus)
	exceptionReturn := new(models.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Status = body.Status

	_, err := db.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400123"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return tm, nil
}

/*
GetAll

Gets all rows from transaction status table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.TransactionStatus: List of Transaction status objects from database.
		*models.Exception: Exception payload.
*/
func (as *TransactionStatusService) GetAll(ctx context.Context, embed string) (*[]models.TransactionStatus, *models.Exception) {
	db := as.Db.WithContext(ctx)

	wm := new([]models.TransactionStatus)
	exceptionReturn := new(models.Exception)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, embed).Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400124"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return wm, nil
}
