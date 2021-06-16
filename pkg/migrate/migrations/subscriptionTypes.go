package migrations

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type SubscriptionTypesMigration struct {
	Db *pg.DB
}

func (am *SubscriptionTypesMigration) Create() error {
	models := []interface{}{
		(*models.SubscriptionType)(nil),
	}

	for _, model := range models {
		err := am.Db.Model(model).CreateTable(&orm.CreateTableOptions{
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

func (am *SubscriptionTypesMigration) Populate() error {
	daily := new(models.SubscriptionType)
	weekly := new(models.SubscriptionType)
	monthly := new(models.SubscriptionType)
	yearly := new(models.SubscriptionType)

	daily.Init()
	daily.Name = "Daily"
	daily.Type = "daily"

	weekly.Init()
	weekly.Name = "Weekly"
	weekly.Type = "weekly"

	monthly.Init()
	monthly.Name = "Monthly"
	monthly.Type = "monthly"

	yearly.Init()
	yearly.Name = "Yearly"
	yearly.Type = "yearly"

	_, err := am.Db.Model(daily).Where("? = ?", pg.Ident("type"), daily.Type).SelectOrInsert()

	_, err = am.Db.Model(weekly).Where("? = ?", pg.Ident("type"), weekly.Type).SelectOrInsert()

	_, err = am.Db.Model(monthly).Where("? = ?", pg.Ident("type"), monthly.Type).SelectOrInsert()

	_, err = am.Db.Model(yearly).Where("? = ?", pg.Ident("type"), yearly.Type).SelectOrInsert()

	return err
}
