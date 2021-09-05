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

/*
New new row into transaction table

Inserts
   	Args:
		context.Context: Application context
		*models.NewTransactionBody: Transaction body object
	Returns:
		*models.Transaction: Transaction object
*/
func (as *TransactionService) New(ctx context.Context, body *models.NewTransactionBody) *models.Transaction {
	db := as.Db.WithContext(ctx)

	tm := new(models.Transaction)
	transactionStatus := new(models.TransactionStatus)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "completed").Select()

	amount, _ := body.Amount.Float64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate
	tm.Amount = float32(math.Round(amount*100) / 100)
	tm.TransactionStatusID = transactionStatus.Id

	if body.TransactionDate.IsZero() {
		tm.TransactionDate = time.Now()
	}

	tx.Model(tm).Insert()
	tx.Commit()

	return tm
}

/*
GetAll

Gets all rows from subscription type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.SubscriptionType: List of subscription type objects.
*/
// Gets filtered rows from transaction table.
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

/*
Edit

Updates row in transaction table by id.
   	Args:
		context.Context: Application context
		*models.TransactionEdit: Object to edit
		string: id to search
	Returns:
		*models.Transaction: Transaction object from database.
*/
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

/*
Get

Gets row from transaction table by id.
   	Args:
		context.Context: Application context
		*models.Auth: Authentication object
		string: id to search
		*model.Params: url query parameters
	Returns:
		*models.Transaction: Transaction object from database.
*/
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
