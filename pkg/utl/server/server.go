package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Start() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	r.Run(":" + port)
	return r
}
