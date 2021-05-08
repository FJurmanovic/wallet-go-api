package models

type UserModel struct {
	tableName struct{} `pg:"users,alias:users"`
	CommonModel
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserReturnInfoModel struct {
	tableName struct{} `pg:"users,alias:users"`
	CommonModel
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (um *UserModel) Payload() UserReturnInfoModel {
	payload := UserReturnInfoModel{
		CommonModel: um.CommonModel,
		Username:    um.Username,
		Email:       um.Email,
	}
	return payload
}
