package api

import (
	"wallet-api/pkg/controllers"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/configs"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func Routes(s *gin.Engine, db *pg.DB) {
	ver := s.Group(configs.Prefix)

	api := ver.Group("api")
	register := ver.Group("register")

	apiService := services.ApiService{Db: db}
	registerService := services.RegisterService{Db: db}

	controllers.NewApiController(&apiService, api)
	controllers.NewRegisterController(&registerService, register)
}
