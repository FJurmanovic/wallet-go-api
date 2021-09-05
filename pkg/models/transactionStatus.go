package models

type TransactionStatus struct {
	tableName struct{} `pg:"transactionStatus,alias:transactionStatus"`
	BaseModel
	Name string `json:"name" pg:"name"`
	Status string `json:"status" pg:"status,notnull"`
}

type NewTransactionStatusBody struct {
	Name string `json:"name" form:"name"`
	Status string `json:"status" form:"status"`
}
