package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type TransactionStatusController struct {
	service *services.TransactionStatusService
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
func NewTransactionStatusController(as *services.TransactionStatusService, routeGroups *common.RouteGroups) *TransactionStatusController {
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
	body := new(models.NewTransactionStatusBody)
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
GetAll
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /transaction-status)
func (wc *TransactionStatusController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")

	wm, exception := wc.service.GetAll(c, embed)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}
