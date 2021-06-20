package services

import (
	"math"
	"time"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type SubscriptionService struct {
	Db *pg.DB
}

func (as *SubscriptionService) New(body *models.NewSubscriptionBody) *models.Subscription {
	tm := new(models.Subscription)

	amount, _ := body.Amount.Float64()
	customRange, _ := body.CustomRange.Int64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.SubscriptionTypeID = body.SubscriptionTypeID
	tm.CustomRange = int(customRange)
	tm.Description = body.Description
	tm.StartDate = body.StartDate
	tm.HasEnd = body.HasEnd
	tm.EndDate = body.EndDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	if body.StartDate.IsZero() {
		tm.StartDate = time.Now()
	}

	as.Db.Model(tm).Insert()

	as.SubToTrans(tm)

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

func (as *SubscriptionService) SubToTrans(subModel *models.Subscription) {
	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)

	startDate := subModel.StartDate.Local()
	stopDate := firstOfNextMonth
	if subModel.HasEnd && subModel.EndDate.Local().Before(firstOfNextMonth) {
		stopDate = subModel.EndDate.Local()
	}

	transactions := new([]models.Transaction)

	if subModel.SubscriptionType == nil {
		st := new(models.SubscriptionType)
		as.Db.Model(st).Where("? = ?", pg.Ident("id"), subModel.SubscriptionTypeID).Select()
		subModel.SubscriptionType = st
	}

	for startDate.Before(stopDate) {
		trans := subModel.ToTrans()
		trans.TransactionDate = startDate
		if startDate.After(subModel.LastTransactionDate) {
			*transactions = append(*transactions, *trans)
		}
		if subModel.SubscriptionType.Type == "monthly" {
			startDate = startDate.AddDate(0, subModel.CustomRange, 0)
		} else if subModel.SubscriptionType.Type == "weekly" {
			startDate = startDate.AddDate(0, 0, 7*subModel.CustomRange)
		} else if subModel.SubscriptionType.Type == "daily" {
			startDate = startDate.AddDate(0, 0, subModel.CustomRange)
		} else {
			startDate = startDate.AddDate(subModel.CustomRange, 0, 0)
		}
	}

	if len(*transactions) > 0 {
		for _, trans := range *transactions {
			_, err := as.Db.Model(&trans).Where("? = ?", pg.Ident("transaction_date"), trans.TransactionDate).Where("? = ?", pg.Ident("subscription_id"), trans.SubscriptionID).OnConflict("DO NOTHING").SelectOrInsert()
			if err != nil {
				as.Db.Model(subModel).Set("? = ?", pg.Ident("last_transaction_date"), trans.TransactionDate).WherePK().Update()
			}
		}
	}
}
