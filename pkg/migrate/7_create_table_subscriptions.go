package migrate

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"wallet-api/pkg/models"
)

func CreateTableSubscriptions(db pg.DB) error {
	models := []interface{}{
		(*models.Subscription)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   false,
			FKConstraints: true,
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
