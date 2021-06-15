package models

type Wallet struct {
	tableName struct{} `pg:"wallets,alias:wallets"`
	BaseModel
	Name   string          `json:"name" pg:"name"`
	UserID string          `json:"userId" pg:"user_id"`
	User   *UserReturnInfo `json:"user" pg:"rel:has-one,fk:user_id"`
}

type NewWalletBody struct {
	Name   string `json:"name" form:"name"`
	UserID string `json:"userId" form:"userId"`
}

type WalletHeader struct {
	WalletId       string `json:"walletId"`
	CurrentBalance int    `json:"currentBalance"`
	LastMonth      int    `json:"lastMonth"`
	NextMonth      int    `json:"nextMonth"`
	Currency       string `json:"currency"`
}

type WalletTransactions struct {
	WalletId     string
	Transactions []Transaction
}
