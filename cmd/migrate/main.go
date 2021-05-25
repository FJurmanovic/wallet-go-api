package main

import (
	"os"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/utl/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")

	conn := db.CreateConnection(dbUrl)
	migrate.Start(conn)
}
