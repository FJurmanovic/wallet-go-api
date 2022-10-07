package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type SubscriptionTypeController struct {
	service *service.SubscriptionTypeService
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
func NewSubscriptionTypeController(as *service.SubscriptionTypeService, routeGroups *common.RouteGroups) *SubscriptionTypeController {
	wc := &SubscriptionTypeController{
		service: as,
	}

	routeGroups.SubscriptionType.POST("", wc.New)
	routeGroups.SubscriptionType.GET("", wc.GetAll)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /subscription-types)
func (wc *SubscriptionTypeController) New(c *gin.Context) {
	body := new(model.NewSubscriptionTypeBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mdl := body.ToSubscriptionType()

	wm, exception := wc.service.New(c, mdl)
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

	flt := filter.NewSubscriptionTypeFilter(model.Params{Embed: embed})

	wm, exception := wc.service.GetAll(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}
	c.JSON(200, wm)
}
