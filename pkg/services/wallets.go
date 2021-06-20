package services

import (
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

func (as *WalletService) New(am *models.NewWalletBody) *models.Wallet {

	walletModel := new(models.Wallet)
	walletModel.Init()
	walletModel.UserID = am.UserID
	walletModel.Name = am.Name
	as.Db.Model(walletModel).Insert()
	return walletModel
}

func (as *WalletService) Get(am *models.Auth, embed string) *models.Wallet {
	wm := new(models.Wallet)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	common.GenerateEmbed(query, embed).Select()

	return wm
}

func (as *WalletService) GetAll(am *models.Auth, filtered *models.FilteredResponse) {
	wm := new([]models.Wallet)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	FilteredResponse(query, wm, filtered)
}

func (as *WalletService) GetHeader(am *models.Auth, embed string, walletId string) *models.WalletHeader {
	wm := new(models.WalletHeader)
	wallets := new([]models.WalletTransactions)
	var wg sync.WaitGroup
	transactions := new([]models.Transaction)
	subscriptions := new([]models.Subscription)

	query2 := as.Db.Model(subscriptions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType").Relation("SubscriptionType")
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
		as.Ss.SubToTrans(&sub)
	}

	query := as.Db.Model(transactions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType")
	if walletId != "" {
		query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query.Select()

	for _, trans := range *transactions {
		addWhere(wallets, trans.WalletID, trans)
	}

	for i, wallet := range *wallets {
		for _, trans := range wallet.Transactions {
			if trans.TransactionDate.Local().Before(firstOfNextMonth) && trans.TransactionDate.Local().After(firstOfMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].CurrentBalance -= trans.Amount
				} else {
					(*wallets)[i].CurrentBalance += trans.Amount
				}
			} else if trans.TransactionDate.Local().Before(firstOfMonthAfterNext) && trans.TransactionDate.Local().After(firstOfNextMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].NextMonth -= trans.Amount
				} else {
					(*wallets)[i].NextMonth += trans.Amount
				}
			} else if trans.TransactionDate.Local().Before(firstOfMonth) {
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

	return wm
}

func addWhere(s *[]models.WalletTransactions, walletId string, e models.Transaction) {
	var exists bool
	for a, _ := range *s {
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
