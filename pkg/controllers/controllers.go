package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

// Gets query parameters and populates FilteredResponse model.
func FilteredResponse(c *gin.Context) *models.FilteredResponse {
	filtered := new(models.FilteredResponse)
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
