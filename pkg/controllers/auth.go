package controllers

import (
	"net/http"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UsersService *services.UsersService
}

func NewAuthController(rs *services.UsersService, s *gin.RouterGroup) *AuthController {
	rc := new(AuthController)
	rc.UsersService = rs

	s.POST("login", rc.PostLogin)
	s.POST("register", rc.PostRegister)
	s.DELETE("deactivate", middleware.Auth, rc.Delete)

	return rc
}

func (rc *AuthController) PostLogin(c *gin.Context) {
	body := new(models.Login)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, exceptionReturn := rc.UsersService.Login(body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
	} else {
		c.JSON(200, returnedUser)
	}
}

func (rc *AuthController) PostRegister(c *gin.Context) {
	body := new(models.User)
	body.Init()
	body.IsActive = true
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, exceptionReturn := rc.UsersService.Create(body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
	} else {
		c.JSON(200, returnedUser.Payload())
	}
}
func (rc *AuthController) Delete(c *gin.Context) {
	auth := new(models.Auth)
	authGet := c.MustGet("auth")
	auth.Id = authGet.(*models.Auth).Id

	mr, er := rc.UsersService.Deactivate(auth)

	if er.Message != "" {
		c.JSON(er.StatusCode, er)
	} else {
		c.JSON(200, mr)
	}
}