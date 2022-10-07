package controller

import (
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type WalletHeaderController struct {
	service *service.WalletService
}

/*
NewWalletHeaderController

Initializes WalletHeaderController.

	Args:
		*services.WalletService: Wallet service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*WalletHeaderController: Controller for "wallet/wallet-header" route interactions
*/
func NewWalletHeaderController(as *service.WalletService, routeGroups *common.RouteGroups) *WalletHeaderController {
	wc := &WalletHeaderController{
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
func (wc *WalletHeaderController) Get(c *gin.Context) {
	walletId, _ := c.GetQuery("walletId")

	auth := c.MustGet("auth")
	userId := auth.(*model.Auth).Id
	embed, _ := c.GetQuery("embed")

	flt := filter.NewWalletHeaderFilter(model.Params{Embed: embed})
	flt.UserId = userId
	flt.WalletId = walletId

	wm, exception := wc.service.GetHeader(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}
