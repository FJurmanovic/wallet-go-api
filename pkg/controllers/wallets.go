package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type WalletsController struct {
	WalletService *services.WalletService
}

func NewWalletsController(as *services.WalletService, s *gin.RouterGroup) *WalletsController {
	wc := new(WalletsController)
	wc.WalletService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

func (wc *WalletsController) New(c *gin.Context) {
	body := new(models.NewWalletBody)

	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	get := c.MustGet("auth")
	body.UserID = get.(*models.Auth).Id

	wm := wc.WalletService.New(body)
	c.JSON(200, wm)
}

func (wc *WalletsController) Get(c *gin.Context) {
	body := new(models.Auth)

	embed, _ := c.GetQuery("embed")
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	wm := wc.WalletService.Get(body, embed)

	c.JSON(200, wm)
}

func (wc *WalletsController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)

	wc.WalletService.GetAll(body, fr)

	c.JSON(200, fr)

}
