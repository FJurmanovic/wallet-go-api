package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

/*
CreateTableCurrencies

Creates Currencies table if it does not exist.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with table creation
*/
func CreateTableCurrencies(db *pg.Tx) error {
	models := []interface{}{
		(*model.Currency)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"currencies\": %s", err)
			return err
		} else {
			fmt.Println("Table \"currencies\" created successfully")
		}
	}
	return nil
}
