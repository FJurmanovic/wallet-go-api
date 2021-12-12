package controllers

import (
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	ApiService *services.ApiService
}

/*
NewApiController

Initializes ApiController.
	Args:
		*services.ApiService: API service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*ApiController: Controller for "api" interactions
*/
func NewApiController(as *services.ApiService, s *gin.RouterGroup) *ApiController {
	ac := new(ApiController)
	ac.ApiService = as

	s.GET("", ac.getFirst)
	s.POST("migrate", middleware.SecretCode, ac.postMigrate)

	return ac
}

/*
getFirst
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /api).
func (ac *ApiController) getFirst(c *gin.Context) {
	apiModel := ac.ApiService.GetFirst(c)
	c.JSON(200, apiModel)
}

/*
postMigrate

Requires "SECRET_CODE", "VERSION" (optional) from body.
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /api/migrate).
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
