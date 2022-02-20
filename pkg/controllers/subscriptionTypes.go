package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionTypeController struct {
	SubscriptionTypeService *services.SubscriptionTypeService
}

/*
NewSubscriptionTypeController

Initializes SubscriptionTypeController.
	Args:
		*services.SubscriptionTypeService: Subscription type service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*SubscriptionTypeController: Controller for "subscription-types" route interactions
*/
func NewSubscriptionTypeController(as *services.SubscriptionTypeService, s *gin.RouterGroup) *SubscriptionTypeController {
	wc := new(SubscriptionTypeController)
	wc.SubscriptionTypeService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /subscription-types)
func (wc *SubscriptionTypeController) New(c *gin.Context) {
	body := new(models.NewSubscriptionTypeBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm, exception := wc.SubscriptionTypeService.New(c, body)
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
// ROUTE (GET /subscription-types)
func (wc *SubscriptionTypeController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")

	wm, exception := wc.SubscriptionTypeService.GetAll(c, embed)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}
	c.JSON(200, wm)
}
