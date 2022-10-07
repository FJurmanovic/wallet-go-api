package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type TransactionTypeController struct {
	service *service.TransactionTypeService
}

/*
NewTransactionTypeController

Initializes TransactionTypeController.

	Args:
		*services.TransactionTypeService: Transaction Type service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*TransactionTypeController: Controller for "transaction-types" route interactions
*/
func NewTransactionTypeController(as *service.TransactionTypeService, routeGroups *common.RouteGroups) *TransactionTypeController {
	wc := &TransactionTypeController{
		service: as,
	}

	routeGroups.TransactionType.POST("", wc.New)
	routeGroups.TransactionType.GET("", wc.GetAll)

	return wc
}

/*
New
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /transaction-types)
func (wc *TransactionTypeController) New(c *gin.Context) {
	body := new(model.NewTransactionTypeBody)
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mdl := body.ToTransactionType()

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
// ROUTE (GET /transaction-types)
func (wc *TransactionTypeController) GetAll(c *gin.Context) {
	embed, _ := c.GetQuery("embed")

	flt := filter.NewTransactionTypeFilter(model.Params{Embed: embed})

	wm, exception := wc.service.GetAll(c, flt)
	if exception != nil {
		c.JSON(exception.StatusCode, exception)
		return
	}

	c.JSON(200, wm)
}
