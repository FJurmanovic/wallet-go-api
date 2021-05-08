package models

import "time"

type CommonModel struct {
	Id          string    `json:"id" pg:"id"`
	DateCreated time.Time `json:"dateCreated" pg:"datecreated"`
	DateUpdated time.Time `json:"dateUpdated" pg:"dateupdated"`
}
