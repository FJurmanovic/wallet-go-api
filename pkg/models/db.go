package models

import (
	"time"

	"github.com/google/uuid"
)

type CommonModel struct {
	Id          string    `json:"id" pg:"id,pk"`
	DateCreated time.Time `json:"dateCreated" pg:"datecreated"`
	DateUpdated time.Time `json:"dateUpdated" pg:"dateupdated"`
}

func (cm *CommonModel) Init() {
	date := time.Now()
	cm.Id = uuid.NewString()
	cm.DateCreated = date
	cm.DateUpdated = date
}
