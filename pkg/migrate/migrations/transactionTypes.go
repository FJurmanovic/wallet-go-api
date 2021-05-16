package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type TransactionTypesMigration struct {
	Db *pg.DB
}

func (am *TransactionTypesMigration) Create() {
	models := []interface{}{
		(*models.TransactionType)(nil),
	}

	for _, model := range models {
		err := am.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
		})
		if err != nil {
			log.Printf("Error Creating Table: %s", err)
		} else {
			fmt.Println("Table created successfully")
		}
	}
}
