package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	SubscriptionService *services.SubscriptionService
}

func NewSubscriptionController(as *services.SubscriptionService, s *gin.RouterGroup) *SubscriptionController {
	wc := new(SubscriptionController)
	wc.SubscriptionService = as

	s.PUT("/end/:id", wc.End)
	s.POST("", wc.New)
	s.PUT("/:id", wc.Edit)
	s.GET("", wc.GetAll)
	s.GET("/:id", wc.Get)

	return wc
}

func (wc *SubscriptionController) New(c *gin.Context) {
	body := new(models.NewSubscriptionBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm := wc.SubscriptionService.New(c, body)
	c.JSON(200, wm)
}

func (wc *SubscriptionController) Edit(c *gin.Context) {
	body := new(models.SubscriptionEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	wm := wc.SubscriptionService.Edit(c, body, id)
	c.JSON(200, wm)
}

func (wc *SubscriptionController) Get(c *gin.Context) {
	body := new(models.Auth)
	params := new(models.Params)

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr := wc.SubscriptionService.Get(c, body, id, params)

	c.JSON(200, fr)
}

func (wc *SubscriptionController) End(c *gin.Context) {
	body := new(models.Auth)

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	id := c.Param("id")

	fr := wc.SubscriptionService.End(c, id)

	c.JSON(200, fr)
}

func (wc *SubscriptionController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	wc.SubscriptionService.GetAll(c, body, wallet, fr)

	c.JSON(200, fr)
}
