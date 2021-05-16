package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id          string    `json:"id" pg:"id,pk"`
	DateCreated time.Time `json:"dateCreated" pg:"date_created"`
	DateUpdated time.Time `json:"dateUpdated" pg:"date_updated"`
}

func (cm *BaseModel) Init() {
	date := time.Now()
	cm.Id = uuid.NewString()
	cm.DateCreated = date
	cm.DateUpdated = date
}
