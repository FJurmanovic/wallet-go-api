package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	service *services.SubscriptionService
}

/*
NewSubscriptionController

Initializes SubscriptionController.

	Args:
		*services.SubscriptionService: Subscription service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*SubscriptionController: Controller for "subscription" route interactions
*/
func NewSubscriptionController(as *services.SubscriptionService, routeGroups *common.RouteGroups) *SubscriptionController {
	wc := &SubscriptionController{
		service: as,
	}

	routeGroups.Subscription.POST("", wc.New)
	routeGroups.Subscription.PUT("/:id", wc.Edit)
	routeGroups.Subscription.GET("/:id", wc.Get)
	routeGroups.Subscription.GET("", wc.GetAll)

	se := routeGroups.Subscription.Group("/end")
	{
		se.PUT("/:id", wc.End)
	}

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /subscription)
func (wc *SubscriptionController) New(c *gin.Context) {
	body := new(models.NewSubscriptionBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm, exception := wc.service.New(c, body)

	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}

/*
Edit
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (PUT /subscription/:id)
func (wc *SubscriptionController) Edit(c *gin.Context) {
	body := new(models.SubscriptionEdit)
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
// ROUTE (GET /subscription/:id)
func (wc *SubscriptionController) Get(c *gin.Context) {
	body := new(models.Auth)
	params := new(models.Params)

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr, exception := wc.service.Get(c, body, id, params)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

// ROUTE (PUT /subscription/end/:id)
func (wc *SubscriptionController) End(c *gin.Context) {
	body := new(models.Auth)

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	end := new(models.SubscriptionEnd)
	if err := c.ShouldBind(end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	fr, exception := wc.service.End(c, id)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

/*
GetAll
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /subscription)
func (wc *SubscriptionController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	exception := wc.service.GetAll(c, body, wallet, fr)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
