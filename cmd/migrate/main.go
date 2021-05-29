package main

import (
	"fmt"
	"os"
	"wallet-api/pkg/migrate"
	"wallet-api/pkg/utl/db"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Start migrate")
	godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	fmt.Println("Database: ", dbUrl)

	conn := db.CreateConnection(dbUrl)
	migrate.Start(conn)
}
