package server

import (
	"fmt"
	"os"

	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine) *gin.Engine {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	err := r.Run(":" + port)
	if err != nil {
		msg := fmt.Sprintf("Running on %s:%s", common.GetIP(), port)
		println(msg)
	}
	return r
}
