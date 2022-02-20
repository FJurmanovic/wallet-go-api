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

/*
NewTransactionController

Initializes TransactionController.
	Args:
		*services.TransactionService: Transaction service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*TransactionController: Controller for "transaction" route interactions
*/
func NewTransactionController(as *services.TransactionService, s *gin.RouterGroup) *TransactionController {
	wc := new(TransactionController)
	wc.TransactionService = as

	s.POST("", wc.New)
	s.GET("", wc.GetAll)
	s.PUT("/:id", wc.Edit)
	s.GET("/:id", wc.Get)

	bulkGroup := s.Group("bulk")
	{
		bulkGroup.PUT("", wc.BulkEdit)
	}

	checkGroup := s.Group("check")
	{
		checkGroup.GET("", wc.Check)
	}

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /transactions)
func (wc *TransactionController) New(c *gin.Context) {
	body := new(models.NewTransactionBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm, exception := wc.TransactionService.New(c, body)
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
// ROUTE (GET /transactions)
func (wc *TransactionController) GetAll(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	noPendingQry, _ := c.GetQuery("noPending")
	noPending := noPendingQry != ""

	exception := wc.TransactionService.GetAll(c, body, wallet, fr, noPending)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

/*
Check
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /transactions)
func (wc *TransactionController) Check(c *gin.Context) {
	body := new(models.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	exception := wc.TransactionService.Check(c, body, wallet, fr)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}

/*
Edit
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (PUT /transactions/:id)
func (wc *TransactionController) Edit(c *gin.Context) {
	body := new(models.TransactionEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	wm, exception := wc.TransactionService.Edit(c, body, id)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}
	c.JSON(200, wm)
}

/*
BulkEdit
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (PUT /transactions/:id)
func (wc *TransactionController) BulkEdit(c *gin.Context) {
	body := new([]models.TransactionEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm, exception := wc.TransactionService.BulkEdit(c, body)
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
// ROUTE (GET /transactions/:id)
func (wc *TransactionController) Get(c *gin.Context) {
	body := new(models.Auth)
	params := new(models.Params)

	auth := c.MustGet("auth")
	body.Id = auth.(*models.Auth).Id

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	fr, exception := wc.TransactionService.Get(c, body, id, params)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
