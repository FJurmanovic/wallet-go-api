package service

import (
	"context"
	"fmt"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
)

type TransactionTypeService struct {
	repository *repository.TransactionTypeRepository
}

func NewTransactionTypeService(repository *repository.TransactionTypeRepository) *TransactionTypeService {
	return &TransactionTypeService{
		repository: repository,
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
func (as *TransactionTypeService) New(ctx context.Context, tm *model.TransactionType) (*model.TransactionType, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.New(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400125"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
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
func (as *TransactionTypeService) GetAll(ctx context.Context, flt *filter.TransactionTypeFilter) (*[]model.TransactionType, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.GetAll(ctx, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400133"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
