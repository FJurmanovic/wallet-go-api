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
	s.PUT("/:id", wc.Edit)
	s.GET("/:id", wc.Get)

	return wc
}

func (wc *WalletsController) New(c *gin.Context) {
	body := new(models.NewWalletBody)

	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	get := c.MustGet("auth")
	body.UserID = get.(*models.Auth).Id

	wm := wc.WalletService.New(c, body)
	c.JSON(200, wm)
}

func (wc *WalletsController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)

	wc.WalletService.GetAll(c, body, fr)

	c.JSON(200, fr)
}

func (wc *WalletsController) Edit(c *gin.Context) {
	body := new(models.WalletEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	wm := wc.WalletService.Edit(c, body, id)
	c.JSON(200, wm)
}

func (wc *WalletsController) Get(c *gin.Context) {
	params := new(models.Params)

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr := wc.WalletService.Get(c, id, params)

	c.JSON(200, fr)
}
