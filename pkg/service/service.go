package service

import (
	"wallet-api/pkg/repository"

	"go.uber.org/dig"
)

/*
InitializeServices

Initializes Dependency Injection modules for services

	Args:
		*dig.Container: Dig Container
*/
func InitializeServices(c *dig.Scope) {
	repository.InitializeRepositories(c)

	c.Provide(NewApiService)
	c.Provide(NewSubscriptionService)
	c.Provide(NewSubscriptionTypeService)
	c.Provide(NewTransactionService)
	c.Provide(NewTransactionStatusService)
	c.Provide(NewTransactionTypeService)
	c.Provide(NewUserService)
	c.Provide(NewWalletService)
	c.Provide(NewCurrencyService)
}
