package models

import "github.com/gin-gonic/gin"

type FilteredResponse struct {
	Items        interface{} `json:"items"`
	SortBy       string      `json:"sortBy"`
	Embed        string      `json:"embed"`
	Page         int         `json:"page"`
	Rpp          int         `json:"rpp"`
	TotalRecords int         `json:"totalRecords"`
}

type ResponseFunc func(*gin.Context) *[]interface{}

type MessageResponse struct {
	Message string `json:"message"`
}
