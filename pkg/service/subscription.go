package service

import (
	"context"
	"fmt"
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
		repository: repository,
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
func (as *SubscriptionService) New(ctx context.Context, tm *model.Subscription) (*model.Subscription, *model.Exception) {
	exceptionReturn := new(model.Exception)

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
func (as *SubscriptionService) Get(ctx context.Context, flt filter.SubscriptionFilter) (*model.Subscription, *model.Exception) {
	exceptionReturn := new(model.Exception)
	response, err := as.repository.Get(ctx, flt)
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
func (as *SubscriptionService) GetAll(ctx context.Context, flt *filter.SubscriptionFilter) (*model.FilteredResponse, *model.Exception) {
	exceptionReturn := new(model.Exception)

	filtered, err := as.repository.GetAll(ctx, flt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400110"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	return filtered, nil
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
func (as *SubscriptionService) Edit(ctx context.Context, tm *model.Subscription) (*model.Subscription, *model.Exception) {

	exceptionReturn := new(model.Exception)

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
