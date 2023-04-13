package db

import (
	"context"
	"log"
	"os"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

type LoggerHook struct {
	QueryLogger *log.Logger
}

func (e *LoggerHook) BeforeQuery(ctx context.Context, qe *pg.QueryEvent) (context.Context, error) {
	return context.Background(), nil
}
func (e *LoggerHook) AfterQuery(ctx context.Context, qe *pg.QueryEvent) error {
	if qe.Err != nil {
		e.QueryLogger.Println(qe.Err)
	} else {
		fmtdQry, err := qe.FormattedQuery()
		if err != nil {
			return err
		}
		e.QueryLogger.Println(string(fmtdQry))
	}

	return nil
}

func CreateConnection(ctx context.Context, dbUrl string) *pg.DB {
	opt, err := pg.ParseURL(dbUrl)
	common.CheckError(err)
	file, err := os.OpenFile("querys.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	common.CheckError(err)
	QueryLogger := log.New(file, "Query: ", log.Ldate|log.Ltime|log.Lshortfile)
	conn := pg.Connect(opt)
	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	loggerHook := LoggerHook{QueryLogger}
	conn.AddQueryHook(&loggerHook)
	db := conn.WithContext(ctx)

	log.Printf("Successfully connected to %s!", dbUrl)
	return db
}
