package repository

import (
	"context"
	"os"
	"time"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/common"
	"wallet-api/pkg/utl/configs"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-pg/pg/v10"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

/*
Get

Gets row from transaction table by id.

	   	Args:
			context.Context: Application context
			*model.Auth: Authentication object
			string: id to search
			*model.Params: url query parameters
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (us *UserRepository) Get(ctx context.Context, flt *filter.UserFilter, tx *pg.Tx) (*model.User, error) {
	wm := new(model.User)
	wm.Id = flt.Id

	commit := false
	if tx == nil {
		commit = true
		db := us.db.WithContext(ctx)
		tx, _ := db.Begin()
		defer tx.Rollback()
	}

	qry := tx.Model(wm)
	err := common.GenerateEmbed(qry, flt.Embed).WherePK().Select()
	if err != nil {
		return nil, err
	}

	if commit {
		tx.Commit()
	}

	return wm, nil
}

/*
Edit

Updates row in transaction table by id.

	   	Args:
			context.Context: Application context
			*model.TransactionEdit: Object to edit
			string: id to search
		Returns:
			*model.Transaction: Transaction object from database.
			*model.Exception: Exception payload.
*/
func (us *UserRepository) Edit(ctx context.Context, tm *model.User, tx *pg.Tx) (*model.User, error) {
	commit := false
	if tx == nil {
		commit = true
		db := us.db.WithContext(ctx)
		tx, _ := db.Begin()
		defer tx.Rollback()
	}

	_, err := tx.Model(tm).WherePK().UpdateNotZero()

	if err != nil {
		return nil, err
	}

	err = tx.Model(tm).WherePK().Select()
	if err != nil {
		return nil, err
	}

	if commit {
		tx.Commit()
	}

	return tm, nil
}

/*
Check

Inserts new row to users table.

	   	Args:
			context.Context: Application context
			*model.User: User object to create
		Returns:
			*model.User: User object from database
			*model.Exception
*/
func (us *UserRepository) Check(ctx context.Context, tx *pg.Tx, checkBody *model.User) (*model.User, error) {
	check := new(model.User)

	commit := false
	if tx == nil {
		commit = true
		db := us.db.WithContext(ctx)
		tx, _ := db.Begin()
		defer tx.Rollback()
	}

	err := tx.Model(check).Where("? = ?", pg.Ident("username"), checkBody.Username).WhereOr("? = ?", pg.Ident("email"), checkBody.Email).Select()
	if err != nil {
		return nil, err
	}
	if commit {
		tx.Commit()
	}
	return check, nil
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
func (us *UserRepository) Create(ctx context.Context, tx *pg.Tx, registerBody *model.User) (*model.User, error) {
	commit := false
	if tx == nil {
		commit = true
		db := us.db.WithContext(ctx)
		tx, _ := db.Begin()
		defer tx.Rollback()
	}

	_, err := tx.Model(registerBody).Insert()

	if err != nil {
		return nil, err
	}

	if commit {
		tx.Commit()
	}

	return registerBody, nil
}

func (us *UserRepository) CreateTx(ctx context.Context) (*pg.Tx, error) {
	db := us.db.WithContext(ctx)
	return db.Begin()
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
func (us *UserRepository) Login(ctx context.Context, loginBody *model.Login) (*model.Token, *model.Exception) {
	db := us.db.WithContext(ctx)

	check := new(model.User)
	exceptionReturn := new(model.Exception)
	tokenPayload := new(model.Token)

	tx, _ := db.Begin()
	defer tx.Rollback()

	tx.Model(check).Where("? = ?", pg.Ident("email"), loginBody.Email).Select()

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

	token, err := CreateToken(check, loginBody.RememberMe)
	common.CheckError(err)

	tokenPayload.Token = token

	tx.Commit()

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
func (us *UserRepository) Deactivate(ctx context.Context, auth *model.Auth) (*model.MessageResponse, *model.Exception) {
	db := us.db.WithContext(ctx)

	mm := new(model.MessageResponse)
	me := new(model.Exception)
	um := new(model.User)

	tx, _ := db.Begin()
	defer tx.Rollback()

	err := tx.Model(um).Where("? = ?", pg.Ident("id"), auth.Id).Select()

	if err != nil {
		me.ErrorCode = "404101"
		me.Message = "User not found"
		me.StatusCode = 404
		return mm, me
	}
	um.IsActive = false
	_, err = tx.Model(um).Where("? = ?", pg.Ident("id"), auth.Id).Update()

	if err != nil {
		me.ErrorCode = "400105"
		me.Message = "Could not deactivate user"
		me.StatusCode = 400
		return mm, me
	}

	mm.Message = "User successfully deactivated."

	tx.Commit()

	return mm, me
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
