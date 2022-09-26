package services

import (
	"context"
	"fmt"
	"math"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionService struct {
	db                  *pg.DB
	subscriptionService *SubscriptionService
}

func NewTransactionService(db *pg.DB, ss *SubscriptionService) *TransactionService {
	return &TransactionService{
		db:                  db,
		subscriptionService: ss,
	}
}

/*
New row into transaction table

Inserts

	   	Args:
			context.Context: Application context
			*models.NewTransactionBody: Transaction body object
		Returns:
			*models.Transaction: Transaction object
			*models.Exception: Exception payload.
*/
func (as *TransactionService) New(ctx context.Context, body *models.NewTransactionBody) (*models.Transaction, *models.Exception) {
	db := as.db.WithContext(ctx)
	exceptionReturn := new(models.Exception)

	tm := new(models.Transaction)
	transactionStatus := new(models.TransactionStatus)

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "completed").Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400115"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionsStatus\" table: %s", err)
		return nil, exceptionReturn
	}

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

	_, err = tx.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400116"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}
	tx.Commit()

	return tm, nil
}

/*
GetAll

Gets all rows from subscription type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*models.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) GetAll(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse, noPending bool) *models.Exception {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(models.Exception)
	wm := new([]models.Transaction)
	transactionStatus := new(models.TransactionStatus)

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "completed").Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400117"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionStatus\" table: %s", err)
		return exceptionReturn
	}

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	if noPending {
		query = query.Where("? = ?", pg.Ident("transaction_status_id"), transactionStatus.Id)
	}

	err = FilteredResponse(query, wm, filtered)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400118"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row(s) in \"transaction\" table: %s", err)
		return exceptionReturn
	}

	tx.Commit()
	return nil
}

/*
Check

Checks subscriptions and create transacitons.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*models.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionService) Check(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse) *models.Exception {
	db := as.db.WithContext(ctx)

	wm := new([]models.Transaction)
	sm := new([]models.Subscription)
	transactionStatus := new(models.TransactionStatus)
	exceptionReturn := new(models.Exception)

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "pending").Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400119"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionStatus\" table: %s", err)
		return exceptionReturn
	}
	query2 := tx.Model(sm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query2 = query2.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query2.Select()

	for _, sub := range *sm {
		if sub.HasNew() {
			as.subscriptionService.SubToTrans(&sub, tx)
		}
	}

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query = query.Where("? = ?", pg.Ident("transaction_status_id"), transactionStatus.Id)

	err = FilteredResponse(query, wm, filtered)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400120"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transaction\" table: %s", err)
		return exceptionReturn
	}

	tx.Commit()
	return nil
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
			*models.Exception: Exception payload.
*/
func (as *TransactionService) Edit(ctx context.Context, body *models.TransactionEdit, id string) (*models.Transaction, *models.Exception) {
	db := as.db.WithContext(ctx)

	amount, _ := body.Amount.Float64()

	exceptionReturn := new(models.Exception)

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

	_, err := tx.Model(tm).WherePK().UpdateNotZero()

	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400107"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

	err = tx.Model(tm).WherePK().Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400108"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return tm, nil
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
			*models.Exception: Exception payload.
*/
func (as *TransactionService) BulkEdit(ctx context.Context, body *[]models.TransactionEdit) (*[]models.Transaction, *models.Exception) {
	db := as.db.WithContext(ctx)
	tx, _ := db.Begin()
	defer tx.Rollback()

	transactions := new([]models.Transaction)
	exceptionReturn := new(models.Exception)

	for _, transaction := range *body {

		amount, _ := transaction.Amount.Float64()

		tm := new(models.Transaction)
		tm.Id = transaction.Id
		tm.Description = transaction.Description
		tm.WalletID = transaction.WalletID
		tm.TransactionTypeID = transaction.TransactionTypeID
		tm.TransactionDate = transaction.TransactionDate
		tm.Amount = float32(math.Round(amount*100) / 100)

		*transactions = append(*transactions, *tm)
	}

	_, err := tx.Model(transactions).WherePK().UpdateNotZero()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400121"
		exceptionReturn.Message = fmt.Sprintf("Error updating rows in \"transactions\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return transactions, nil
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
			*models.Exception: Exception payload.
*/
func (as *TransactionService) Get(ctx context.Context, am *models.Auth, id string, params *models.Params) (*models.Transaction, *models.Exception) {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(models.Exception)
	wm := new(models.Transaction)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	err := common.GenerateEmbed(qry, params.Embed).WherePK().Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400122"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactions\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return wm, nil
}
