package services

import (
	"context"
	"math"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionService struct {
	Db *pg.DB
	Ss *SubscriptionService
}

func (as *TransactionService) New(ctx context.Context, body *models.NewTransactionBody) *models.Transaction {
	db := as.Db.WithContext(ctx)

	tm := new(models.Transaction)

	tx, _ := db.Begin()
	defer tx.Rollback()

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

	tx.Model(tm).Insert()
	tx.Commit()

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

func (as *TransactionService) Edit(ctx context.Context, body *models.TransactionEdit, id string) *models.Transaction {
	db := as.Db.WithContext(ctx)

	amount, _ := body.Amount.Float64()

	tm := new(models.Transaction)
	tm.Id = id
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.TransactionDate = body.TransactionDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(tm).WherePK().UpdateNotZero()

	tx.Commit()

	return tm
}

func (as *TransactionService) Get(ctx context.Context, am *models.Auth, id string, params *models.Params) *models.Transaction {
	db := as.Db.WithContext(ctx)

	wm := new(models.Transaction)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	common.GenerateEmbed(qry, params.Embed).WherePK().Select()

	tx.Commit()

	return wm
}
