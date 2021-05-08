package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type UsersMigration struct {
	Db *pg.DB
}

func (am *UsersMigration) Create() {
	models := []interface{}{
		(*models.UserModel)(nil),
	}

	for _, model := range models {
		err := am.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			log.Printf("Error Creating Table: %s", err)
		} else {
			fmt.Println("Table created successfully")
		}
	}
}
