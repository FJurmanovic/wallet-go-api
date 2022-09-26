package controllers

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type WalletsHeaderController struct {
	service *services.WalletService
}

/*
NewWalletsHeaderController

Initializes WalletsHeaderController.

	Args:
		*services.WalletService: Wallet service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*WalletsHeaderController: Controller for "wallet/wallet-header" route interactions
*/
func NewWalletsHeaderController(as *services.WalletService, routeGroups *common.RouteGroups) *WalletsHeaderController {
	wc := &WalletsHeaderController{
		service: as,
	}

	routeGroups.WalletHeader.GET("", wc.Get)

	return wc
}

/*
Get
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /wallet/wallet-header)
func (wc *WalletsHeaderController) Get(c *gin.Context) {
	body := new(models.Auth)

	walletId, _ := c.GetQuery("walletId")

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	wm, exception := wc.service.GetHeader(c, body, walletId)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}
