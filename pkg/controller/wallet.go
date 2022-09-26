package controller

import (
	"net/http"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service *service.WalletService
}

/*
NewWalletController

Initializes WalletController.

	Args:
		*services.WalletService: Wallet service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*WalletController: Controller for "wallet" route interactions
*/
func NewWalletController(as *service.WalletService, routeGroups *common.RouteGroups) *WalletController {
	wc := &WalletController{
		service: as,
	}

	routeGroups.Wallet.POST("", wc.New)
	routeGroups.Wallet.GET("", wc.GetAll)
	routeGroups.Wallet.PUT("/:id", wc.Edit)
	routeGroups.Wallet.GET("/:id", wc.Get)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /wallet)
func (wc *WalletController) New(c *gin.Context) {
	body := new(model.NewWalletBody)

	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	get := c.MustGet("auth")
	body.UserID = get.(*model.Auth).Id

	wm, exception := wc.service.New(c, body)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}
	c.JSON(200, wm)
}

/*
GetAll
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /wallet)
func (wc *WalletController) GetAll(c *gin.Context) {
	body := new(model.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*model.Auth).Id

	fr := FilteredResponse(c)

	exception := wc.service.GetAll(c, body, fr)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

/*
Edit
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (PUT /wallet/:id)
func (wc *WalletController) Edit(c *gin.Context) {
	body := new(model.WalletEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	wm, exception := wc.service.Edit(c, body, id)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}
	c.JSON(200, wm)
}

/*
Get
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /wallet/:id)
func (wc *WalletController) Get(c *gin.Context) {
	params := new(model.Params)

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr, exception := wc.service.Get(c, id, params)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
