package common

import (
	"strings"

	"github.com/go-pg/pg/v10/orm"
)

func GenerateEmbed(qr *orm.Query, embed string) *orm.Query{
	if embed != "" {
		rels := strings.Split(embed, ",")
		for _, rel := range rels {
			qr = qr.Relation(rel)
		}
	}
	return qr
}
