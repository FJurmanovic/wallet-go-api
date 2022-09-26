package services

import (
	"context"
	"fmt"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type TransactionTypeService struct {
	db *pg.DB
}

func NewTransactionTypeService(db *pg.DB) *TransactionTypeService {
	return &TransactionTypeService{
		db: db,
	}
}

/*
New

Inserts new row to transaction type table.

	   	Args:
			context.Context: Application context
			*models.NewTransactionTypeBody: object to create
		Returns:
			*models.TransactionType: Transaction Type object from database.
			*models.Exception: Exception payload.
*/
func (as *TransactionTypeService) New(ctx context.Context, body *models.NewTransactionTypeBody) (*models.TransactionType, *models.Exception) {
	db := as.db.WithContext(ctx)

	tm := new(models.TransactionType)
	exceptionReturn := new(models.Exception)

	tm.Init()
	tm.Name = body.Name
	tm.Type = body.Type

	_, err := db.Model(tm).Insert()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400125"
		exceptionReturn.Message = fmt.Sprintf("Error inserting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return tm, nil
}

/*
GetAll

Gets all rows from transaction type table.

	   	Args:
			context.Context: Application context
			string: Relations to embed
		Returns:
			*[]models.TransactionType: List of Transaction type objects from database.
			*models.Exception: Exception payload.
*/
func (as *TransactionTypeService) GetAll(ctx context.Context, embed string) (*[]models.TransactionType, *models.Exception) {
	db := as.db.WithContext(ctx)

	wm := new([]models.TransactionType)
	exceptionReturn := new(models.Exception)

	query := db.Model(wm)
	err := common.GenerateEmbed(query, embed).Select()
	if err != nil {
		exceptionReturn.StatusCode = 400
		exceptionReturn.ErrorCode = "400133"
		exceptionReturn.Message = fmt.Sprintf("Error selecting row in \"transactionTypes\" table: %s", err)
		return nil, exceptionReturn
	}

	return wm, nil
}
