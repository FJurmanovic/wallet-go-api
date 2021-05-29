package main

import (
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

	dbUrl := os.Getenv("DATABASE_URL")
	r := gin.New()
	r.Use(middleware.CORSMiddleware())

	conn := db.CreateConnection(dbUrl)
	api.Init(r, conn)

	server.Start(r)
}
