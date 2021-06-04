package models

type Token struct {
	Token string `json:"token"`
}

type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Auth struct {
	Id string
}

type CheckToken struct {
	Valid bool `json:"valid"`
}
