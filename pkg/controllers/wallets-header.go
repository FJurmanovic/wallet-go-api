package controllers

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type WalletsHeaderController struct {
	WalletService *services.WalletService
}

// Initializes WalletsHeaderController.
func NewWalletsHeaderController(as *services.WalletService, s *gin.RouterGroup) *WalletsHeaderController {
	wc := new(WalletsHeaderController)
	wc.WalletService = as

	s.GET("", wc.Get)

	return wc
}

// ROUTE (GET /wallet/wallet-header)
func (wc *WalletsHeaderController) Get(c *gin.Context) {
	body := new(models.Auth)

	walletId, _ := c.GetQuery("walletId")

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	wm := wc.WalletService.GetHeader(c, body, walletId)

	c.JSON(200, wm)
}
