package filter

import "go.uber.org/dig"

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
	Id       string
	WalletId string
}
