package controllers

import (
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	ApiService *services.ApiService
}

func NewApiController(as *services.ApiService, s *gin.RouterGroup) *ApiController {
	ac := new(ApiController)
	ac.ApiService = as

	s.GET("", ac.getFirst)
	s.POST("migrate", middleware.SecretCode, ac.postMigrate)

	return ac
}

func (ac *ApiController) getFirst(c *gin.Context) {
	apiModel := ac.ApiService.GetFirst()
	c.JSON(200, apiModel)
}

func (ac *ApiController) postMigrate(c *gin.Context) {
	mr, er := ac.ApiService.PostMigrate()

	if er.Message != "" {
		c.JSON(er.StatusCode, er)
	} else {
		c.JSON(200, mr)
	}
}
