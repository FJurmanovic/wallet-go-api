package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	RegisterService *services.RegisterService
}

func NewRegisterController(rs *services.RegisterService, s *gin.RouterGroup) *RegisterController {
	rc := new(RegisterController)
	rc.RegisterService = rs

	s.POST("", rc.Post)

	return rc
}

func (rc *RegisterController) Post(c *gin.Context) {
	registerBody := createModel()
	if err := c.ShouldBindJSON(&registerBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, returnException := rc.RegisterService.Create(&registerBody)

	if returnException.Message != "" {
		c.JSON(returnException.StatusCode, returnException)
	} else {
		c.JSON(200, returnedUser.Payload())
	}

}

func createModel() models.UserModel {
	commonModel := common.CreateDbModel()
	userModel := models.UserModel{CommonModel: commonModel}
	return userModel
}
