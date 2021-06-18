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

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

func (wc *SubscriptionController) New(c *gin.Context) {
	body := new(models.NewSubscriptionBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm := wc.SubscriptionService.New(body)
	c.JSON(200, wm)
}

func (wc *SubscriptionController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	wc.SubscriptionService.GetAll(body, wallet, fr)

	c.JSON(200, fr)
}
