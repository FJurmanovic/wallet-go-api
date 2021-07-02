package migrate

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10/orm"
)


func CreateTableUsers(db pg.DB) error {
	models := []interface{}{
		(*models.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
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
