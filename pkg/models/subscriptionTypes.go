package models

type SubscriptionType struct {
	tableName struct{} `pg:"subscriptionTypes,alias:subscriptionTypes"`
	BaseModel
	Name string `json:"name" pg:"name"`
	Type string `json:"type" pg:"type,notnull"`
}

type NewSubscriptionTypeBody struct {
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
}
