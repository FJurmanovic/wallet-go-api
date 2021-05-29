package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	Db *pg.DB
}

func (as *WalletService) New(am *models.NewWalletBody) *models.Wallet {

	walletModel := new(models.Wallet)
	walletModel.Init()
	walletModel.UserID = am.UserID
	walletModel.Name = am.Name
	as.Db.Model(walletModel).Insert()
	return walletModel
}

func (as *WalletService) Get(am *models.Auth, embed string) *models.Wallet {
	wm := new(models.Wallet)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	common.GenerateEmbed(query, embed).Select()

	return wm
}

func (as *WalletService) GetAll(am *models.Auth, filtered *models.FilteredResponse) {
	wm := new([]models.Wallet)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	FilteredResponse(query, wm, filtered)
}
