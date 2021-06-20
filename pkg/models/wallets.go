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
	WalletId       string  `json:"walletId"`
	CurrentBalance float32 `json:"currentBalance"`
	LastMonth      float32 `json:"lastMonth"`
	NextMonth      float32 `json:"nextMonth"`
	Currency       string  `json:"currency"`
}

type WalletTransactions struct {
	WalletId       string
	Transactions   []Transaction
	CurrentBalance float32
	LastMonth      float32
	NextMonth      float32
}
