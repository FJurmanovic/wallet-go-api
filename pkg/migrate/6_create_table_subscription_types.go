package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

// Creates subscriptionTypes table if it does not exist.
func CreateTableSubscriptionTypes(db pg.DB) error {
	models := []interface{}{
		(*models.SubscriptionType)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"subscriptionTypes\": %s", err)
			return err
		} else {
			fmt.Println("Table \"subscriptionTypes\" created successfully")
		}
	}
	return nil
}
