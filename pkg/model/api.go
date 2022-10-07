package model

type ApiModel struct {
	BaseModel
	tableName struct{} `pg:"api,alias:api"`
	Api       string   `json:"api"`
}
