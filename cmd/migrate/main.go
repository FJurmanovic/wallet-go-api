package main

import (
	"log"
	"os"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/utl/db"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	conn := db.CreateConnection(dbUrl)
	migrate.Start(conn)
}
