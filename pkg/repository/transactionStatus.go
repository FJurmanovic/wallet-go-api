package repository

import (
	"context"
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
func (as *TransactionStatusRepository) New(ctx context.Context, tm *model.TransactionStatus) (*model.TransactionStatus, error) {
	db := as.db.WithContext(ctx)

	_, err := db.Model(tm).Insert()
	if err != nil {
		return nil, err
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
func (as *TransactionStatusRepository) GetAll(ctx context.Context, flt *filter.TransactionStatusFilter, tx *pg.Tx) (*[]model.TransactionStatus, error) {
	var commit = false
	if tx == nil {
		commit = true
		db := as.db.WithContext(ctx)
		tx, _ = db.Begin()
	}

	if commit {
		defer tx.Rollback()
	}

	wm := new([]model.TransactionStatus)

	query := tx.Model(wm)
	as.OnBeforeGetTransactionStatusFilter(query, flt)
	err := common.GenerateEmbed(query, flt.Embed).Select()
	if err != nil {
		return nil, err
	}

	if commit {
		tx.Commit()
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
func (as *TransactionStatusRepository) Get(ctx context.Context, flt *filter.TransactionStatusFilter, tx *pg.Tx) (*model.TransactionStatus, error) {
	am := new(model.TransactionStatus)
	commit := false
	if tx == nil {
		commit = true
		db := as.db.WithContext(ctx)
		tx, _ = db.Begin()
		defer tx.Rollback()
	}

	qry := tx.Model(am)
	err := common.GenerateEmbed(qry, flt.Embed).Select()
	if err != nil {
		return nil, err
	}

	if commit {
		tx.Commit()
	}

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
func (as *TransactionStatusRepository) GetTx(tx *pg.Tx, flt *filter.TransactionStatusFilter) (*model.TransactionStatus, error) {
	am := new(model.TransactionStatus)
	qry := tx.Model(am)
	as.OnBeforeGetTransactionStatusFilter(qry, flt)
	err := common.GenerateEmbed(qry, flt.Embed).Select()
	if err != nil {
		return nil, err
	}

	return am, nil
}

func (as *TransactionStatusRepository) OnBeforeGetTransactionStatusFilter(qry *orm.Query, flt *filter.TransactionStatusFilter) {
	if flt.Id != "" {
		qry.Where("? = ?", pg.Ident("id"), flt.Id)
	}
	if flt.Status != "" {
		qry.Where("? = ?", pg.Ident("status"), flt.Status)
	}
}
