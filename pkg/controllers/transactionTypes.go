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

func NewTransactionTypeController(as *services.TransactionTypeService, s *gin.RouterGroup) *TransactionTypeController {
	wc := new(TransactionTypeController)
	wc.TransactionTypeService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

func (wc *TransactionTypeController) New(c *gin.Context) {
	body := new(models.NewTransactionTypeBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm := wc.TransactionTypeService.New(body)
	c.JSON(200, wm)
}

func (wc *TransactionTypeController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")

	wm := wc.TransactionTypeService.GetAll(embed)

	c.JSON(200, wm)
}
