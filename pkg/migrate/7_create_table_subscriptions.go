package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Creates subscriptions table if it does not exist.
func CreateTableSubscriptions(db pg.DB) error {
	models := []interface{}{
		(*models.Subscription)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"subscriptions\": %s", err)
			return err
		} else {
			fmt.Println("Table \"subscriptions\" created successfully")
		}
	}
	return nil
}
