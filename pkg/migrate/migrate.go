package migrate

import (
	"context"

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
			PopulateTransactionTypes,
			CreateTableSubscriptionTypes,
			PopulateSubscriptionTypes,
		},
	}
	migration002 := Migration{
		Version: "002",
		Migrations: []interface{}{
			CreateTableTransactionStatus,
		},
	}
	migration003 := Migration{
		Version: "003",
		Migrations: []interface{}{
			PopulateTransactionStatus,
		},
	}
	migration004 := Migration{
		Version: "004",
		Migrations: []interface{}{
			CreateTableSubscriptions,
			CreateTableTransactions,
		},
	}

	migrationsMap := []Migration{
		migration001,
		migration002,
		migration003,
		migration004,
	}

	var errors []error

	ctx := context.Background()

	tx, _ := conn.BeginContext(ctx)

	defer tx.Rollback()

	for _, migrationCol := range migrationsMap {
		if version != "" && version == migrationCol.Version || version == "" {
			for _, migration := range migrationCol.Migrations {
				mgFunc, isFunc := migration.(func(*pg.Tx) error)
				if isFunc {
					err := mgFunc(tx)
					if err != nil {
						errors = append(errors, err)
					}
				}
			}
		}
	}

	tx.CommitContext(ctx)

	return errors
}

type Migration struct {
	Version    string
	Migrations []interface{}
}
