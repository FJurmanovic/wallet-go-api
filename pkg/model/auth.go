package model

type Token struct {
	Token string `json:"token"`
}

type Login struct {
	Email      string `json:"email" form:"email"`
	Password   string `json:"password" form:"password"`
	RememberMe bool   `json:"rememberMe" form:"rememberMe"`
}

type Auth struct {
	Id string
}

type CheckToken struct {
	Valid bool `json:"valid"`
}
