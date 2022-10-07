package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

/*
CreateTableSubscriptions

Creates subscriptions table if it does not exist.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with table creation
*/
func CreateTableSubscriptions(db pg.DB) error {
	models := []interface{}{
		(*model.Subscription)(nil),
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
