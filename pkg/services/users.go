package services

import (
	"os"
	"time"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/common"
	"wallet-api/pkg/utl/configs"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-pg/pg/v10"
)

type UsersService struct {
	Db *pg.DB
}

func (us *UsersService) Create(registerBody *models.User) (*models.User, *models.Exception) {
	check := new(models.User)
	exceptionReturn := new(models.Exception)

	us.Db.Model(check).Where("? = ?", pg.Ident("username"), registerBody.Username).WhereOr("? = ?", pg.Ident("email"), registerBody.Email).Select()
	if check.Username != "" || check.Email != "" {
		exceptionReturn.Message = "User already exists"
		exceptionReturn.ErrorCode = "400101"
		exceptionReturn.StatusCode = 400
		return check, exceptionReturn
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

func (us *UsersService) Login(loginBody *models.Login) (*models.Token, *models.Exception) {
	check := new(models.User)
	exceptionReturn := new(models.Exception)
	tokenPayload := new(models.Token)

	us.Db.Model(check).Where("? = ?", pg.Ident("email"), loginBody.Email).Select()
	if check.Email == "" {
		exceptionReturn.Message = "Email not found"
		exceptionReturn.ErrorCode = "400103"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	if !check.IsActive {
		exceptionReturn.Message = "Can't log in. User is deactivated."
		exceptionReturn.ErrorCode = "400106"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	if bcrypt.CompareHashAndPassword([]byte(check.Password), []byte(loginBody.Password)) != nil {
		exceptionReturn.Message = "Incorrect password"
		exceptionReturn.ErrorCode = "400104"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	token, err := CreateToken(check)
	common.CheckError(err)

	tokenPayload.Token = token

	return tokenPayload, exceptionReturn
}

func (us *UsersService) Deactivate(auth *models.Auth) (*models.MessageResponse, *models.Exception) {
	mm := new(models.MessageResponse)
	me := new(models.Exception)
	um := new(models.User)

	err := us.Db.Model(um).Where("? = ?", pg.Ident("id"), auth.Id).Select()

	if err != nil {
		me.ErrorCode = "404101"
		me.Message = "User not found"
		me.StatusCode = 404
		return mm, me
	}
	um.IsActive = false
	_, err = us.Db.Model(um).Where("? = ?", pg.Ident("id"), auth.Id).Update()

	if err != nil {
		me.ErrorCode = "400105"
		me.Message = "Could not deactivate user"
		me.StatusCode = 400
		return mm, me
	}

	mm.Message = "User successfully deactivated."

	return mm, me
}

func CreateToken(user *models.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.Id
	atClaims["exp"] = time.Now().Add(time.Minute).Unix()

	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		secret = configs.Secret
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))

	return token, err
}
