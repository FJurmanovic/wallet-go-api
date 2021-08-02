package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type TransactionTypeController struct {
	TransactionTypeService *services.TransactionTypeService
}

// Initializes TransactionTypeController.
func NewTransactionTypeController(as *services.TransactionTypeService, s *gin.RouterGroup) *TransactionTypeController {
	wc := new(TransactionTypeController)
	wc.TransactionTypeService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

// ROUTE (POST /transaction-types)
func (wc *TransactionTypeController) New(c *gin.Context) {
	body := new(models.NewTransactionTypeBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm := wc.TransactionTypeService.New(c, body)
	c.JSON(200, wm)
}

// ROUTE (GET /transaction-types)
func (wc *TransactionTypeController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")

	wm := wc.TransactionTypeService.GetAll(c, embed)

	c.JSON(200, wm)
}
