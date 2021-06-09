package services

import (
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	Db *pg.DB
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
	//wallets := new([]models.Wallet)
	transactions := new([]models.Transaction)

	query := as.Db.Model(transactions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)

	if walletId != "" {
		query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}

	query.Select()

	currentBalance := 0
	lastMonthBalance := 0
	nextMonth := 0

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	firstOfMonthAfterNext := time.Date(currentYear, currentMonth+2, 1, 0, 0, 0, 0, currentLocation)

	for _, trans := range *transactions {
		if trans.TransactionDate.Before(firstOfNextMonth) && trans.TransactionDate.After(firstOfMonth) {
			currentBalance += trans.Amount
		} else if trans.TransactionDate.Before(firstOfMonthAfterNext) && trans.TransactionDate.After(firstOfNextMonth) {
			nextMonth += trans.Amount
		} else if trans.TransactionDate.Before(firstOfMonth) {
			lastMonthBalance += trans.Amount
			currentBalance += trans.Amount
		}
	}

	wm.CurrentBalance = currentBalance
	wm.LastMonth = lastMonthBalance
	wm.NextMonth = currentBalance + nextMonth
	wm.Currency = "USD"
	wm.WalletId = walletId

	//common.GenerateEmbed(query, embed).Select()

	return wm
}
