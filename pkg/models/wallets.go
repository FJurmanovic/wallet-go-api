package models

type WalletModel struct {
	tableName struct{} `pg:"wallets,alias:wallets"`
	CommonModel
	WalletTypeID string               `json:"walletTypeId" pg:"wallet_type_id"`
	WalletType   *WalletTypeModel     `json:"walletType" pg:"rel:has-one,fk:wallet_type_id"`
	UserID       string               `json:"userId" pg:"user_id"`
	User         *UserReturnInfoModel `json:"user" pg:"rel:has-one,fk:user_id"`
}

type WalletTypeModel struct {
	tableName struct{} `pg:"walletTypes,alias:walletTypes"`
	CommonModel
	Name string `json:"name"`
}
