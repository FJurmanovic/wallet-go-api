package pg

import (
	"fmt"
	"os"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

func CreateConnection() func() *pg.DB {
	opt, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	common.CheckError(err)

	db := pg.Connect(opt)
	fmt.Println("Successfully connected!")
	return func() *pg.DB {
		return db
	}
}
