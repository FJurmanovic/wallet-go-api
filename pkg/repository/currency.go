package repository

import (
	"context"
	"log"
	"time"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
)

type CurrencyRepository struct {
	db     *pg.DB
	logger *log.Logger
}

func NewCurrencyRepository(db *pg.DB, logger *log.Logger) *CurrencyRepository {
	return &CurrencyRepository{
		db:     db,
		logger: logger,
	}
}

/*
GetFirst

Gets first row from Currency table.

	   	Args:
	   		context.Context: Application context
		Returns:
			model.CurrencyModel: Currency object from database.
*/
func (as CurrencyRepository) Sync(ctx context.Context, rate *model.Rate, tx *pg.Tx) *model.Currency {
	currency := new(model.Currency)
	currency.Name = rate.Code

	any := tx.Model(currency).Where("name = ?", rate.Code).First()
	currency.Rate = float32(rate.Rate)
	if any != nil {
		currency.Init()
		_, err := tx.Model(currency).Insert()
		if err != nil {
			as.logger.Println(err)
			return nil
		}
	} else {
		currency.DateUpdated = time.Now()
		_, err := tx.Model(currency).WherePK().Update()
		if err != nil {
			as.logger.Println(err)
			return nil
		}
	}
	return currency
}

/*
GetFirst

Gets first row from Currency table.

	   	Args:
	   		context.Context: Application context
		Returns:
			model.CurrencyModel: Currency object from database.
*/
func (as CurrencyRepository) SyncBulk(ctx context.Context, rates *[]model.Rate) *[]model.Currency {
	tx, _ := as.db.BeginContext(ctx)
	defer tx.Rollback()

	currencies := new([]model.Currency)

	for _, r := range *rates {
		currency := as.Sync(ctx, &r, tx)
		if currency != nil {
			*currencies = append(*currencies, *currency)
		}
	}

	tx.Commit()

	return currencies
}
