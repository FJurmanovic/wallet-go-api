package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"github.com/go-pg/pg/v10"
)

func FilteredResponse(qry *pg.Query, mdl interface{}, filtered *models.FilteredResponse) {
	qry = qry.Limit(filtered.Rpp).Offset((filtered.Page - 1) * filtered.Rpp)
	common.GenerateEmbed(qry, filtered.Embed)
	count, err := qry.SelectAndCount()
	common.CheckError(err)

	filtered.TotalRecords = count
	filtered.Items = mdl
}
