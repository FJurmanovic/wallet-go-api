package services

import (
	"time"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type TransactionService struct {
	Db *pg.DB
}

func (as *TransactionService) New(body *models.NewTransactionBody) *models.Transaction {
	tm := new(models.Transaction)

	amount, _ := body.Amount.Int64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate
	tm.Amount = int(amount)

	if body.TransactionDate.IsZero() {
		tm.TransactionDate = time.Now()
	}

	as.Db.Model(tm).Insert()

	return tm
}

func (as *TransactionService) GetAll(am *models.Auth, walletId string, filtered *models.FilteredResponse) {
	wm := new([]models.Transaction)

	query := as.Db.Model((wm)).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	FilteredResponse(query, wm, filtered)
}
