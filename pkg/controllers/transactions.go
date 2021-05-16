package controllers

import (
	"net/http"
	"wallet-api/pkg/models"
	"wallet-api/pkg/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(as *services.TransactionService, s *gin.RouterGroup) *TransactionController {
	wc := new(TransactionController)
	wc.TransactionService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)

	return wc
}

func (wc *TransactionController) New(c *gin.Context) {
	body := new(models.NewTransactionBody)
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm := wc.TransactionService.New(body)
	c.JSON(200, wm)
}

func (wc *TransactionController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")
	wallet, _ := c.GetQuery("walletId")

	wm := wc.TransactionService.GetAll(wallet, embed)

	c.JSON(200, wm)
}
