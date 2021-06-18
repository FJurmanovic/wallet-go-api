package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type TransactionTypesMigration struct {
	Db *pg.DB
}

func (am *TransactionTypesMigration) Create() error {
	models := []interface{}{
		(*models.TransactionType)(nil),
	}

	for _, model := range models {
		err := am.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error Creating Table: %s", err)
			return err
		} else {
			fmt.Println("Table created successfully")
		}
	}
	return nil
}

func (am *TransactionTypesMigration) Populate() error {
	gain := new(models.TransactionType)
	expense := new(models.TransactionType)

	gain.Init()
	gain.Name = "Gain"
	gain.Type = "gain"

	expense.Init()
	expense.Name = "Expense"
	expense.Type = "expense"

	_, err := am.Db.Model(gain).Where("? = ?", pg.Ident("type"), gain.Type).SelectOrInsert()

	_, err = am.Db.Model(expense).Where("? = ?", pg.Ident("type"), expense.Type).SelectOrInsert()

	return err
}
