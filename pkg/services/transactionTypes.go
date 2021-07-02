package services

import (
	"context"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionTypeService struct {
	Db *pg.DB
}

func (as *TransactionTypeService) New(ctx context.Context, body *models.NewTransactionTypeBody) *models.TransactionType {
	db := as.Db.WithContext(ctx)

	tm := new(models.TransactionType)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	db.Model(tm).Insert()

	return tm
}

func (as *TransactionTypeService) GetAll(ctx context.Context, embed string) *[]models.TransactionType {
	db := as.Db.WithContext(ctx)

	wm := new([]models.TransactionType)

	query := db.Model(wm)
	common.GenerateEmbed(query, embed).Select()

	return wm
}
