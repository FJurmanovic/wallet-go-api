package model

type TransactionType struct {
	tableName struct{} `pg:"transactionTypes,alias:transactionTypes"`
	BaseModel
	Name string `json:"name" pg:"name"`
	Type string `json:"type" pg:"type,notnull"`
}

type NewTransactionTypeBody struct {
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}
