package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type TransactionStatusController struct {
	TransactionStatusService *services.TransactionStatusService
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
func NewTransactionStatusController(as *services.TransactionStatusService, s *gin.RouterGroup) *TransactionStatusController {
	wc := new(TransactionStatusController)
	wc.TransactionStatusService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

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

	wm := wc.TransactionStatusService.New(c, body)
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

	wm := wc.TransactionStatusService.GetAll(c, embed)

	c.JSON(200, wm)
}
