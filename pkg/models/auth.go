package models

type TokenModel struct {
	Token string `json:"token"`
}

type LoginModel struct {
	Email    string
	Password string
}

type AuthModel struct {
	Id string
}
