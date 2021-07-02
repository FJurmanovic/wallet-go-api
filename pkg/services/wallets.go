package services

import (
	"context"
	"sync"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	Db *pg.DB
	Ss *SubscriptionService
}

func (as *WalletService) New(ctx context.Context, am *models.NewWalletBody) *models.Wallet {
	db := as.Db.WithContext(ctx)

	walletModel := new(models.Wallet)
	walletModel.Init()
	walletModel.UserID = am.UserID
	walletModel.Name = am.Name
	db.Model(walletModel).Insert()
	return walletModel
}

func (as *WalletService) Get(ctx context.Context, am *models.Auth, embed string) *models.Wallet {
	db := as.Db.WithContext(ctx)

	wm := new(models.Wallet)

	query := db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	common.GenerateEmbed(query, embed).Select()

	return wm
}

func (as *WalletService) GetAll(ctx context.Context, am *models.Auth, filtered *models.FilteredResponse) {
	db := as.Db.WithContext(ctx)
	wm := new([]models.Wallet)

	query := db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	FilteredResponse(query, wm, filtered)
}

func (as *WalletService) GetHeader(ctx context.Context, am *models.Auth, walletId string) *models.WalletHeader {
	db := as.Db.WithContext(ctx)

	wm := new(models.WalletHeader)
	wallets := new([]models.WalletTransactions)
	var wg sync.WaitGroup
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

	for _, trans := range *transactions {
		addWhere(wallets, trans.WalletID, trans)
	}

	for i, wallet := range *wallets {
		for _, trans := range wallet.Transactions {
			// tzFirstOfMonthAfterNext := firstOfMonthAfterNext.In(trans.TransactionDate.Location())
			// tzFirstOfNextMonth := firstOfNextMonth.In(trans.TransactionDate.Location())
			// tzFirstOfMonth := firstOfMonth.In(trans.TransactionDate.Location())
			if trans.TransactionDate.Before(firstOfMonth) && trans.TransactionDate.After(firstOfMonth) || trans.TransactionDate.Equal(firstOfMonth) {
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

	wg.Wait()

	for _, wallet := range *wallets {
		wm.LastMonth += wallet.LastMonth
		wm.CurrentBalance += wallet.CurrentBalance + wallet.LastMonth
		wm.NextMonth += wallet.NextMonth + wallet.CurrentBalance + wallet.LastMonth
	}

	wm.Currency = "USD"
	wm.WalletId = walletId

	tx.Commit()

	return wm
}

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
