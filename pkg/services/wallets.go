package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type WalletService struct {
	Db *pg.DB
}

func (as *WalletService) New(am *models.AuthModel) *models.WalletModel {
	walletType := as.GetType()

	walletModel := new(models.WalletModel)
	walletModel.Init()
	walletModel.UserID = am.Id
	walletModel.WalletTypeID = walletType.Id
	as.Db.Model(walletModel).Insert()
	return walletModel
}

func (as *WalletService) Get(am *models.AuthModel, embed string) *models.WalletModel {
	wm := new(models.WalletModel)

	query := as.Db.Model(wm).Where("? = ?", pg.Ident("user_id"), am.Id)
	common.GenerateEmbed(query, embed).Select()

	return wm
}

func (as *WalletService) GetType() *models.WalletTypeModel {
	wt := new(models.WalletTypeModel)

	as.Db.Model(wt).Select()

	return wt
}
