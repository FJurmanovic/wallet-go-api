package services

import (
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-pg/pg/v10"
)

type RegisterService struct {
	Db *pg.DB
}

func (rs *RegisterService) Create(registerBody *models.UserModel) (*models.UserModel, models.ExceptionModel) {
	var checkModel models.UserModel
	var exceptionReturn models.ExceptionModel

	rs.Db.Model(&checkModel).Where("? = ?", pg.Ident("username"), registerBody.Username).WhereOr("? = ?", pg.Ident("email"), registerBody.Email).Select()
	if checkModel.Username != "" || checkModel.Email != "" {
		exceptionReturn.Message = "User already exists"
		exceptionReturn.ErrorCode = "400101"
		exceptionReturn.StatusCode = 400
		return &checkModel, exceptionReturn
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	common.CheckError(err)

	registerBody.Password = string(hashedPassword)
	_, err = rs.Db.Model(registerBody).Insert()

	if err != nil {
		exceptionReturn.Message = "Error creating user"
		exceptionReturn.ErrorCode = "400102"
		exceptionReturn.StatusCode = 400
	}

	return registerBody, exceptionReturn
}
