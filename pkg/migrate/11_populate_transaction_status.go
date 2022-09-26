package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
)

/*
PopulateTransactionStatus

Populates transaction_status table if it exists.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with populating table
*/
func PopulateTransactionStatus(db pg.DB) error {
	completed := new(model.TransactionStatus)
	pending := new(model.TransactionStatus)
	deleted := new(model.TransactionStatus)

	completed.Init()
	completed.Name = "Completed"
	completed.Status = "completed"

	pending.Init()
	pending.Name = "Pending"
	pending.Status = "pending"

	deleted.Init()
	deleted.Name = "Deleted"
	deleted.Status = "deleted"

	_, err := db.Model(completed).Where("? = ?", pg.Ident("status"), completed.Status).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"transactionStatus\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"transactionStatus\" table.")
	}

	_, err = db.Model(pending).Where("? = ?", pg.Ident("status"), pending.Status).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"transactionStatus\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"transactionStatus\" table.")
	}

	_, err = db.Model(deleted).Where("? = ?", pg.Ident("status"), pending.Status).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"transactionStatus\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"transactionStatus\" table.")
	}

	return err
}
