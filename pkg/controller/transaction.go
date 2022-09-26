package controller

import (
	"net/http"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service *service.TransactionService
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
func NewTransactionController(as *service.TransactionService, routeGroups *common.RouteGroups) *TransactionController {
	wc := &TransactionController{
		service: as,
	}

	routeGroups.Transaction.POST("", wc.New)
	routeGroups.Transaction.GET("", wc.GetAll)
	routeGroups.Transaction.PUT("/:id", wc.Edit)
	routeGroups.Transaction.GET("/:id", wc.Get)

	bulkGroup := routeGroups.Transaction.Group("bulk")
	{
		bulkGroup.PUT("", wc.BulkEdit)
	}

	checkGroup := routeGroups.Transaction.Group("check")
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
	body := new(model.NewTransactionBody)
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
// ROUTE (GET /transactions)
func (wc *TransactionController) GetAll(c *gin.Context) {
	body := new(model.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*model.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	noPendingQry, _ := c.GetQuery("noPending")
	noPending := noPendingQry != ""

	exception := wc.service.GetAll(c, body, wallet, fr, noPending)
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
	body := new(model.Auth)
	auth := c.MustGet("auth")
	body.Id = auth.(*model.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	exception := wc.service.Check(c, body, wallet, fr)
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
	body := new(model.TransactionEdit)
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
BulkEdit
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (PUT /transactions/:id)
func (wc *TransactionController) BulkEdit(c *gin.Context) {
	body := new([]model.TransactionEdit)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wm, exception := wc.service.BulkEdit(c, body)
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
	body := new(model.Auth)
	params := new(model.Params)

	auth := c.MustGet("auth")
	body.Id = auth.(*model.Auth).Id

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
