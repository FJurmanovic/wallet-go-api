package controller

import (
	"fmt"
	"go.uber.org/dig"
	"strconv"
	"strings"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

/*
InitializeControllers

Initializes Dependency Injection modules and registers controllers

	Args:
		*dig.Container: Dig Container
*/
func InitializeControllers(c *dig.Container) {
	service.InitializeServices(c)

	c.Invoke(NewApiController)
	c.Invoke(NewUserController)
	c.Invoke(NewWalletController)
	c.Invoke(NewWalletHeaderController)
	c.Invoke(NewTransactionController)
	c.Invoke(NewTransactionStatusController)
	c.Invoke(NewTransactionTypeController)
	c.Invoke(NewSubscriptionController)
	c.Invoke(NewSubscriptionTypeController)

}

/*
FilteredResponse

Gets query parameters and populates FilteredResponse model.

	Args:
		*gin.Context: Gin Application Context
	Returns:
		*model.FilteredResponse: Filtered response
*/
func FilteredResponse(c *gin.Context) *model.FilteredResponse {
	filtered := new(model.FilteredResponse)
	page, _ := c.GetQuery("page")
	rpp, _ := c.GetQuery("rpp")
	sortBy, _ := c.GetQuery("sortBy")

	dividers := [5]string{"|", " ", ".", "/", ","}

	for _, div := range dividers {
		sortArr := strings.Split(sortBy, div)

		if len(sortArr) >= 2 {
			sortBy = fmt.Sprintf("%s %s", common.ToSnakeCase(sortArr[0]), strings.ToUpper(sortArr[1]))
		}
	}

	filtered.Embed, _ = c.GetQuery("embed")
	filtered.Page, _ = strconv.Atoi(page)
	filtered.Rpp, _ = strconv.Atoi(rpp)
	filtered.SortBy = sortBy

	return filtered
}
