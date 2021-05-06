package migrate

import (
	"wallet-api/pkg/migrate/migrations"

	"github.com/go-pg/pg/v10"
)

func Start(conn *pg.DB) {
	apiMigration := migrations.ApiMigration{Db: conn}

	apiMigration.Create()
}
