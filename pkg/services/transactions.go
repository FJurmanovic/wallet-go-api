package services

import (
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type TransactionService struct {
	Db *pg.DB
}

func (as *TransactionService) New(body *models.NewTransactionBody) *models.Transaction {
	tm := new(models.Transaction)

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate

	as.Db.Model(tm).Insert()

	return tm
}

func (as *TransactionService) GetAll(walletId string, filtered *models.FilteredResponse) {
	wm := new([]models.Transaction)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("wallet_id"), walletId)
	FilteredResponse(query, wm, filtered)
}
