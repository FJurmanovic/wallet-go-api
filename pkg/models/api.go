package models

type ApiModel struct {
	tableName struct{} `pg:"api,alias:api"`
	Api       string   `json:"api"`
}
