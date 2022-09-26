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

type SubscriptionService struct {
	repository *repository.SubscriptionRepository
}

func NewSubscriptionService(repository *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repository,
	}
}

/*
New

Inserts new row to subscription table.

	   	Args:
	   		context.Context: Application context
			*model.NewSubscriptionBody: Request body
		Returns:
			*model.Subscription: Created Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionService) New(ctx context.Context, body *model.NewSubscriptionBody) (*model.Subscription, *model.Exception) {
	tm := new(model.Subscription)
	exceptionReturn := new(model.Exception)

	amount, _ := body.Amount.Float64()
	customRange, _ := body.CustomRange.Int64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.SubscriptionTypeID = body.SubscriptionTypeID
	tm.CustomRange = int(customRange)
	tm.Description = body.Description
	tm.StartDate = body.StartDate
	tm.HasEnd = body.HasEnd
	tm.EndDate = body.EndDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	if body.StartDate.IsZero() {
		tm.StartDate = time.Now()
	}

	response, err := as.repository.New(ctx, tm)

	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400109"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}
	return response, nil
}

/*
Get

Gets row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication model
			string: subscription id to search
			params: *model.Params
		Returns:
			*model.Subscription: Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionService) Get(ctx context.Context, am *model.Auth, flt filter.SubscriptionFilter) (*model.Subscription, *model.Exception) {
	exceptionReturn := new(model.Exception)
	wm := new(model.Subscription)
	wm.Id = flt.Id
	response, err := as.repository.Get(ctx, wm, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400129"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}

/*
GetAll

Gets filtered rows from subscription table.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication object
			string: Wallet id to search
			*model.FilteredResponse: filter options
		Returns:
			*model.Exception: Exception payload.
*/
func (as *SubscriptionService) GetAll(ctx context.Context, flt *filter.SubscriptionFilter, filtered *model.FilteredResponse) *model.Exception {
	wm := new([]model.Subscription)
	exceptionReturn := new(model.Exception)

	err := as.repository.GetAll(ctx, wm, filtered, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400110"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"subscription\" table: %s", err)
		return exceptionReturn
	}

	return nil
}

/*
Edit

Updates row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*model.SubscriptionEdit: Values to edit
			string: id to search
		Returns:
			*model.Subscription: Edited Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionService) Edit(ctx context.Context, body *model.SubscriptionEdit, id string) (*model.Subscription, *model.Exception) {
	amount, _ := body.Amount.Float64()
	exceptionReturn := new(model.Exception)

	tm := new(model.Subscription)
	tm.Id = id
	tm.EndDate = body.EndDate
	tm.HasEnd = body.HasEnd
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.Amount = float32(math.Round(amount*100) / 100)

	response, err := as.repository.Edit(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400111"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}

/*
End

Updates row in subscription table by id.

Ends subscription with current date.

	   	Args:
	   		context.Context: Application context
			string: id to search
		Returns:
			*model.Subscription: Created Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionService) End(ctx context.Context, id string) (*model.Subscription, *model.Exception) {
	exceptionReturn := new(model.Exception)

	tm := new(model.Subscription)
	tm.Id = id
	tm.EndDate = time.Now()
	tm.HasEnd = true

	response, err := as.repository.End(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400112"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
