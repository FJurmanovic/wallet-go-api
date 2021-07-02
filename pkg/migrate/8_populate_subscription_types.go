package migrate

import (
	"github.com/go-pg/pg/v10"
	"wallet-api/pkg/models"
)

func PopulateSubscriptionTypes(db pg.DB) error {
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

	_, err := db.Model(daily).Where("? = ?", pg.Ident("type"), daily.Type).SelectOrInsert()

	_, err = db.Model(weekly).Where("? = ?", pg.Ident("type"), weekly.Type).SelectOrInsert()

	_, err = db.Model(monthly).Where("? = ?", pg.Ident("type"), monthly.Type).SelectOrInsert()

	_, err = db.Model(yearly).Where("? = ?", pg.Ident("type"), yearly.Type).SelectOrInsert()

	return err
}
