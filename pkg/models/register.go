package models

type User struct {
	tableName struct{} `pg:"users,alias:users"`
	BaseModel
	IsActive bool   `json:"isActive" pg:"is_active" form:"isActive"`
	Username string `json:"username" pg:"username" form:"username"`
	Password string `json:"password" pg:"password" form:"password"`
	Email    string `json:"email" pg:"email" form:"email"`
}

type UserReturnInfo struct {
	tableName struct{} `pg:"users,alias:users"`
	BaseModel
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (um *User) Payload() *UserReturnInfo {
	payload := new(UserReturnInfo)
	payload.BaseModel = um.BaseModel
	payload.Username = um.Username
	payload.Email = um.Email

	return payload
}
