package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type TransactionStatusController struct {
	service *service.TransactionStatusService
}

/*
NewTransactionStatusController

Initializes TransactionStatusController.

	Args:
		*services.TransactionStatusService: Transaction Staus service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*TransactionStatusController: Controller for "transaction-status" route interactions
*/
func NewTransactionStatusController(as *service.TransactionStatusService, routeGroups *common.RouteGroups) *TransactionStatusController {
	wc := &TransactionStatusController{
		service: as,
	}

	routeGroups.TransactionStatus.POST("", wc.New)
	routeGroups.TransactionStatus.GET("", wc.GetAll)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /transaction-status)
func (wc *TransactionStatusController) New(c *gin.Context) {
	body := new(model.NewTransactionStatusBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mdl := body.ToTransactionStatus()

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
// ROUTE (GET /transaction-status)
func (wc *TransactionStatusController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")
	flt := filter.NewTransactionStatusFilter(model.Params{Embed: embed})

	wm, exception := wc.service.GetAll(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}
