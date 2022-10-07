package service

import (
	"context"
	"os"
	"time"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
	"wallet-api/pkg/utl/common"
	"wallet-api/pkg/utl/configs"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

/*
Create

Inserts new row to users table.

	   	Args:
			context.Context: Application context
			*model.User: User object to create
		Returns:
			*model.User: User object from database
			*model.Exception
*/
func (us *UserService) Create(ctx context.Context, registerBody *model.User) (*model.User, *model.Exception) {
	exceptionReturn := new(model.Exception)

	tx, _ := us.repository.CreateTx(ctx)
	defer tx.Rollback()

	check, err := us.repository.Check(ctx, tx, registerBody)
	if err != nil {
		exceptionReturn.Message = "Error checking user"
		exceptionReturn.ErrorCode = "400139"
		exceptionReturn.StatusCode = 400
		return nil, exceptionReturn
	}

	if check.Username != "" || check.Email != "" {
		exceptionReturn.Message = "User already exists"
		exceptionReturn.ErrorCode = "400101"
		exceptionReturn.StatusCode = 400
		return nil, exceptionReturn
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	common.CheckError(err)

	registerBody.Password = string(hashedPassword)
	us.repository.Create(ctx, tx, registerBody)
	if err != nil {
		exceptionReturn.Message = "Error creating user"
		exceptionReturn.ErrorCode = "400102"
		exceptionReturn.StatusCode = 400
		return nil, exceptionReturn
	}

	tx.Commit()

	return registerBody, nil
}

/*
Login

Gets row from users table by email and valid password.

	   	Args:
			context.Context: Application context
			*model.Login: object to search
		Returns:
			*model.Token: new session token
			*model.Exception
*/
func (us *UserService) Login(ctx context.Context, body *model.Login) (*model.Token, *model.Exception) {
	exceptionReturn := new(model.Exception)
	tokenPayload := new(model.Token)
	loginBody := new(model.User)

	loginBody.Email = body.Email
	loginBody.Username = body.Email

	tx, _ := us.repository.CreateTx(ctx)
	defer tx.Rollback()

	check, err := us.repository.Check(ctx, tx, loginBody)
	if err != nil {
		exceptionReturn.Message = "Error checking user"
		exceptionReturn.ErrorCode = "400139"
		exceptionReturn.StatusCode = 400
		return nil, exceptionReturn
	}

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

	if bcrypt.CompareHashAndPassword([]byte(check.Password), []byte(body.Password)) != nil {
		exceptionReturn.Message = "Incorrect password"
		exceptionReturn.ErrorCode = "400104"
		exceptionReturn.StatusCode = 400
		return tokenPayload, exceptionReturn
	}

	token, err := CreateToken(check, body.RememberMe)
	common.CheckError(err)

	tokenPayload.Token = token

	return tokenPayload, exceptionReturn
}

/*
Deactivate

Updates row in users table.

IsActive column is set to false

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
		Returns:
			*model.MessageResponse
			*model.Exception
*/
func (us *UserService) Deactivate(ctx context.Context, flt *filter.UserFilter) (*model.MessageResponse, *model.Exception) {
	mm := new(model.MessageResponse)
	me := new(model.Exception)

	tx, _ := us.repository.CreateTx(ctx)
	defer tx.Rollback()

	um, err := us.repository.Get(ctx, flt, tx)
	if err != nil {
		me.ErrorCode = "404101"
		me.Message = "User not found"
		me.StatusCode = 404
		return nil, me
	}

	um.IsActive = false
	_, err = us.repository.Edit(ctx, um, tx)
	if err != nil {
		me.ErrorCode = "400105"
		me.Message = "Could not deactivate user"
		me.StatusCode = 400
		return nil, me
	}

	mm.Message = "User successfully deactivated."

	tx.Commit()

	return mm, nil
}

/*
CreateToken

Generates new jwt token.

It encodes the user id. Based on rememberMe it is valid through 48hours or 2hours.

	   	Args:
			*model.User: User object to encode
			bool: Should function generate longer lasting token (48hrs)
		Returns:
			string: Generated token
			error: Error that occured in the process
*/
func CreateToken(user *model.User, rememberMe bool) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.Id
	if rememberMe {
		atClaims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	} else {
		atClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	}

	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		secret = configs.Secret
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))

	return token, err
}
