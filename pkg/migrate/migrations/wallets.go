package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type WalletsMigration struct {
	Db *pg.DB
}

func (am *WalletsMigration) Create() {
	models := []interface{}{
		(*models.WalletTypeModel)(nil),
		(*models.WalletModel)(nil),
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

func (am *WalletsMigration) PopulateTypes() {
	walletTypeModel := new(models.WalletTypeModel)
	walletTypeModel.Init()
	walletTypeModel.Name = "Test"
	am.Db.Model(walletTypeModel).Insert()
}
