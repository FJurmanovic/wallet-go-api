package service

import (
	"go.uber.org/dig"
	"wallet-api/pkg/repository"
)

/*
InitializeServices

Initializes Dependency Injection modules for services

	Args:
		*dig.Container: Dig Container
*/
func InitializeServices(c *dig.Container) {
	repository.InitializeRepositories(c)

	c.Provide(NewApiService)
	c.Provide(NewSubscriptionService)
	c.Provide(NewSubscriptionTypeService)
	c.Provide(NewTransactionService)
	c.Provide(NewTransactionStatusService)
	c.Provide(NewTransactionTypeService)
	c.Provide(NewUserService)
	c.Provide(NewWalletService)
}
