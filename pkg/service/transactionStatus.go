package service

import (
	"context"
	"fmt"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
)

type TransactionStatusService struct {
	repository *repository.TransactionStatusRepository
}

func NewTransactionStatusService(repository *repository.TransactionStatusRepository) *TransactionStatusService {
	return &TransactionStatusService{
		repository: repository,
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
func (as *TransactionStatusService) New(ctx context.Context, tm *model.TransactionStatus) (*model.TransactionStatus, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.New(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400123"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
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
func (as *TransactionStatusService) GetAll(ctx context.Context, flt *filter.TransactionStatusFilter) (*[]model.TransactionStatus, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.GetAll(ctx, flt, nil)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400124"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
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
	response, err := as.repository.Get(ctx, flt, nil)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400129"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
