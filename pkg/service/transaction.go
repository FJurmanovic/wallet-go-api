package service

import (
	"context"
	"fmt"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
	"wallet-api/pkg/utl/common"
)

type TransactionService struct {
	repository                  *repository.TransactionRepository
	subscriptionRepository      *repository.SubscriptionRepository
	transactionStatusRepository *repository.TransactionStatusRepository
}

func NewTransactionService(repository *repository.TransactionRepository, sr *repository.SubscriptionRepository, tsr *repository.TransactionStatusRepository) *TransactionService {
	return &TransactionService{
		repository:                  repository,
		subscriptionRepository:      sr,
		transactionStatusRepository: tsr,
	}
}

/*
New row into transaction table

Inserts

	   	Args:
			context.Context: Application context
			*model.NewTransactionBody: Transaction body object
		Returns:
			*model.Transaction: Transaction object
			*model.Exception: Exception payload.
*/
func (as *TransactionService) New(ctx context.Context, tm *model.Transaction) (*model.Transaction, *model.Exception) {
	exceptionReturn := new(model.Exception)

	tx, err := as.repository.CreateTx(ctx)
	defer tx.Rollback()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400136"
		exceptionReturn.Message = fmt.Sprintf("Error beginning transaction: %s", err)
		return nil, exceptionReturn
	}

	tsFlt := filter.NewTransactionStatusFilter(model.Params{})
	tsFlt.Status = "completed"
	transactionStatuses, err := as.transactionStatusRepository.GetAll(ctx, tsFlt, tx)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400138"
		exceptionReturn.Message = fmt.Sprintf("Error fetching transactionStatus: %s", err)
		return nil, exceptionReturn
	}

	var transactionStatus = common.Find[model.TransactionStatus](transactionStatuses, func(status *model.TransactionStatus) bool {
		return status.Status == "completed"
	})
	tm.TransactionStatusID = transactionStatus.Id

	response, err := as.repository.New(ctx, tm, tx)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400116"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}
	tx.Commit()

	return response, nil
}

/*
GetAll

Gets all rows from subscription type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*model.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) GetAll(ctx context.Context, flt *filter.TransactionFilter) (*model.FilteredResponse, *model.Exception) {
	exceptionReturn := new(model.Exception)

	filtered, err := as.repository.GetAll(ctx, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400118"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row(s) in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}
	return filtered, nil
}

/*
Check

Checks subscriptions and create transactions.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*model.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) Check(ctx context.Context, flt *filter.TransactionFilter) (*model.FilteredResponse, *model.Exception) {
	exceptionReturn := new(model.Exception)

	filtered, err := as.repository.Check(ctx, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400120"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}
	return filtered, nil
}

/*
Edit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			*model.TransactionEdit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionService) Edit(ctx context.Context, tm *model.Transaction) (*model.Transaction, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.Edit(ctx, tm)

	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400107"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}

/*
BulkEdit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			?[]model.Transaction Bulk Edit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionService) BulkEdit(ctx context.Context, transactions *[]model.Transaction) (*[]model.Transaction, *model.Exception) {
	exceptionReturn := new(model.Exception)

	response, err := as.repository.BulkEdit(ctx, transactions)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400121"
		exceptionReturn.Message = fmt.Sprintf("Error updating rows in \"transactions\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}

/*
Get

Gets row from transaction table by id.

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
			string: id to search
			*model.Params: url query parameters
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionService) Get(ctx context.Context, flt *filter.TransactionFilter) (*model.Transaction, *model.Exception) {

	exceptionReturn := new(model.Exception)

	response, err := as.repository.Get(ctx, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400122"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactions\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
