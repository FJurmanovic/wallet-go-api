package controllers

import (
	"net/http"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.UsersService
}

/*
NewAuthController

Initializes AuthController.

	Args:
		*services.UsersService: Users service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*AuthController: Controller for "auth" interactions
*/
func NewAuthController(rs *services.UsersService, routeGroups *common.RouteGroups) *AuthController {
	rc := &AuthController{
		service: rs,
	}

	routeGroups.Auth.POST("login", rc.PostLogin)
	routeGroups.Auth.POST("register", rc.PostRegister)
	routeGroups.Auth.DELETE("deactivate", middleware.Auth, rc.Delete)
	routeGroups.Auth.GET("check-token", rc.CheckToken)

	return rc
}

/*
PostLogin
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /auth/login).
func (rc *AuthController) PostLogin(c *gin.Context) {
	body := new(models.Login)
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, exceptionReturn := rc.service.Login(c, body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
		return
	} else {
		c.JSON(200, returnedUser)
	}
}

/*
PostRegister
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /auth/register).
func (rc *AuthController) PostRegister(c *gin.Context) {
	body := new(models.User)
	body.Init()
	body.IsActive = true
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, exceptionReturn := rc.service.Create(c, body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
		return
	} else {
		c.JSON(200, returnedUser.Payload())
	}
}

/*
Delete
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (DELETE /auth/deactivate).
func (rc *AuthController) Delete(c *gin.Context) {
	auth := new(models.Auth)
	authGet := c.MustGet("auth")
	auth.Id = authGet.(*models.Auth).Id

	mr, er := rc.service.Deactivate(c, auth)

	if er.Message != "" {
		c.JSON(er.StatusCode, er)
		return
	} else {
		c.JSON(200, mr)
	}
}

/*
CheckToken
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /auth/check-token).
func (rc *AuthController) CheckToken(c *gin.Context) {
	token, _ := c.GetQuery("token")
	re := new(models.CheckToken)

	_, err := middleware.CheckToken(token)

	if err != nil {
		re.Valid = false
		c.AbortWithStatusJSON(400, re)
		return
	}

	re.Valid = true

	c.JSON(200, re)
}
