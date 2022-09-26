package api

import (
	"wallet-api/pkg/controllers"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"
	"wallet-api/pkg/utl/configs"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"go.uber.org/dig"
)

/*
Routes

Initializes web api controllers and its corresponding routes.

	Args:
		*gin.Engine: Gin Engine
		*pg.DB: Postgres database client
*/
func Routes(s *gin.Engine, db *pg.DB) {
	c := dig.New()
	ver := s.Group(configs.Prefix)

	routeGroups := &common.RouteGroups{
		Api:               ver.Group("api"),
		Auth:              ver.Group("auth"),
		Wallet:            ver.Group("wallet", middleware.Auth),
		WalletHeader:      ver.Group("wallet/wallet-header", middleware.Auth),
		Transaction:       ver.Group("transaction", middleware.Auth),
		TransactionType:   ver.Group("transaction-type", middleware.Auth),
		Subscription:      ver.Group("subscription", middleware.Auth),
		SubscriptionType:  ver.Group("subscription-type", middleware.Auth),
		TransactionStatus: ver.Group("transaction-status", middleware.Auth),
	}

	s.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	c.Provide(func() *common.RouteGroups {
		return routeGroups
	})
	c.Provide(func() *pg.DB {
		return db
	})
	c.Provide(services.NewApiService)
	c.Provide(services.NewSubscriptionService)
	c.Provide(services.NewSubscriptionTypeService)
	c.Provide(services.NewTransactionService)
	c.Provide(services.NewTransactionStatusService)
	c.Provide(services.NewTransactionTypeService)
	c.Provide(services.NewUsersService)
	c.Provide(services.NewWalletService)

	c.Invoke(controllers.NewApiController)
	c.Invoke(controllers.NewAuthController)
	c.Invoke(controllers.NewWalletsController)
	c.Invoke(controllers.NewWalletsHeaderController)
	c.Invoke(controllers.NewTransactionController)
	c.Invoke(controllers.NewTransactionStatusController)
	c.Invoke(controllers.NewTransactionTypeController)
	c.Invoke(controllers.NewSubscriptionController)
	c.Invoke(controllers.NewSubscriptionTypeController)
}
