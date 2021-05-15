package models

type UserModel struct {
	tableName struct{} `pg:"users,alias:users"`
	CommonModel
	Username string `json:"username" pg:"username"`
	Password string `json:"password" pg:"password"`
	Email    string `json:"email" pg:"email"`
}

type UserReturnInfoModel struct {
	tableName struct{} `pg:"users,alias:users"`
	CommonModel
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (um *UserModel) Payload() *UserReturnInfoModel {
	payload := new(UserReturnInfoModel)
	payload.CommonModel = um.CommonModel
	payload.Username = um.Username
	payload.Email = um.Email

	return payload
}
