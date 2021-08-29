package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

/*
FilteredResponse

Adds filters to query and executes it.
   	Args:
   		*pg.Query: postgres query
		interface{}: model to be mapped from query execution.
		*models.FilteredResponse: filter options.
*/
func FilteredResponse(qry *pg.Query, mdl interface{}, filtered *models.FilteredResponse) {
	if filtered.Page == 0 {
		filtered.Page = 1
	}
	if filtered.Rpp == 0 {
		filtered.Rpp = 20
	}
	if filtered.SortBy == "" {
		filtered.SortBy = "date_created DESC"
	}
	qry = qry.Limit(filtered.Rpp).Offset((filtered.Page - 1) * filtered.Rpp).Order(filtered.SortBy)
	common.GenerateEmbed(qry, filtered.Embed)
	count, err := qry.SelectAndCount()
	common.CheckError(err)

	filtered.TotalRecords = count
	filtered.Items = mdl
}
