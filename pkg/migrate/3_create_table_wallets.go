package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/model"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

/*
CreateTableWallets

Creates wallets table if it does not exist.

	   	Args:
	   		*pg.DB: Postgres database client
		Returns:
			error: Returns if there is an error with table creation
*/
func CreateTableWallets(db *pg.Tx) error {
	models := []interface{}{
		(*model.Wallet)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error creating table \"wallets\": %s", err)
			return err
		} else {
			fmt.Println("Table \"wallets\" created successfully")
		}
	}
	return nil
}
