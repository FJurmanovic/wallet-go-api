package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
)

/*
PopulateTransactionTypes

Populates transaction_types table if it exists.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with populating table
*/
func PopulateTransactionTypes(db *pg.Tx) error {
	gain := new(model.TransactionType)
	expense := new(model.TransactionType)

	gain.Init()
	gain.Name = "Gain"
	gain.Type = "gain"

	expense.Init()
	expense.Name = "Expense"
	expense.Type = "expense"

	_, err := db.Model(gain).Where("? = ?", pg.Ident("type"), gain.Type).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"transactionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"transactionTypes\" table.")
	}

	_, err = db.Model(expense).Where("? = ?", pg.Ident("type"), expense.Type).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"transactionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"transactionTypes\" table.")
	}

	return err
}
