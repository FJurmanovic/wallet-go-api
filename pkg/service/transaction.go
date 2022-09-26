package service

import (
	"context"
	"fmt"
	"math"
	"time"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
)

type TransactionService struct {
	repository               *repository.TransactionRepository
	subscriptionRepository   *repository.SubscriptionRepository
	transactionStatusService *TransactionStatusService
}

func NewTransactionService(repository *repository.TransactionRepository, sr *repository.SubscriptionRepository, tss *TransactionStatusService) *TransactionService {
	return &TransactionService{
		repository:               repository,
		subscriptionRepository:   sr,
		transactionStatusService: tss,
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
func (as *TransactionService) New(ctx context.Context, body *model.NewTransactionBody) (*model.Transaction, *model.Exception) {
	exceptionReturn := new(model.Exception)
	tm := new(model.Transaction)

	tsFlt := filter.NewTransactionStatusFilter(model.Params{})
	tsFlt.Status = "completed"
	transactionStatus, exceptionReturn := as.transactionStatusService.Get(ctx, tsFlt)

	if exceptionReturn != nil {
		return nil, exceptionReturn
	}

	amount, _ := body.Amount.Float64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate
	tm.Amount = float32(math.Round(amount*100) / 100)
	tm.TransactionStatusID = transactionStatus.Id

	if body.TransactionDate.IsZero() {
		tm.TransactionDate = time.Now()
	}

	response, err := as.repository.New(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400116"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

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
func (as *TransactionService) GetAll(ctx context.Context, filtered *model.FilteredResponse, flt *filter.TransactionFilter) *model.Exception {
	exceptionReturn := new(model.Exception)

	err := as.repository.GetAll(ctx, filtered, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400118"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row(s) in \"transaction\" table: %s", err)
		return exceptionReturn
	}
	return nil
}

/*
Check

Checks subscriptions and create transacitons.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*model.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) Check(ctx context.Context, flt *filter.TransactionFilter) *model.Exception {
	exceptionReturn := new(model.Exception)
	filtered := new(model.FilteredResponse)

	err := as.repository.Check(ctx, filtered, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400120"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transaction\" table: %s", err)
		return exceptionReturn
	}
	return nil
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
func (as *TransactionService) Edit(ctx context.Context, body *model.TransactionEdit, id string) (*model.Transaction, *model.Exception) {
	amount, _ := body.Amount.Float64()

	exceptionReturn := new(model.Exception)

	tm := new(model.Transaction)
	tm.Id = id
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.TransactionDate = body.TransactionDate
	tm.TransactionStatusID = body.TransactionStatusID
	tm.Amount = float32(math.Round(amount*100) / 100)

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
Bulk Edit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			?[]model.Transaction Bulk Edit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionService) BulkEdit(ctx context.Context, body *[]model.TransactionEdit) (*[]model.Transaction, *model.Exception) {
	transactions := new([]model.Transaction)
	exceptionReturn := new(model.Exception)

	for _, transaction := range *body {

		amount, _ := transaction.Amount.Float64()

		tm := new(model.Transaction)
		tm.Id = transaction.Id
		tm.Description = transaction.Description
		tm.WalletID = transaction.WalletID
		tm.TransactionTypeID = transaction.TransactionTypeID
		tm.TransactionDate = transaction.TransactionDate
		tm.Amount = float32(math.Round(amount*100) / 100)

		*transactions = append(*transactions, *tm)
	}

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
