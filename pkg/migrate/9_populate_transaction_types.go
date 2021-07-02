package migrate

import (
	"github.com/go-pg/pg/v10"
	"wallet-api/pkg/models"
)

func PopulateTransactionTypes(db pg.DB) error {
	gain := new(models.TransactionType)
	expense := new(models.TransactionType)

	gain.Init()
	gain.Name = "Gain"
	gain.Type = "gain"

	expense.Init()
	expense.Name = "Expense"
	expense.Type = "expense"

	_, err := db.Model(gain).Where("? = ?", pg.Ident("type"), gain.Type).SelectOrInsert()

	_, err = db.Model(expense).Where("? = ?", pg.Ident("type"), expense.Type).SelectOrInsert()

	return err
}
