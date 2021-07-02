package migrate

import (
	"github.com/go-pg/pg/v10"
)

func Start(conn *pg.DB, version string) {
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

	migrationsMap := []Migration{
		migration001,
		migration002,
	}

	for _, migrationCol := range migrationsMap {
		if version != "" && version == migrationCol.Version || version == "" {
			for _, migration := range migrationCol.Migrations {
				mgFunc, isFunc := migration.(func(pg.DB) error)
				if isFunc {
					mgFunc(*conn)
				}
			}
		}
	}
}

type Migration struct {
	Version string
	Migrations []interface{}
}

