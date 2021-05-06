package db

import (
	"fmt"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

func CreateConnection(dbUrl string) *pg.DB {
	opt, err := pg.ParseURL(dbUrl)
	common.CheckError(err)
	conn := pg.Connect(opt)

	fmt.Println("Successfully connected!")
	return conn
}
