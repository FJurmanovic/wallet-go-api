package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	UsersService *services.UsersService
}

func NewLoginController(rs *services.UsersService, s *gin.RouterGroup) *LoginController {
	rc := new(LoginController)
	rc.UsersService = rs

	s.POST("", rc.Post)

	return rc
}

func (rc *LoginController) Post(c *gin.Context) {
	loginBody := new(models.LoginModel)
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, returnException := rc.UsersService.Login(loginBody)

	if returnException.Message != "" {
		c.JSON(returnException.StatusCode, returnException)
	} else {
		c.JSON(200, returnedUser)
	}

}
