package services

import (
	"context"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	Db *pg.DB
	Ss *SubscriptionService
}

// Inserts row to wallets table.
func (as *WalletService) New(ctx context.Context, am *models.NewWalletBody) *models.Wallet {
	db := as.Db.WithContext(ctx)

	walletModel := new(models.Wallet)
	walletModel.Init()
	walletModel.UserID = am.UserID
	walletModel.Name = am.Name
	db.Model(walletModel).Insert()
	return walletModel
}

// Updates row in wallets table by id.
func (as *WalletService) Edit(ctx context.Context, body *models.WalletEdit, id string) *models.Wallet {
	db := as.Db.WithContext(ctx)

	tm := new(models.Wallet)
	tm.Id = id
	tm.Name = body.Name

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(tm).WherePK().UpdateNotZero()

	tx.Commit()

	return tm
}

// Gets row in wallets table by id.
func (as *WalletService) Get(ctx context.Context, id string, params *models.Params) *models.Wallet {
	db := as.Db.WithContext(ctx)

	wm := new(models.Wallet)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	common.GenerateEmbed(qry, params.Embed).WherePK().Select()

	tx.Commit()

	return wm
}

// Gets filtered rows from wallets table.
func (as *WalletService) GetAll(ctx context.Context, am *models.Auth, filtered *models.FilteredResponse) {
	db := as.Db.WithContext(ctx)
	wm := new([]models.Wallet)

	query := db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	FilteredResponse(query, wm, filtered)
}

// Gets row from wallets table.
//
// Calculates previous month, current and next month totals.
func (as *WalletService) GetHeader(ctx context.Context, am *models.Auth, walletId string) *models.WalletHeader {
	db := as.Db.WithContext(ctx)

	wm := new(models.WalletHeader)
	wallets := new([]models.WalletTransactions)
	transactions := new([]models.Transaction)
	subscriptions := new([]models.Subscription)

	tx, _ := db.Begin()
	defer tx.Rollback()

	query2 := tx.Model(subscriptions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType").Relation("SubscriptionType")
	if walletId != "" {
		query2.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query2.Select()

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	firstOfMonthAfterNext := time.Date(currentYear, currentMonth+2, 1, 0, 0, 0, 0, currentLocation)

	for _, sub := range *subscriptions {
		if sub.HasNew() {
			as.Ss.SubToTrans(&sub, tx)
		}
	}

	query := tx.Model(transactions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType")
	if walletId != "" {
		query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query.Select()
	tx.Commit()

	for _, sub := range *subscriptions {
		stopDate := firstOfMonthAfterNext
		if sub.HasEnd && sub.EndDate.Before(firstOfMonthAfterNext) {
			stopDate = sub.EndDate
		}
		startDate := sub.StartDate
		for startDate.Before(stopDate) {
			trans := sub.ToTrans()
			trans.TransactionDate = startDate
			if startDate.After(firstOfNextMonth) || startDate.Equal(firstOfNextMonth) {
				*transactions = append(*transactions, *trans)
			}
			if sub.SubscriptionType.Type == "monthly" {
				startDate = startDate.AddDate(0, sub.CustomRange, 0)
			} else if sub.SubscriptionType.Type == "weekly" {
				startDate = startDate.AddDate(0, 0, 7*sub.CustomRange)
			} else if sub.SubscriptionType.Type == "daily" {
				startDate = startDate.AddDate(0, 0, sub.CustomRange)
			} else {
				startDate = startDate.AddDate(sub.CustomRange, 0, 0)
			}
		}
	}

	for _, trans := range *transactions {
		addWhere(wallets, trans.WalletID, trans)
	}

	for i, wallet := range *wallets {
		for _, trans := range wallet.Transactions {
			// tzFirstOfMonthAfterNext := firstOfMonthAfterNext.In(trans.TransactionDate.Location())
			// tzFirstOfNextMonth := firstOfNextMonth.In(trans.TransactionDate.Location())
			// tzFirstOfMonth := firstOfMonth.In(trans.TransactionDate.Location())
			if trans.TransactionDate.Before(firstOfNextMonth) && trans.TransactionDate.After(firstOfMonth) || trans.TransactionDate.Equal(firstOfMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].CurrentBalance -= trans.Amount
				} else {
					(*wallets)[i].CurrentBalance += trans.Amount
				}
			} else if trans.TransactionDate.Before(firstOfMonthAfterNext) && trans.TransactionDate.After(firstOfNextMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].NextMonth -= trans.Amount
				} else {
					(*wallets)[i].NextMonth += trans.Amount
				}
			} else if trans.TransactionDate.Before(firstOfMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].LastMonth -= trans.Amount
				} else {
					(*wallets)[i].LastMonth += trans.Amount
				}
			}

		}
	}

	for _, wallet := range *wallets {
		wm.LastMonth += wallet.LastMonth
		wm.CurrentBalance += wallet.CurrentBalance + wallet.LastMonth
		wm.NextMonth += wallet.NextMonth + wallet.CurrentBalance + wallet.LastMonth
	}

	wm.Currency = "USD"
	wm.WalletId = walletId

	return wm
}

// Appends Transaction to the belonging walletId
//
// If missing, it creates the item list.
func addWhere(s *[]models.WalletTransactions, walletId string, e models.Transaction) {
	var exists bool
	for a := range *s {
		if (*s)[a].WalletId == walletId {
			(*s)[a].Transactions = append((*s)[a].Transactions, e)
			exists = true
		}
	}
	if !exists {
		var walletTransaction models.WalletTransactions
		walletTransaction.WalletId = walletId
		walletTransaction.Transactions = append(walletTransaction.Transactions, e)
		*s = append(*s, walletTransaction)
	}
}
