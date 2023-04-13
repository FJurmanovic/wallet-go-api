package controller

import (
	"fmt"
	"strconv"
	"strings"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"go.uber.org/dig"

	"github.com/gin-gonic/gin"
)

/*
InitializeControllers

Initializes Dependency Injection modules and registers controllers

	Args:
		*dig.Container: Dig Container
*/
func InitializeControllers(c *dig.Container) {
	controllerContainer := c.Scope("controller")
	service.InitializeServices(controllerContainer)

	controllerContainer.Invoke(NewApiController)
	controllerContainer.Invoke(NewUserController)
	controllerContainer.Invoke(NewWalletController)
	controllerContainer.Invoke(NewWalletHeaderController)
	controllerContainer.Invoke(NewTransactionController)
	controllerContainer.Invoke(NewTransactionStatusController)
	controllerContainer.Invoke(NewTransactionTypeController)
	controllerContainer.Invoke(NewSubscriptionController)
	controllerContainer.Invoke(NewSubscriptionTypeController)

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
