package migrate

import (
	"github.com/go-pg/pg/v10"
)

/*
Start

Starts database migration.
   	Args:
   		*pg.DB: Postgres database client
	Returns:
		error: Returns if there is an error with populating table
*/
func Start(conn *pg.DB, version string) []error {
	migration001 := Migration{
		Version: "001",
		Migrations: []interface{}{
			CreateTableApi,
			CreateTableUsers,
			CreateTableWallets,
			CreateTableTransactionTypes,
			CreateTableTransactions,
			CreateTableSubscriptionTypes,
			CreateTableSubscriptions,
		},
	}
	migration002 := Migration{
		Version: "002",
		Migrations: []interface{}{
			PopulateSubscriptionTypes,
			PopulateTransactionTypes,
		},
	}
	migration003 := Migration{
		Version: "003",
		Migrations: []interface{}{
			CreateTableTransactionStatus,
		},
	}
	migration004 := Migration{
		Version: "004",
		Migrations: []interface{}{
			PopulateTransactionStatus,
		},
	}

	migrationsMap := []Migration{
		migration001,
		migration002,
		migration003,
		migration004,
	}

	var errors []error

	for _, migrationCol := range migrationsMap {
		if version != "" && version == migrationCol.Version || version == "" {
			for _, migration := range migrationCol.Migrations {
				mgFunc, isFunc := migration.(func(pg.DB) error)
				if isFunc {
					err := mgFunc(*conn)
					if err != nil {
						errors = append(errors, err)
					}
				}
			}
		}
	}

	return errors
}

type Migration struct {
	Version    string
	Migrations []interface{}
}
