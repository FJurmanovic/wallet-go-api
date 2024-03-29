package model

type User struct {
	tableName struct{} `pg:"users,alias:users"`
	BaseModel
	IsActive bool   `json:"isActive" pg:"is_active" form:"isActive"`
	Username string `json:"username" pg:"username,notnull" form:"username"`
	Password string `json:"password" pg:"password,notnull" form:"password"`
	Email    string `json:"email" pg:"email,notnull" form:"email"`
}

type UserReturnInfo struct {
	tableName struct{} `pg:"users,alias:users"`
	BaseModel
	Username string `json:"username"`
	Email    string `json:"email"`
}

/*
Payload

Maps User object to UserReturnInfo object.

	Returns:
		*UserReturnInfo: mapped UserReturnInfo object
*/
func (um *User) Payload() *UserReturnInfo {
	payload := new(UserReturnInfo)
	payload.BaseModel = um.BaseModel
	payload.Username = um.Username
	payload.Email = um.Email

	return payload
}
