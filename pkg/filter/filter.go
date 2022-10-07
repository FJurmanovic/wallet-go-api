package filter

import (
	"go.uber.org/dig"
	"wallet-api/pkg/model"
)

/*
InitializeFilters

Initializes Dependency Injection modules for filters

	Args:
		*dig.Container: Dig Container
*/
func InitializeFilters(c *dig.Container) {
	c.Provide(NewApiFilter)
	c.Provide(NewSubscriptionFilter)
	c.Provide(NewSubscriptionTypeFilter)
	c.Provide(NewTransactionFilter)
	c.Provide(NewTransactionStatusFilter)
	c.Provide(NewTransactionTypeFilter)
	c.Provide(NewUserFilter)
	c.Provide(NewWalletFilter)
}

type BaseFilter struct {
	model.Params
	Id       string
	WalletId string
	UserId   string
}

type BBa interface {
	BaseFilter
	some() bool
}
