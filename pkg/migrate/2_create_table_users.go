package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

// Creates api users if it does not exist.
func CreateTableUsers(db pg.DB) error {
	models := []interface{}{
		(*models.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Printf("Error creating table \"users\": %s", err)
			return err
		} else {
			fmt.Println("Table \"users\" created successfully")
		}
	}
	return nil
}
