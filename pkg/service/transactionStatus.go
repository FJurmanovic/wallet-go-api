package service

import (
	"context"
	"fmt"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionStatusService struct {
	db *pg.DB
}

func NewTransactionStatusService(db *pg.DB) *TransactionStatusService {
	return &TransactionStatusService{
		db: db,
	}
}

/*
New

Inserts new row to transaction status table.

	   	Args:
			context.Context: Application context
			*model.NewTransactionStatusBody: object to create
		Returns:
			*model.TransactionType: Transaction Type object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusService) New(ctx context.Context, body *model.NewTransactionStatusBody) (*model.TransactionStatus, *model.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(model.TransactionStatus)
	exceptionReturn := new(model.Exception)

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
			*[]model.TransactionStatus: List of Transaction status objects from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusService) GetAll(ctx context.Context, embed string) (*[]model.TransactionStatus, *model.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]model.TransactionStatus)
	exceptionReturn := new(model.Exception)

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

func (as *TransactionStatusService) Get(ctx context.Context, flt *filter.TransactionStatusFilter) (*model.TransactionStatus, *model.Exception) {
	transactionStatus := new(model.TransactionStatus)
	exceptionReturn := new(model.Exception)
	if flt.Id != "" {
		transactionStatus.Id = flt.Id
	}
	if flt.Status != "" {
		transactionStatus.Status = flt.Status
	}
	response, err := as.repository.Get(ctx, wm, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400129"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
