package controllers

import (
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
	s.GET("", wc.Get)

	return wc
}

func (wc *WalletsController) New(c *gin.Context) {
	body := new(models.AuthModel)

	get := c.MustGet("auth")
	body.Id = get.(*models.AuthModel).Id

	wm := wc.WalletService.New(body)
	c.JSON(200, wm)
}

func (wc *WalletsController) Get(c *gin.Context) {
	body := new(models.AuthModel)

	embed, _ := c.GetQuery("embed")
	auth := c.MustGet("auth")
	body.Id = auth.(*models.AuthModel).Id

	wm := wc.WalletService.Get(body, embed)

	c.JSON(200, wm)
}
