package models

import (
	"encoding/json"
	"time"
)

type Subscription struct {
	tableName struct{} `pg:"subscriptions,alias:subscriptions"`
	BaseModel
	Description         string            `json:"description" pg:"description"`
	StartDate           time.Time         `json:"startDate" pg:"start_date"`
	SubscriptionTypeID  string            `json:"subscriptionTypeId" pg:"subscription_type_id"`
	SubscriptionType    *SubscriptionType `json:"subscriptionType", pg:"rel:has-one, fk:subscription_type_id"`
	CustomRange         int               `json:"customRange", pg:"custom_range"`
	WalletID            string            `json:"walletId", pg:"wallet_id"`
	Wallet              *Wallet           `json:"wallet" pg:"rel:has-one, fk:wallet_id"`
	TransactionTypeID   string            `json:"transactionTypeId", pg:"transaction_type_id"`
	TransactionType     *TransactionType  `json:"transactionType", pg:"rel:has-one, fk:transaction_type_id"`
	LastTransactionDate time.Time         `json:"lastTransactionDate", pg:"last_transaction_date"`
	Amount              int               `json:"amount", pg:"amount"`
}

type NewSubscriptionBody struct {
	WalletID           string      `json:"walletId" form:"walletId"`
	TransactionTypeID  string      `json:"transactionTypeId" form:"transactionTypeId"`
	SubscriptionTypeID string      `json:"subscriptionTypeId" pg:"subscription_type_id"`
	CustomRange        json.Number `json:"customRange", pg:"custom_range"`
	StartDate          time.Time   `json:"startDate" pg:"start_date"`
	Description        string      `json:"description" form:"description"`
	Amount             json.Number `json:"amount" form:"amount"`
}

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
	return trans
}
