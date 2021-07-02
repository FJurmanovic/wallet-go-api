package services

import (
	"context"
	"math"
	"time"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

type TransactionService struct {
	Db *pg.DB
	Ss *SubscriptionService
}

func (as *TransactionService) New(ctx context.Context, body *models.NewTransactionBody) *models.Transaction {
	db := as.Db.WithContext(ctx)

	tm := new(models.Transaction)

	amount, _ := body.Amount.Float64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	if body.TransactionDate.IsZero() {
		tm.TransactionDate = time.Now()
	}

	db.Model(tm).Insert()

	return tm
}

func (as *TransactionService) GetAll(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse) {
	db := as.Db.WithContext(ctx)

	wm := new([]models.Transaction)
	sm := new([]models.Subscription)

	tx, _ := db.Begin()
	defer tx.Rollback()

	query2 := tx.Model(sm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query2 = query2.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query2.Select()

	for _, sub := range *sm {
		if sub.HasNew() {
			as.Ss.SubToTrans(&sub, tx)
		}
	}

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}

	FilteredResponse(query, wm, filtered)

	tx.Commit()
}
