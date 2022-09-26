package model

type Exception struct {
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
