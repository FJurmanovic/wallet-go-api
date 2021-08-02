package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

// Creates transactionTypes table if it does not exist.
func CreateTableTransactionTypes(db pg.DB) error {
	models := []interface{}{
		(*models.TransactionType)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"api\": %s", err)
			return err
		} else {
			fmt.Println("Table \"api\" created successfully")
		}
	}
	return nil
}
