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

type TransactionStatusRepository struct {
	db *pg.DB
}

func NewTransactionStatusRepository(db *pg.DB) *TransactionStatusRepository {
	return &TransactionStatusRepository{
		db: db,
	}
}

/*
New

Inserts new row to transaction status table.

	   	Args:
			context.Context: Application context
			*model.NewTransactionStatusBody: object to create
		Returns:
			*model.TransactionType: Transaction Type object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusRepository) New(ctx context.Context, body *model.NewTransactionStatusBody) (*model.TransactionStatus, *model.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(model.TransactionStatus)
	exceptionReturn := new(model.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Status = body.Status

	_, err := db.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400123"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return tm, nil
}

/*
GetAll

Gets all rows from transaction status table.

	   	Args:
			context.Context: Application context
			string: Relations to embed
		Returns:
			*[]model.TransactionStatus: List of Transaction status objects from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusRepository) GetAll(ctx context.Context, embed string) (*[]model.TransactionStatus, *model.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]model.TransactionStatus)
	exceptionReturn := new(model.Exception)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, embed).Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400124"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"transactionStatus\" table: %s", err)
		return nil, exceptionReturn
	}

	return wm, nil
}

/*
Get

Gets row from transactionStatus table by id.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication model
			string: transactionStatus id to search
			params: *model.Params
		Returns:
			*model.Subscription: Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusRepository) Get(ctx context.Context, am *model.TransactionStatus, flt filter.TransactionStatusFilter) (*model.TransactionStatus, error) {
	db := as.db.WithContext(ctx)
	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(am)
	err := common.GenerateEmbed(qry, flt.Embed).WherePK().Select()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return am, nil
}

/*
GetTx

Gets row from transactionStatus table by id.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication model
			string: transactionStatus id to search
			params: *model.Params
		Returns:
			*model.Subscription: Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *TransactionStatusRepository) GetTx(tx *pg.Tx, am *model.TransactionStatus, flt *filter.TransactionStatusFilter) (*model.TransactionStatus, error) {
	qry := tx.Model(am)
	as.OnBeforeGetTransactionStatusFilter(qry, flt)
	err := common.GenerateEmbed(qry, flt.Embed).WherePK().Select()
	if err != nil {
		return nil, err
	}

	return am, nil
}

func (as *TransactionStatusRepository) OnBeforeGetTransactionStatusFilter(qry *orm.Query, flt *filter.TransactionStatusFilter) {
	if flt.Status != "" {
		qry.Where("? = ?", pg.Ident("status"), flt.Status)
	}
}
