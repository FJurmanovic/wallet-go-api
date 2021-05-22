package models

import "time"

type Transaction struct {
	tableName struct{} `pg:"transactions,alias:transactions"`
	BaseModel
	Description       string           `json:"description" pg:"description"`
	TransactionTypeID string           `json:"transactionTypeId", pg:"transaction_type_id"`
	TransactionType   *TransactionType `json:"transactionType", pg:"rel:has-one, fk:transaction_type_id"`
	WalletID          string           `json:"walletId", pg:"wallet_id"`
	Wallet            *Wallet          `json:"wallet" pg:"rel:has-one, fk:wallet_id"`
	TransactionDate   time.Time        `json:"transactionDate" pg:"transaction_date"`
}

type NewTransactionBody struct {
	WalletID          string    `json:"walletId"`
	TransactionTypeID string    `json:"transactionTypeId"`
	TransactionDate   time.Time `json:"transactionDate"`
	Description       string    `json:"description"`
}