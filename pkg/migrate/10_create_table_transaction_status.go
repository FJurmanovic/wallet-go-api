package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

/*
CreateTableTransactionStatus

Creates transaction_status table if it does not exist.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with table creation
*/
func CreateTableTransactionStatus(db *pg.Tx) error {
	models := []interface{}{
		(*model.TransactionStatus)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"transaction_status\": %s", err)
			return err
		} else {
			fmt.Println("Table \"transaction_status\" created successfully")
		}
	}
	return nil
}
