package services

import (
	"os"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-pg/pg/v10"
)

type UsersService struct {
	Db *pg.DB
}

func (us *UsersService) Create(registerBody *models.UserModel) (*models.UserModel, models.ExceptionModel) {
	var checkModel models.UserModel
	var exceptionReturn models.ExceptionModel

	us.Db.Model(&checkModel).Where("? = ?", pg.Ident("username"), registerBody.Username).WhereOr("? = ?", pg.Ident("email"), registerBody.Email).Select()
	if checkModel.Username != "" || checkModel.Email != "" {
		exceptionReturn.Message = "User already exists"
		exceptionReturn.ErrorCode = "400101"
		exceptionReturn.StatusCode = 400
		return &checkModel, exceptionReturn
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	common.CheckError(err)

	registerBody.Password = string(hashedPassword)
	_, err = us.Db.Model(registerBody).Insert()

	if err != nil {
		exceptionReturn.Message = "Error creating user"
		exceptionReturn.ErrorCode = "400102"
		exceptionReturn.StatusCode = 400
	}

	return registerBody, exceptionReturn
}

func (us *UsersService) Login(loginBody *models.LoginModel) (models.TokenModel, models.ExceptionModel) {
	var checkModel models.UserModel
	var exceptionReturn models.ExceptionModel
	var tokenPayload models.TokenModel

	us.Db.Model(&checkModel).Where("? = ?", pg.Ident("email"), loginBody.Email).Select()
	if checkModel.Email == "" {
		exceptionReturn.Message = "Email not found"
		exceptionReturn.ErrorCode = "400103"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	if bcrypt.CompareHashAndPassword([]byte(checkModel.Password), []byte(loginBody.Password)) != nil {
		exceptionReturn.Message = "Incorrect password"
		exceptionReturn.ErrorCode = "400104"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	token, err := CreateToken(checkModel)
	common.CheckError(err)

	tokenPayload.Token = token

	return tokenPayload, exceptionReturn
}

func CreateToken(user models.UserModel) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.Id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		secret = "Dond3sta"
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))

	return token, err
}
