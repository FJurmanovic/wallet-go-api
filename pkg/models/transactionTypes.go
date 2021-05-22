package models

type TransactionType struct {
	tableName struct{} `pg:"transactionTypes,alias:transactionTypes"`
	BaseModel
	Name string `json:"name" pg:"name"`
	Type string `json:"type" pg:"type"`
}

type NewTransactionTypeBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}