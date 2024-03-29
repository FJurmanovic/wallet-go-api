package model

import (
	"encoding/json"
	"math"
	"time"
)

type Transaction struct {
	tableName struct{} `pg:"transactions,alias:transactions"`
	BaseModel
	Description         string             `json:"description" pg:"description"`
	TransactionTypeID   string             `json:"transactionTypeId", pg:"transaction_type_id"`
	TransactionType     *TransactionType   `json:"transactionType", pg:"rel:has-one, fk:transaction_type_id"`
	TransactionStatusID string             `json:"transactionStatusId", pg:"transaction_status_id"`
	TransactionStatus   *TransactionStatus `json:"transactionStatus", pg:"rel:has-one, fk:transaction_status_id"`
	WalletID            string             `json:"walletId", pg:"wallet_id"`
	Amount              float32            `json:"amount", pg:"amount,default:0"`
	Wallet              *Wallet            `json:"wallet" pg:"rel:has-one, fk:wallet_id"`
	TransactionDate     time.Time          `json:"transactionDate" pg:"transaction_date, type:timestamptz"`
	SubscriptionID      string             `json:"subscriptionId", pg:"subscription_id"`
	Subscription        *Subscription      `json:"subscription", pg:"rel:has-one, fk:subscription_id"`
}

type NewTransactionBody struct {
	WalletID          string      `json:"walletId" form:"walletId"`
	TransactionTypeID string      `json:"transactionTypeId" form:"transactionTypeId"`
	TransactionDate   time.Time   `json:"transactionDate" form:"transactionDate"`
	Description       string      `json:"description" form:"description"`
	Amount            json.Number `json:"amount" form:"amount"`
}

func (body *NewTransactionBody) ToTransaction() *Transaction {
	tm := new(Transaction)
	amount, _ := body.Amount.Float64()

	tm.Init()
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.Description = body.Description
	tm.TransactionDate = body.TransactionDate
	tm.Amount = float32(math.Round(amount*100) / 100)

	if body.TransactionDate.IsZero() {
		tm.TransactionDate = time.Now()
	}
	return tm
}

type TransactionEdit struct {
	Id                  string      `json:"id" form:"id"`
	WalletID            string      `json:"walletId" form:"walletId"`
	TransactionTypeID   string      `json:"transactionTypeId" form:"transactionTypeId"`
	TransactionDate     time.Time   `json:"transactionDate" form:"transactionDate"`
	TransactionStatusID string      `json:"transactionStatusId" form:"transactionStatusId"`
	Description         string      `json:"description" form:"description"`
	Amount              json.Number `json:"amount" form:"amount"`
}

func (body *TransactionEdit) ToTransaction() *Transaction {
	tm := new(Transaction)
	amount, _ := body.Amount.Float64()

	tm.Id = body.Id
	tm.Description = body.Description
	tm.WalletID = body.WalletID
	tm.TransactionTypeID = body.TransactionTypeID
	tm.TransactionDate = body.TransactionDate
	tm.TransactionStatusID = body.TransactionStatusID
	tm.Amount = float32(math.Round(amount*100) / 100)

	return tm
}
