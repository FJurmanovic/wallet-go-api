package service

import (
	"context"
	"fmt"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
)

type SubscriptionTypeService struct {
	repository *repository.SubscriptionTypeRepository
}

func NewSubscriptionTypeService(repository *repository.SubscriptionTypeRepository) *SubscriptionTypeService {
	return &SubscriptionTypeService{
		repository,
	}
}

/*
New

Inserts new row to subscription type table.

	   	Args:
			context.Context: Application context
			*model.NewSubscriptionTypeBody: Values to create new row
		Returns:
			*model.SubscriptionType: Created row from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionTypeService) New(ctx context.Context, body *model.NewSubscriptionTypeBody) (*model.SubscriptionType, *model.Exception) {
	tm := new(model.SubscriptionType)
	exceptionReturn := new(model.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	response, err := as.repository.New(ctx, tm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400114"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscriptionTypes\" table: %s", err)
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
			*[]model.SubscriptionType: List of subscription type objects.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionTypeService) GetAll(ctx context.Context, embed string) (*[]model.SubscriptionType, *model.Exception) {
	wm := new([]model.SubscriptionType)
	exceptionReturn := new(model.Exception)

	flt := filter.NewSubscriptionTypeFilter(model.Params{
		Embed: embed,
	})
	response, err := as.repository.GetAll(ctx, flt, wm)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400135"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"subscriptionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return response, nil
}
