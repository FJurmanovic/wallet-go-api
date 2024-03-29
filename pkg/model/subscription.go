package model

import (
	"encoding/json"
	"math"
	"time"
)

type Subscription struct {
	tableName struct{} `pg:"subscriptions,alias:subscriptions"`
	BaseModel
	Description         string            `json:"description" pg:"description"`
	StartDate           time.Time         `json:"startDate" pg:"start_date", type:timestamptz`
	EndDate             time.Time         `json:"endDate" pg:"end_date", type:timestamptz`
	HasEnd              bool              `json:"hasEnd" pg:"hasEnd"`
	SubscriptionTypeID  string            `json:"subscriptionTypeId" pg:"subscription_type_id"`
	SubscriptionType    *SubscriptionType `json:"subscriptionType", pg:"rel:has-one, fk:subscription_type_id"`
	CustomRange         int               `json:"customRange", pg:"custom_range,default:1"`
	WalletID            string            `json:"walletId", pg:"wallet_id"`
	Wallet              *Wallet           `json:"wallet" pg:"rel:has-one, fk:wallet_id"`
	TransactionTypeID   string            `json:"transactionTypeId", pg:"transaction_type_id"`
	TransactionType     *TransactionType  `json:"transactionType", pg:"rel:has-one, fk:transaction_type_id"`
	LastTransactionDate time.Time         `json:"lastTransactionDate", pg:"last_transaction_date", type:timestamptz`
	Amount              float32           `json:"amount", pg:"amount,default:0"`
}

type SubscriptionEdit struct {
	tableName   struct{}    `pg:"subscriptions,alias:subscriptions"`
	Id          string      `json:"id" form:"id"`
	Description string      `json:"description" form:"description"`
	EndDate     time.Time   `json:"endDate" form:"endDate" `
	HasEnd      bool        `json:"hasEnd" form:"hasEnd"`
	WalletID    string      `json:"walletId" form:"walletId"`
	Amount      json.Number `json:"amount" form:"amount"`
}

type NewSubscriptionBody struct {
	WalletID           string      `json:"walletId" form:"walletId"`
	TransactionTypeID  string      `json:"transactionTypeId" form:"transactionTypeId"`
	SubscriptionTypeID string      `json:"subscriptionTypeId" pg:"subscription_type_id"`
	CustomRange        json.Number `json:"customRange", pg:"custom_range"`
	StartDate          time.Time   `json:"startDate" pg:"start_date"`
	EndDate            time.Time   `json:"endDate" pg:"end_date"`
	HasEnd             bool        `json:"hasEnd" pg:"hasEnd"`
	Description        string      `json:"description" form:"description"`
	Amount             json.Number `json:"amount" form:"amount"`
}

func (body *SubscriptionEdit) ToSubscription() *Subscription {
	amount, _ := body.Amount.Float64()
	tm := new(Subscription)
	tm.Id = body.Id
	tm.EndDate = body.EndDate
	tm.HasEnd = body.HasEnd
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.Amount = float32(math.Round(amount*100) / 100)

	return tm
}

func (body *NewSubscriptionBody) ToSubscription() *Subscription {
	tm := new(Subscription)
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
	return tm
}

type SubscriptionEnd struct {
	Id string `json:"id" form:"id"`
}

/*
ToTrans

Maps Subscription object to Transaction object.

	Returns:
		*Transaction: mapped Transaction object
*/
func (cm *Subscription) ToTrans() *Transaction {
	trans := new(Transaction)
	trans.Init()
	trans.Amount = cm.Amount
	trans.Description = cm.Description
	trans.WalletID = cm.WalletID
	trans.Wallet = cm.Wallet
	trans.TransactionTypeID = cm.TransactionTypeID
	trans.TransactionType = cm.TransactionType
	trans.DateCreated = cm.DateCreated
	trans.SubscriptionID = cm.Id
	return trans
}

/*
HasNew

Checks if Subscription reached new transaction interval.

	Returns:
		bool: Is new transaction interval reached
*/
func (cm *Subscription) HasNew() bool {
	trans := cm.TransactionType
	if trans != nil {
		switch trans.Type {
		case "monthly":
			lastDate := time.Now().AddDate(0, -cm.CustomRange, 0)
			if cm.LastTransactionDate.Before(lastDate) {
				return true
			}
			return false
		case "weekly":
			lastDate := time.Now().AddDate(0, 0, -(7 * cm.CustomRange))
			if cm.LastTransactionDate.Before(lastDate) {
				return true
			}
			return false
		case "daily":
			lastDate := time.Now().AddDate(0, 0, -cm.CustomRange)
			if cm.LastTransactionDate.Before(lastDate) {
				return true
			}
			return false
		default:
			lastDate := time.Now().AddDate(-cm.CustomRange, 0, 0)
			if cm.LastTransactionDate.Before(lastDate) {
				return true
			}
			return false
		}
	}
	return true
}
