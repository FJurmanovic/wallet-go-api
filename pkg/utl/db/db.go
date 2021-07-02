package db

import (
	"context"
	"fmt"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

func CreateConnection(dbUrl string, ctx context.Context) *pg.DB {
	opt, err := pg.ParseURL(dbUrl)
	common.CheckError(err)
	conn := pg.Connect(opt)
	db := conn.WithContext(ctx)

	fmt.Println("Successfully connected!")
	return db
}
