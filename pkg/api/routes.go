package api

import (
	"wallet-api/pkg/controllers"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/configs"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func Routes(s *gin.Engine, db *pg.DB) {
	ver := s.Group(configs.Prefix)

	api := ver.Group("api")
	auth := ver.Group("auth")
	wallet := ver.Group("wallet", middleware.Auth)
	walletHeader := ver.Group("wallet/wallet-header", middleware.Auth)
	transaction := ver.Group("transaction", middleware.Auth)
	transactionType := ver.Group("transaction-type", middleware.Auth)
	subscription := ver.Group("subscription", middleware.Auth)
	subscriptionType := ver.Group("subscription-type", middleware.Auth)

	apiService := services.ApiService{Db: db}
	usersService := services.UsersService{Db: db}
	walletService := services.WalletService{Db: db}
	transactionService := services.TransactionService{Db: db}
	transactionTypeService := services.TransactionTypeService{Db: db}
	subscriptionService := services.SubscriptionService{Db: db}
	subscriptionTypeService := services.SubscriptionTypeService{Db: db}

	walletService.Ss = &subscriptionService

	controllers.NewApiController(&apiService, api)
	controllers.NewAuthController(&usersService, auth)
	controllers.NewWalletsController(&walletService, wallet)
	controllers.NewWalletsHeaderController(&walletService, walletHeader)
	controllers.NewTransactionController(&transactionService, transaction)
	controllers.NewTransactionTypeController(&transactionTypeService, transactionType)
	controllers.NewSubscriptionController(&subscriptionService, subscription)
	controllers.NewSubscriptionTypeController(&subscriptionTypeService, subscriptionType)
}
