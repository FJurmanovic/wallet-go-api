package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type ApiMigration struct {
	Db *pg.DB
}

func (am *ApiMigration) Create() error {
	models := []interface{}{
		(*models.ApiModel)(nil),
	}

	for _, model := range models {
		err := am.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Printf("Error Creating Table: %s", err)
			return err
		} else {
			fmt.Println("Table created successfully")
		}
	}
	return nil
}
