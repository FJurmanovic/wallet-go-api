package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

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
	registerBody := createUserModel()
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, returnException := rc.UsersService.Create(&registerBody)

	if returnException.Message != "" {
		c.JSON(returnException.StatusCode, returnException)
	} else {
		c.JSON(200, returnedUser.Payload())
	}

}

func createUserModel() models.UserModel {
	commonModel := common.CreateDbModel()
	userModel := models.UserModel{CommonModel: commonModel}
	return userModel
}
