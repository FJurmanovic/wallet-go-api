package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

/*
Init

Initializes Web API Routes.
	Args:
		*gin.Engine: Gin Engine.
		*pg.DB: Postgres Database Client.
*/
func Init(s *gin.Engine, db *pg.DB) {
	Routes(s, db)
}

type API struct {
	Api string `json:"api"`
}
