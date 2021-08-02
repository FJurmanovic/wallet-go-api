package controllers

import (
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	ApiService *services.ApiService
}

// Initializes ApiController.
func NewApiController(as *services.ApiService, s *gin.RouterGroup) *ApiController {
	ac := new(ApiController)
	ac.ApiService = as

	s.GET("", ac.getFirst)
	s.POST("migrate", middleware.SecretCode, ac.postMigrate)

	return ac
}

// ROUTE (GET /api).
func (ac *ApiController) getFirst(c *gin.Context) {
	apiModel := ac.ApiService.GetFirst(c)
	c.JSON(200, apiModel)
}

// ROUTE (POST /api/migrate).
//
// Requires "SECRET_CODE", "VERSION" (optional) from body.
func (ac *ApiController) postMigrate(c *gin.Context) {
	migrateModel := c.MustGet("migrate")
	version := migrateModel.(middleware.SecretCodeModel).Version
	mr, er := ac.ApiService.PostMigrate(c, version)

	if er.Message != "" {
		c.JSON(er.StatusCode, er)
	} else {
		c.JSON(200, mr)
	}
}
