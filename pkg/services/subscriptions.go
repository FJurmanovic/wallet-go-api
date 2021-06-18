package services

import (
	"time"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type SubscriptionService struct {
	Db *pg.DB
}

func (as *SubscriptionService) New(body *models.NewSubscriptionBody) *models.Subscription {
	tm := new(models.Subscription)

	amount, _ := body.Amount.Int64()
	customRange, _ := body.CustomRange.Int64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.SubscriptionTypeID = body.SubscriptionTypeID
	tm.CustomRange = int(customRange)
	tm.Description = body.Description
	tm.StartDate = body.StartDate
	tm.Amount = int(amount)

	if body.StartDate.IsZero() {
		tm.StartDate = time.Now()
	}

	as.Db.Model(tm).Insert()

	return tm
}

func (as *SubscriptionService) GetAll(am *models.Auth, walletId string, filtered *models.FilteredResponse) {
	wm := new([]models.Subscription)

	query := as.Db.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	FilteredResponse(query, wm, filtered)
}
