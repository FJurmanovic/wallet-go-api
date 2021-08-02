package main

import (
	"context"
	"os"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/utl/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	conn := db.CreateConnection(ctx, dbUrl)
	migrate.Start(conn, "")
}
