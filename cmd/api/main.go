package main

import (
	"context"
	"log"
	"os"
	"wallet-api/pkg/api"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/utl/db"
	"wallet-api/pkg/utl/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx := context.Background()
	dbUrl := os.Getenv("DATABASE_URL")
	r := gin.New()
	r.Use(middleware.CORSMiddleware())

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Print("Cannot open file logs.txt")
	}
	log.SetOutput(file)

	conn := db.CreateConnection(ctx, dbUrl)
	api.Init(r, conn)

	server.Start(r)
}
