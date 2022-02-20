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

/*
NewWalletsController

Initializes WalletsController.
	Args:
		*services.WalletService: Wallet service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*WalletsController: Controller for "wallet" route interactions
*/
func NewWalletsController(as *services.WalletService, s *gin.RouterGroup) *WalletsController {
	wc := new(WalletsController)
	wc.WalletService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)
	s.PUT("/:id", wc.Edit)
	s.GET("/:id", wc.Get)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /wallet)
func (wc *WalletsController) New(c *gin.Context) {
	body := new(models.NewWalletBody)

	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	get := c.MustGet("auth")
	body.UserID = get.(*models.Auth).Id

	wm, exception := wc.WalletService.New(c, body)
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
func (wc *WalletsController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)

	exception := wc.WalletService.GetAll(c, body, fr)
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
func (wc *WalletsController) Edit(c *gin.Context) {
	body := new(models.WalletEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	wm, exception := wc.WalletService.Edit(c, body, id)
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
func (wc *WalletsController) Get(c *gin.Context) {
	params := new(models.Params)

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr, exception := wc.WalletService.Get(c, id, params)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
