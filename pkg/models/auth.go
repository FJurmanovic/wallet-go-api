package models

type Token struct {
	Token string `json:"token"`
}

type Login struct {
	Email    string
	Password string
}

type Auth struct {
	Id string
}
