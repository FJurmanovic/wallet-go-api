package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionTypeService struct {
	Db *pg.DB
}

func (as *TransactionTypeService) New(body *models.NewTransactionTypeBody) *models.TransactionType {
	tm := new(models.TransactionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	as.Db.Model(tm).Insert()

	return tm
}

func (as *TransactionTypeService) GetAll(embed string) *[]models.TransactionType {
	wm := new([]models.TransactionType)

	query := as.Db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
