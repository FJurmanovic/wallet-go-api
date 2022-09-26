package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"time"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type SubscriptionRepository struct {
	db *pg.DB
}

func NewSubscriptionRepository(db *pg.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

/*
New

Inserts new row to subscription table.

	   	Args:
	   		context.Context: Application context
			*model.NewSubscriptionBody: Request body
		Returns:
			*model.Subscription: Created Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) New(ctx context.Context, tm *model.Subscription) (*model.Subscription, error) {
	db := as.db.WithContext(ctx)

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).Insert()
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return tm, nil
}

/*
Get

Gets row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication model
			string: subscription id to search
			params: *model.Params
		Returns:
			*model.Subscription: Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) Get(ctx context.Context, am *model.Subscription, flt filter.SubscriptionFilter) (*model.Subscription, error) {
	db := as.db.WithContext(ctx)
	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(am)
	as.OnBeforeGetSubscriptionFilter(qry, &flt)
	err := common.GenerateEmbed(qry, flt.Embed).WherePK().Select()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return am, nil
}

/*
GetAll

Gets filtered rows from subscription table.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication object
			string: Wallet id to search
			*model.FilteredResponse: filter options
		Returns:
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) GetAll(ctx context.Context, am *[]model.Subscription, filtered *model.FilteredResponse, flt *filter.SubscriptionFilter) error {
	db := as.db.WithContext(ctx)

	tx, _ := db.Begin()
	defer tx.Rollback()

	query := tx.Model(am)
	as.OnBeforeGetSubscriptionFilter(query, flt)

	err := FilteredResponse(query, am, filtered)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

/*
GetAllTx

Gets filtered rows from subscription table.

	   	Args:
	   		context.Context: Application context
			*model.Auth: Authentication object
			string: Wallet id to search
			*model.FilteredResponse: filter options
		Returns:
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) GetAllTx(tx *pg.Tx, am *[]model.Subscription, flt *filter.SubscriptionFilter) error {
	query := tx.Model(am)
	as.OnBeforeGetSubscriptionFilter(query, flt)

	err := query.Select()
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

/*
Edit

Updates row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*model.SubscriptionEdit: Values to edit
			string: id to search
		Returns:
			*model.Subscription: Edited Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) Edit(ctx context.Context, tm *model.Subscription) (*model.Subscription, error) {
	db := as.db.WithContext(ctx)

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return tm, nil
}

/*
End

Updates row in subscription table by id.

Ends subscription with current date.

	   	Args:
	   		context.Context: Application context
			string: id to search
		Returns:
			*model.Subscription: Created Subscription row object from database.
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) End(ctx context.Context, tm *model.Subscription) (*model.Subscription, error) {
	db := as.db.WithContext(ctx)

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return tm, nil
}

/*
SubToTrans

Generates and Inserts new Transaction rows from the subscription model.

	   	Args:
			*model.Subscription: Subscription model to generate new transactions from
			*pg.Tx: Postgres query context
		Returns:
			*model.Exception: Exception payload.
*/
func (as *SubscriptionRepository) SubToTrans(subModel *model.Subscription, tx *pg.Tx) *model.Exception {
	exceptionReturn := new(model.Exception)

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	transactionStatus := new(model.TransactionStatus)
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "pending").Select()
	//tzFirstOfNextMonth := firstOfNextMonth.In(subModel.StartDate.Location())

	startDate := subModel.StartDate
	stopDate := firstOfNextMonth
	if subModel.HasEnd && subModel.EndDate.Before(firstOfNextMonth) {
		stopDate = subModel.EndDate
	}

	transactions := new([]model.Transaction)

	if subModel.SubscriptionType == nil {
		st := new(model.SubscriptionType)
		tx.Model(st).Where("? = ?", pg.Ident("id"), subModel.SubscriptionTypeID).Select()
		subModel.SubscriptionType = st
	}

	for startDate.Before(stopDate) {
		trans := subModel.ToTrans()
		trans.TransactionDate = startDate
		trans.TransactionStatusID = transactionStatus.Id
		if startDate.After(subModel.LastTransactionDate) && (startDate.Before(now) || startDate.Equal(now)) {
			*transactions = append(*transactions, *trans)
		}
		if subModel.SubscriptionType.Type == "monthly" {
			startDate = startDate.AddDate(0, subModel.CustomRange, 0)
		} else if subModel.SubscriptionType.Type == "weekly" {
			startDate = startDate.AddDate(0, 0, 7*subModel.CustomRange)
		} else if subModel.SubscriptionType.Type == "daily" {
			startDate = startDate.AddDate(0, 0, subModel.CustomRange)
		} else {
			startDate = startDate.AddDate(subModel.CustomRange, 0, 0)
		}
	}

	var err error
	if len(*transactions) > 0 {
		for _, trans := range *transactions {
			_, err = tx.Model(&trans).Where("? = ?", pg.Ident("transaction_date"), trans.TransactionDate).Where("? = ?", pg.Ident("subscription_id"), trans.SubscriptionID).OnConflict("DO NOTHING").SelectOrInsert()
			if err != nil {
				_, err = tx.Model(subModel).Set("? = ?", pg.Ident("last_transaction_date"), trans.TransactionDate).WherePK().Update()
			}
		}
	}
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400113"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"subscription\" table: %s", err)
		return exceptionReturn
	}
	return nil
}

func (as *SubscriptionRepository) OnBeforeGetSubscriptionFilter(qry *orm.Query, flt *filter.SubscriptionFilter) {
	if flt.Id != "" {
		qry.Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), flt.Id)
	}
	if flt.WalletId != "" {
		qry.Where("? = ?", pg.Ident("wallet_id"), flt.WalletId)
	}
}
