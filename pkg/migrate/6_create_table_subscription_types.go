package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

/*
CreateTableSubscriptionTypes

Creates subscription_types table if it does not exist.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with table creation
*/
func CreateTableSubscriptionTypes(db *pg.Tx) error {
	models := []interface{}{
		(*model.SubscriptionType)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
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
