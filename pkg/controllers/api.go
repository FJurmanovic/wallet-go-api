package controllers

import (
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

	return ac
}

func (ac *ApiController) getFirst(c *gin.Context) {
	apiModel := ac.ApiService.GetFirst()
	c.JSON(200, apiModel)
}
