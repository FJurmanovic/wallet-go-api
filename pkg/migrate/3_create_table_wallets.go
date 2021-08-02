package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"

	"github.com/go-pg/pg/v10/orm"
)

// Creates wallets table if it does not exist.
func CreateTableWallets(db pg.DB) error {
	models := []interface{}{
		(*models.Wallet)(nil),
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
