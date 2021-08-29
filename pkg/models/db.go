package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id          string    `json:"id" pg:"id,pk,notnull"`
	DateCreated time.Time `json:"dateCreated" pg:"date_created"`
	DateUpdated time.Time `json:"dateUpdated" pg:"date_updated"`
}

/*
Init

Initializes base model with DateCreated, DateUpdated, and Id values.
*/
func (cm *BaseModel) Init() {
	date := time.Now()
	cm.Id = uuid.NewString()
	cm.DateCreated = date
	cm.DateUpdated = date
}
