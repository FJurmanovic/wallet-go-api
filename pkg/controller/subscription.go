package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	service *service.SubscriptionService
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
func NewSubscriptionController(as *service.SubscriptionService, routeGroups *common.RouteGroups) *SubscriptionController {
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
	body := new(model.NewSubscriptionBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mdl := body.ToSubscription()

	wm, exception := wc.service.New(c, mdl)

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
	body := new(model.SubscriptionEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	mdl := body.ToSubscription()
	mdl.Id = id

	wm, exception := wc.service.Edit(c, mdl)
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
	params := new(model.Params)

	auth := c.MustGet("auth")
	userId := auth.(*model.Auth).Id

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	flt := filter.NewSubscriptionFilter(*params)
	flt.Id = id
	flt.UserId = userId

	fr, exception := wc.service.Get(c, *flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

// ROUTE (PUT /subscription/end/:id)
func (wc *SubscriptionController) End(c *gin.Context) {
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
	auth := c.MustGet("auth")

	wallet, _ := c.GetQuery("walletId")
	embed, _ := c.GetQuery("embed")

	flt := filter.NewSubscriptionFilter(model.Params{Embed: embed})
	flt.WalletId = wallet
	flt.UserId = auth.(*model.Auth).Id

	fr, exception := wc.service.GetAll(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
