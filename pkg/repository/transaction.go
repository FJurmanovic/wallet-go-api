package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionRepository struct {
	db                          *pg.DB
	subscriptionRepository      *SubscriptionRepository
	transactionStatusRepository *TransactionStatusRepository
}

func NewTransactionRepository(db *pg.DB, ss *SubscriptionRepository, tsr *TransactionStatusRepository) *TransactionRepository {
	return &TransactionRepository{
		db:                          db,
		subscriptionRepository:      ss,
		transactionStatusRepository: tsr,
	}
}

/*
New row into transaction table

Inserts

	   	Args:
			context.Context: Application context
			*model.NewTransactionBody: Transaction body object
		Returns:
			*model.Transaction: Transaction object
			*model.Exception: Exception payload.
*/
func (as *TransactionRepository) New(ctx context.Context, tm *model.Transaction, tx *pg.Tx) (*model.Transaction, error) {
	var commit = false
	if tx == nil {
		commit = true
		db := as.db.WithContext(ctx)
		tx, _ = db.Begin()

	}

	if commit {
		defer tx.Rollback()
	}

	_, err := tx.Model(tm).Insert()
	if err != nil {
		return nil, err
	}
	if commit {
		tx.Commit()
	}

	return tm, nil
}

/*
GetAll

Gets all rows from subscription type table.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*model.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionRepository) GetAll(ctx context.Context, flt *filter.TransactionFilter) (*model.FilteredResponse, *model.Exception) {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(model.Exception)
	wm := new([]model.Transaction)

	tx, _ := db.Begin()
	defer tx.Rollback()

	if flt.NoPending {
		tsFlt := filter.NewTransactionStatusFilter(model.Params{})
		tsFlt.Status = "completed"
		transactionStatus, err := as.transactionStatusRepository.GetTx(tx, tsFlt)

		if err != nil {
			exceptionReturn.StatusCode = 400
			exceptionReturn.ErrorCode = "400117"
			exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionStatus\" table: %s", err)
			return nil, exceptionReturn
		}

		flt.TransactionStatusId = transactionStatus.Id
	}

	query := tx.Model(wm)

	as.OnBeforeGetTransactionFilter(query, flt)
	filtered, err := FilteredResponse(query, wm, flt.Params)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400118"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row(s) in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()
	return filtered, nil
}

/*
GetAllTx

Gets filtered rows from transaction table.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication object
			string: Wallet id to search
			*model.FilteredResponse: filter options
		Returns:
			*model.Exception: Exception payload.
*/
func (as *TransactionRepository) GetAllTx(tx *pg.Tx, flt *filter.TransactionFilter) (*[]model.Transaction, error) {
	am := new([]model.Transaction)
	query := tx.Model(am)
	as.OnBeforeGetTransactionFilter(query, flt)

	common.GenerateEmbed(query, flt.Embed)
	err := query.Select()
	if err != nil {
		return nil, err
	}

	return am, nil
}

/*
Check

Checks subscriptions and create transactions.
   	Args:
		context.Context: Application context
		string: Relations to embed
	Returns:
		*model.Exception: Exception payload.
*/
// Gets filtered rows from transaction table.
func (as *TransactionRepository) Check(ctx context.Context, flt *filter.TransactionFilter) (*model.FilteredResponse, *model.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]model.Transaction)
	exceptionReturn := new(model.Exception)
	filtered := new(model.FilteredResponse)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tsFlt := filter.NewTransactionStatusFilter(model.Params{})
	tsFlt.Status = "pending"
	transactionStatus, err := as.transactionStatusRepository.GetTx(tx, tsFlt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400119"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}
	flt.TransactionStatusId = transactionStatus.Id

	smFlt := filter.NewSubscriptionFilter(model.Params{})
	smFlt.Id = flt.Id
	smFlt.WalletId = flt.WalletId
	sm, err := as.subscriptionRepository.GetAllTx(tx, smFlt)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400137"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	for _, sub := range *sm {
		if sub.HasNew() {
			as.subscriptionRepository.SubToTrans(&sub, tx)
		}
	}

	qry := tx.Model(wm)
	as.OnBeforeGetTransactionFilter(qry, flt)
	filtered, err = FilteredResponse(qry, wm, flt.Params)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400120"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transaction\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()
	return filtered, nil
}

/*
Edit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			*model.TransactionEdit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionRepository) Edit(ctx context.Context, tm *model.Transaction) (*model.Transaction, error) {
	db := as.db.WithContext(ctx)

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()

	if err != nil {
		return nil, err
	}

	err = tx.Model(tm).WherePK().Select()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return tm, nil
}

/*
BulkEdit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			?[]model.Transaction Bulk Edit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionRepository) BulkEdit(ctx context.Context, transactions *[]model.Transaction) (*[]model.Transaction, error) {
	db := as.db.WithContext(ctx)
	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(transactions).WherePK().UpdateNotZero()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return transactions, nil
}

/*
Get

Gets row from transaction table by id.

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
			string: id to search
			*model.Params: url query parameters
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionRepository) Get(ctx context.Context, flt *filter.TransactionFilter) (*model.Transaction, error) {
	db := as.db.WithContext(ctx)
	wm := new(model.Transaction)
	wm.Id = flt.Id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	as.OnBeforeGetTransactionFilter(qry, flt)
	err := common.GenerateEmbed(qry, flt.Embed).WherePK().Select()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return wm, nil
}

func (as *TransactionRepository) OnBeforeGetTransactionFilter(qry *orm.Query, flt *filter.TransactionFilter) {
	if flt.WalletId != "" {
		qry.Where("? = ?", pg.Ident("wallet_id"), flt.WalletId)
	}
	if flt.UserId != "" {
		qry.Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), flt.UserId)
	}
	if flt.TransactionStatusId != "" {
		qry.Where("? = ?", pg.Ident("transaction_status_id"), flt.TransactionStatusId)
	}
}

func (as *TransactionRepository) CreateTx(ctx context.Context) (*pg.Tx, error) {
	db := as.db.WithContext(ctx)
	return db.Begin()
}
