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

type SubscriptionService struct {
	db *pg.DB
}

func NewSubscriptionService(db *pg.DB) *SubscriptionService {
	return &SubscriptionService{
		db: db,
	}
}

/*
New

Inserts new row to subscription table.

	   	Args:
	   		context.Context: Application context
			*models.NewSubscriptionBody: Request body
		Returns:
			*models.Subscription: Created Subscription row object from database.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) New(ctx context.Context, body *models.NewSubscriptionBody) (*models.Subscription, *models.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(models.Subscription)
	exceptionReturn := new(models.Exception)

	amount, _ := body.Amount.Float64()
	customRange, _ := body.CustomRange.Int64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.SubscriptionTypeID = body.SubscriptionTypeID
	tm.CustomRange = int(customRange)
	tm.Description = body.Description
	tm.StartDate = body.StartDate
	tm.HasEnd = body.HasEnd
	tm.EndDate = body.EndDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	if body.StartDate.IsZero() {
		tm.StartDate = time.Now()
	}

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400109"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}
	tx.Commit()

	return tm, nil
}

/*
Get

Gets row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*models.Auth: Authentication model
			string: subscription id to search
			params: *models.Params
		Returns:
			*models.Subscription: Subscription row object from database.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) Get(ctx context.Context, am *models.Auth, id string, params *models.Params) (*models.Subscription, *models.Exception) {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(models.Exception)
	wm := new(models.Subscription)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	err := common.GenerateEmbed(qry, params.Embed).WherePK().Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400129"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return wm, nil
}

/*
GetAll

Gets filtered rows from subscription table.

	   	Args:
	   		context.Context: Application context
			*models.Auth: Authentication object
			string: Wallet id to search
			*models.FilteredResponse: filter options
		Returns:
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) GetAll(ctx context.Context, am *models.Auth, walletId string, filtered *models.FilteredResponse) *models.Exception {
	db := as.db.WithContext(ctx)

	wm := new([]models.Subscription)
	exceptionReturn := new(models.Exception)

	tx, _ := db.Begin()
	defer tx.Rollback()

	query := tx.Model(wm).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id)
	if walletId != "" {
		query = query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}

	err := FilteredResponse(query, wm, filtered)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400110"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"subscription\" table: %s", err)
		return exceptionReturn
	}
	tx.Commit()

	return nil
}

/*
Edit

Updates row from subscription table by id.

	   	Args:
	   		context.Context: Application context
			*models.SubscriptionEdit: Values to edit
			string: id to search
		Returns:
			*models.Subscription: Edited Subscription row object from database.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) Edit(ctx context.Context, body *models.SubscriptionEdit, id string) (*models.Subscription, *models.Exception) {
	db := as.db.WithContext(ctx)

	amount, _ := body.Amount.Float64()
	exceptionReturn := new(models.Exception)

	tm := new(models.Subscription)
	tm.Id = id
	tm.EndDate = body.EndDate
	tm.HasEnd = body.HasEnd
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.Amount = float32(math.Round(amount*100) / 100)

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400111"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
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
			*models.Subscription: Created Subscription row object from database.
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) End(ctx context.Context, id string) (*models.Subscription, *models.Exception) {
	db := as.db.WithContext(ctx)
	exceptionReturn := new(models.Exception)

	tm := new(models.Subscription)
	tm.Id = id
	tm.EndDate = time.Now()
	tm.HasEnd = true

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400112"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"subscription\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return tm, nil
}

/*
SubToTrans

Generates and Inserts new Transaction rows from the subscription model.

	   	Args:
			*models.Subscription: Subscription model to generate new transactions from
			*pg.Tx: Postgres query context
		Returns:
			*models.Exception: Exception payload.
*/
func (as *SubscriptionService) SubToTrans(subModel *models.Subscription, tx *pg.Tx) *models.Exception {
	exceptionReturn := new(models.Exception)

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	transactionStatus := new(models.TransactionStatus)
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "pending").Select()
	//tzFirstOfNextMonth := firstOfNextMonth.In(subModel.StartDate.Location())

	startDate := subModel.StartDate
	stopDate := firstOfNextMonth
	if subModel.HasEnd && subModel.EndDate.Before(firstOfNextMonth) {
		stopDate = subModel.EndDate
	}

	transactions := new([]models.Transaction)

	if subModel.SubscriptionType == nil {
		st := new(models.SubscriptionType)
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
