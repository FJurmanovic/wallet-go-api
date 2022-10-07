package repository

import (
	"go.uber.org/dig"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

/*
InitializeRepositories

Initializes Dependency Injection modules for repositories

	Args:
		*dig.Container: Dig Container
*/
func InitializeRepositories(c *dig.Container) {
	c.Provide(NewApiRepository)
	c.Provide(NewSubscriptionRepository)
	c.Provide(NewSubscriptionTypeRepository)
	c.Provide(NewTransactionRepository)
	c.Provide(NewTransactionStatusRepository)
	c.Provide(NewTransactionTypeRepository)
	c.Provide(NewUserRepository)
	c.Provide(NewWalletRepository)
}

/*
FilteredResponse

Adds filter to query and executes it.

	   	Args:
	   		*pg.Query: postgres query
			interface{}: model to be mapped from query execution.
			*model.FilteredResponse: filter options.
*/
func FilteredResponse(qry *pg.Query, mdl interface{}, params model.Params) (*model.FilteredResponse, error) {
	filtered := new(model.FilteredResponse)
	filtered.Params = params
	if filtered.Page == 0 {
		filtered.Page = 1
	}
	if filtered.Rpp == 0 {
		filtered.Rpp = 20
	}
	if filtered.SortBy == "" {
		filtered.SortBy = "date_created DESC"
	}
	qry = qry.Limit(filtered.Rpp).Offset((filtered.Page - 1) * filtered.Rpp).Order(filtered.SortBy)
	common.GenerateEmbed(qry, filtered.Embed)
	count, err := qry.SelectAndCount()
	common.CheckError(err)

	filtered.TotalRecords = count
	filtered.Items = mdl

	return filtered, err
}
