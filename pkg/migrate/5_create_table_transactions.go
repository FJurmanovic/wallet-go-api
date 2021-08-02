package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

// Creates api transactions if it does not exist.
func CreateTableTransactions(db pg.DB) error {
	models := []interface{}{
		(*models.Transaction)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"transactions\": %s", err)
			return err
		} else {
			fmt.Println("Table \"transactions\" created successfully")
		}
	}
	return nil
}
