package controllers

import (
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	ApiService *services.ApiService
}

func (ac *ApiController) Init(s *gin.Engine) {
	apiGroup := s.Group("/api")
	apiGroup.GET("", ac.getFirst)
}

func (ac *ApiController) getFirst(c *gin.Context) {
	apiModel := ac.ApiService.GetFirst()
	c.JSON(200, apiModel)
}
