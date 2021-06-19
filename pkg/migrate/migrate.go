package migrate

import (
	"wallet-api/pkg/migrate/migrations"

	"github.com/go-pg/pg/v10"
)

func Start(conn *pg.DB) error {
	apiMigration := migrations.ApiMigration{Db: conn}
	usersMigration := migrations.UsersMigration{Db: conn}
	walletsMigration := migrations.WalletsMigration{Db: conn}
	transactionTypesMigration := migrations.TransactionTypesMigration{Db: conn}
	transactionsMigration := migrations.TransactionsMigration{Db: conn}
	subscriptionTypesMigration := migrations.SubscriptionTypesMigration{Db: conn}
	subscriptionsMigration := migrations.SubscriptionsMigration{Db: conn}

	err := apiMigration.Create()
	err = usersMigration.Create()
	err = walletsMigration.Create()
	err = transactionTypesMigration.Create()
	err = subscriptionTypesMigration.Create()
	err = subscriptionsMigration.Create()
	err = transactionsMigration.Create()

	err = subscriptionTypesMigration.Populate()
	err = transactionTypesMigration.Populate()

	return err
}
