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
func (as *TransactionService) GetAll(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse, noPending bool) {
	db := as.Db.WithContext(ctx)

	wm := new([]models.Transaction)
	transactionStatus := new(models.TransactionStatus)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "completed").Select()

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	if noPending {
		query = query.Where("? = ?", pg.Ident("transaction_status_id"), transactionStatus.Id)
	}

	FilteredResponse(query, wm, filtered)

	tx.Commit()
}

/*
Check

Checks subscriptions and create transacitons.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*[]models.SubscriptionType: List of subscription type objects.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) Check(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse) {
	db := as.Db.WithContext(ctx)

	wm := new([]models.Transaction)
	sm := new([]models.Subscription)
	transactionStatus := new(models.TransactionStatus)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "pending").Select()
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
	query = query.Where("? = ?", pg.Ident("transaction_status_id"), transactionStatus.Id)

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
	tm.TransactionStatusID = body.TransactionStatusID
	tm.Amount = float32(math.Round(amount*100) / 100)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(tm).WherePK().UpdateNotZero()

	tx.Commit()

	return tm
}

/*
Bulk Edit

Updates row in transaction table by id.
   	Args:
		context.Context: Application context
		?[]models.Transaction Bulk Edit: Object to edit
		string: id to search
	Returns:
		*models.Transaction: Transaction object from database.
*/
func (as *TransactionService) BulkEdit(ctx context.Context, body *[]models.TransactionEdit) *[]models.Transaction {
	db := as.Db.WithContext(ctx)
	tx, _ := db.Begin()
	defer tx.Rollback()

	transactions := new([]models.Transaction)

	for _, transaction := range *body {

		amount, _ := transaction.Amount.Float64()

		tm := new(models.Transaction)
		tm.Id = transaction.Id
		tm.Description = transaction.Description
		tm.WalletID = transaction.WalletID
		tm.TransactionTypeID = transaction.TransactionTypeID
		tm.TransactionDate = transaction.TransactionDate
		tm.Amount = float32(math.Round(amount*100) / 100)

		tx.Model(tm).WherePK().UpdateNotZero()
		*transactions = append(*transactions, *tm)
	}

	tx.Commit()

	return transactions
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
