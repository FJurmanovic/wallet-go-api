package migrate

import (
	"fmt"
	"log"
	"wallet-api/pkg/models"

	"github.com/go-pg/pg/v10"
)

/*
PopulateSubscriptionTypes

Populates subscription_types table if it exists.
   	Args:
   		*pg.DB: Postgres database client
	Returns:
		error: Returns if there is an error with populating table
*/
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
	if err != nil {
		log.Printf("Error inserting row into \"subscriptionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"subscriptionTypes\" table.")
	}
	_, err = db.Model(weekly).Where("? = ?", pg.Ident("type"), weekly.Type).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"subscriptionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"subscriptionTypes\" table.")
	}

	_, err = db.Model(monthly).Where("? = ?", pg.Ident("type"), monthly.Type).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"subscriptionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"subscriptionTypes\" table.")
	}

	_, err = db.Model(yearly).Where("? = ?", pg.Ident("type"), yearly.Type).SelectOrInsert()
	if err != nil {
		log.Printf("Error inserting row into \"subscriptionTypes\" table: %s", err)
		return err
	} else {
		fmt.Println("Row inserted successfully into \"subscriptionTypes\" table.")
	}

	return err
}
