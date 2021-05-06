package api

import (
	"wallet-api/pkg/controllers"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func Routes(s *gin.Engine, db *pg.DB) {
	apiService := services.ApiService{Db: db}
	apiController := controllers.ApiController{ApiService: &apiService}
	apiController.Init(s)
}
