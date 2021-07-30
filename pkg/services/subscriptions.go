package services

import (
	"context"
	"math"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionService struct {
	Db *pg.DB
}

func (as *SubscriptionService) New(ctx context.Context, body *models.NewSubscriptionBody) *models.Subscription {
	db := as.Db.WithContext(ctx)

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

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(tm).Insert()

	as.SubToTrans(tm, tx)
	tx.Commit()

	return tm
}

func (as *SubscriptionService) Get(ctx context.Context, am *models.Auth, id string, params *models.Params) *models.Subscription {
	db := as.Db.WithContext(ctx)

	wm := new(models.Subscription)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	common.GenerateEmbed(qry, params.Embed).WherePK().Select()

	if (*wm).HasNew() {
		as.SubToTrans(wm, tx)
	}

	tx.Commit()

	return wm
}

func (as *SubscriptionService) GetAll(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse) {
	db := as.Db.WithContext(ctx)

	wm := new([]models.Subscription)

	tx, _ := db.Begin()
	defer tx.Rollback()

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}

	for _, sub := range *wm {
		if sub.HasNew() {
			as.SubToTrans(&sub, tx)
		}
	}
	FilteredResponse(query, wm, filtered)
	tx.Commit()
}

func (as *SubscriptionService) Edit(ctx context.Context, body *models.SubscriptionEdit, id string) *models.Subscription {
	db := as.Db.WithContext(ctx)

	amount, _ := body.Amount.Float64()

	tm := new(models.Subscription)
	tm.Id = id
	tm.EndDate = body.EndDate
	tm.HasEnd = body.HasEnd
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.Amount = float32(math.Round(amount*100) / 100)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(tm).WherePK().UpdateNotZero()

	tx.Commit()

	return tm
}

func (as *SubscriptionService) SubToTrans(subModel *models.Subscription, tx *pg.Tx) {

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	//tzFirstOfNextMonth := firstOfNextMonth.In(subModel.StartDate.Location())

	startDate := subModel.StartDate
	stopDate := firstOfNextMonth
	if subModel.HasEnd && subModel.EndDate.Before(firstOfNextMonth) {
		stopDate = subModel.EndDate
	}

	transactions := new([]models.Transaction)

	if subModel.SubscriptionType == nil {
		st := new(models.SubscriptionType)
		tx.Model(st).Where("? = ?", pg.Ident("id"), subModel.SubscriptionTypeID).Select()
		subModel.SubscriptionType = st
	}

	for startDate.Before(stopDate) {
		trans := subModel.ToTrans()
		trans.TransactionDate = startDate
		if startDate.After(subModel.LastTransactionDate) && (startDate.Before(now) || startDate.Equal(now)) {
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
			_, err := tx.Model(&trans).Where("? = ?", pg.Ident("transaction_date"), trans.TransactionDate).Where("? = ?", pg.Ident("subscription_id"), trans.SubscriptionID).OnConflict("DO NOTHING").SelectOrInsert()
			if err != nil {
				tx.Model(subModel).Set("? = ?", pg.Ident("last_transaction_date"), trans.TransactionDate).WherePK().Update()
			}
		}
	}
}
