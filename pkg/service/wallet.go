package service

import (
	"context"
	"fmt"
	"time"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	db                  *pg.DB
	subscriptionService *SubscriptionService
}

func NewWalletService(db *pg.DB, ss *SubscriptionService) *WalletService {
	return &WalletService{
		db:                  db,
		subscriptionService: ss,
	}
}

/*
New

Inserts row to wallets table.

	   	Args:
			context.Context: Application context
			*model.NewWalletBody: Object to be inserted
		Returns:
			*model.Wallet: Wallet object from database.
			*model.Exception: Exception payload.
*/
func (as *WalletService) New(ctx context.Context, am *model.NewWalletBody) (*model.Wallet, *model.Exception) {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(model.Exception)
	walletModel := new(model.Wallet)
	walletModel.Init()
	walletModel.UserID = am.UserID
	walletModel.Name = am.Name
	_, err := db.Model(walletModel).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400126"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"wallets\" table: %s", err)
		return nil, exceptionReturn
	}
	return walletModel, nil
}

/*
Edit

Updates row in wallets table by id.

	   	Args:
			context.Context: Application context
			*model.WalletEdit: Object to be edited
			string: id to search
		Returns:
			*model.Wallet: Wallet object from database.
			*model.Exception: Exception payload.
*/
func (as *WalletService) Edit(ctx context.Context, body *model.WalletEdit, id string) (*model.Wallet, *model.Exception) {
	db := as.db.WithContext(ctx)

	exceptionReturn := new(model.Exception)
	tm := new(model.Wallet)
	tm.Id = id
	tm.Name = body.Name

	tx, _ := db.Begin()
	defer tx.Rollback()

	_, err := tx.Model(tm).WherePK().UpdateNotZero()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400127"
		exceptionReturn.Message = fmt.Sprintf("Error updating row in \"wallets\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return tm, nil
}

/*
Get

Gets row in wallets table by id.

	   	Args:
			context.Context: Application context
			string: id to search
			*model.Params: url query parameters
		Returns:
			*model.Wallet: Wallet object from database
			*model.Exception: Exception payload.
*/
func (as *WalletService) Get(ctx context.Context, id string, params *model.Params) (*model.Wallet, *model.Exception) {
	db := as.db.WithContext(ctx)
	exceptionReturn := new(model.Exception)

	wm := new(model.Wallet)
	wm.Id = id

	tx, _ := db.Begin()
	defer tx.Rollback()

	qry := tx.Model(wm)
	err := common.GenerateEmbed(qry, params.Embed).WherePK().Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400128"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"wallets\" table: %s", err)
		return nil, exceptionReturn
	}

	tx.Commit()

	return wm, nil
}

/*
GetAll

Gets filtered rows from wallets table.

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
			*model.FilteredResponse: filter options
		Returns:
			*model.Exception: Exception payload.
*/
func (as *WalletService) GetAll(ctx context.Context, am *model.Auth, filtered *model.FilteredResponse) *model.Exception {
	exceptionReturn := new(model.Exception)
	db := as.db.WithContext(ctx)
	wm := new([]model.Wallet)

	query := db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	err := FilteredResponse(query, wm, filtered)
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400134"
		exceptionReturn.Message = fmt.Sprintf("Error selecting rows in \"wallets\" table: %s", err)
		return exceptionReturn
	}
	return nil
}

/*
GetHeader

Gets row from wallets table.

Calculates previous month, current and next month totals.

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
			string: wallet id to search
		Returns:
			*model.WalletHeader: generated wallet header object
			*model.Exception: Exception payload.
*/
func (as *WalletService) GetHeader(ctx context.Context, am *model.Auth, walletId string) (*model.WalletHeader, *model.Exception) {
	db := as.db.WithContext(ctx)

	wm := new(model.WalletHeader)
	wallets := new([]model.WalletTransactions)
	transactions := new([]model.Transaction)
	subscriptions := new([]model.Subscription)
	transactionStatus := new(model.TransactionStatus)
	exceptionReturn := new(model.Exception)

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := tx.Model(transactionStatus).Where("? = ?", pg.Ident("status"), "completed").Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400130"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionStatuses\" table: %s", err)
		return nil, exceptionReturn
	}
	query2 := tx.Model(subscriptions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType").Relation("SubscriptionType")
	if walletId != "" {
		query2.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query2.Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400131"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"subscriptions\" table: %s", err)
		return nil, exceptionReturn
	}

	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, currentLocation)
	firstOfMonthAfterNext := time.Date(currentYear, currentMonth+2, 1, 0, 0, 0, 0, currentLocation)

	query := tx.Model(transactions).Relation("Wallet").Where("wallet.? = ?", pg.Ident("user_id"), am.Id).Relation("TransactionType")
	if walletId != "" {
		query.Where("? = ?", pg.Ident("wallet_id"), walletId)
	}
	query = query.Where("? = ?", pg.Ident("transaction_status_id"), transactionStatus.Id)
	query.Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400132"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactions\" table: %s", err)
		return nil, exceptionReturn
	}
	tx.Commit()

	for _, sub := range *subscriptions {
		stopDate := firstOfMonthAfterNext
		if sub.HasEnd && sub.EndDate.Before(firstOfMonthAfterNext) {
			stopDate = sub.EndDate
		}
		startDate := sub.StartDate
		for startDate.Before(stopDate) {
			trans := sub.ToTrans()
			trans.TransactionDate = startDate
			if startDate.After(firstOfNextMonth) || startDate.Equal(firstOfNextMonth) {
				*transactions = append(*transactions, *trans)
			}
			if sub.SubscriptionType.Type == "monthly" {
				startDate = startDate.AddDate(0, sub.CustomRange, 0)
			} else if sub.SubscriptionType.Type == "weekly" {
				startDate = startDate.AddDate(0, 0, 7*sub.CustomRange)
			} else if sub.SubscriptionType.Type == "daily" {
				startDate = startDate.AddDate(0, 0, sub.CustomRange)
			} else {
				startDate = startDate.AddDate(sub.CustomRange, 0, 0)
			}
		}
	}

	for _, trans := range *transactions {
		addWhere(wallets, trans.WalletID, trans)
	}

	for i, wallet := range *wallets {
		for _, trans := range wallet.Transactions {
			// tzFirstOfMonthAfterNext := firstOfMonthAfterNext.In(trans.TransactionDate.Location())
			// tzFirstOfNextMonth := firstOfNextMonth.In(trans.TransactionDate.Location())
			// tzFirstOfMonth := firstOfMonth.In(trans.TransactionDate.Location())
			if trans.TransactionDate.Before(firstOfNextMonth) && trans.TransactionDate.After(firstOfMonth) || trans.TransactionDate.Equal(firstOfMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].CurrentBalance -= trans.Amount
				} else {
					(*wallets)[i].CurrentBalance += trans.Amount
				}
			} else if trans.TransactionDate.Before(firstOfMonthAfterNext) && trans.TransactionDate.After(firstOfNextMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].NextMonth -= trans.Amount
				} else {
					(*wallets)[i].NextMonth += trans.Amount
				}
			} else if trans.TransactionDate.Before(firstOfMonth) {
				if trans.TransactionType.Type == "expense" {
					(*wallets)[i].LastMonth -= trans.Amount
				} else {
					(*wallets)[i].LastMonth += trans.Amount
				}
			}

		}
	}

	for _, wallet := range *wallets {
		wm.LastMonth += wallet.LastMonth
		wm.CurrentBalance += wallet.CurrentBalance + wallet.LastMonth
		wm.NextMonth += wallet.NextMonth + wallet.CurrentBalance + wallet.LastMonth
	}

	wm.Currency = "USD"
	wm.WalletId = walletId

	return wm, nil
}

/*
addWhere

Appends Transaction to the belonging walletId.

If missing, it creates the item list.

	   	Args:
			*[]model.WalletTransactions: list to append to
			string: wallet id to check
			model.Transaction: Transaction to append
		Returns:
			*model.Exception: Exception payload.
*/
func addWhere(s *[]model.WalletTransactions, walletId string, e model.Transaction) {
	var exists bool
	for a := range *s {
		if (*s)[a].WalletId == walletId {
			(*s)[a].Transactions = append((*s)[a].Transactions, e)
			exists = true
		}
	}
	if !exists {
		var walletTransaction model.WalletTransactions
		walletTransaction.WalletId = walletId
		walletTransaction.Transactions = append(walletTransaction.Transactions, e)
		*s = append(*s, walletTransaction)
	}
}
