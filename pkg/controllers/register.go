package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	UsersService *services.UsersService
}

func NewRegisterController(rs *services.UsersService, s *gin.RouterGroup) *RegisterController {
	rc := new(RegisterController)
	rc.UsersService = rs

	s.POST("", rc.Post)

	return rc
}

func (rc *RegisterController) Post(c *gin.Context) {
	body := new(models.User)
	body.Init()
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
