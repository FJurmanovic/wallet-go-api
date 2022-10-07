package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
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

	mdl := body.ToTransaction()

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
// ROUTE (GET /transactions)
func (wc *TransactionController) GetAll(c *gin.Context) {
	auth := c.MustGet("auth")
	userId := auth.(*model.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")
	embed, _ := c.GetQuery("embed")

	noPendingQry, _ := c.GetQuery("noPending")
	noPending := noPendingQry != ""

	flt := filter.NewTransactionFilter(model.Params{Embed: embed})
	flt.WalletId = wallet
	flt.NoPending = noPending
	flt.UserId = userId

	fr, exception := wc.service.GetAll(c, flt)
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
	auth := c.MustGet("auth")
	userId := auth.(*model.Auth).Id

	fr := FilteredResponse(c)
	wallet, _ := c.GetQuery("walletId")

	flt := filter.NewTransactionFilter(model.Params{})
	flt.WalletId = wallet
	flt.UserId = userId

	fr, exception := wc.service.Check(c, flt)
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
	mdl := body.ToTransaction()
	mdl.Id = id

	wm, exception := wc.service.Edit(c, mdl)
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

	mdl := new([]model.Transaction)
	for _, transaction := range *body {
		tm := transaction.ToTransaction()
		*mdl = append(*mdl, *tm)
	}

	wm, exception := wc.service.BulkEdit(c, mdl)
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
	params := new(model.Params)

	auth := c.MustGet("auth")
	userId := auth.(*model.Auth).Id

	id := c.Param("id")

	embed, _ := c.GetQuery("embed")
	params.Embed = embed

	flt := filter.NewTransactionFilter(*params)
	flt.Id = id
	flt.UserId = userId

	fr, exception := wc.service.Get(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, fr)
}
