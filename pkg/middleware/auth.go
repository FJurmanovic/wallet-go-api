package middleware

import (
	"errors"
	"os"
	"strings"
	"wallet-api/pkg/model"
	"wallet-api/pkg/utl/configs"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

/*
Auth

Checks if token from header is valid and extracts the id.

	Args:
		*gin.Context: Gin Application Context.
*/
func Auth(c *gin.Context) {
	exceptionReturn := new(model.Exception)
	tokenString := ExtractToken(c)
	token, err := CheckToken(tokenString)
	if err != nil {
		exceptionReturn.ErrorCode = "401001"
		exceptionReturn.StatusCode = 401
		exceptionReturn.Message = "Invalid token"
		c.AbortWithStatusJSON(exceptionReturn.StatusCode, exceptionReturn)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, _ := claims["id"].(string)

		authModel := new(model.Auth)
		authModel.Id = userId

		c.Set("auth", authModel)
	}
	c.Next()
}

/*
ExtractToken

Extracts token from header

	Args:
		*gin.Context: Gin Application Context.
	Returns:
		string: Token extracted from context header
*/
func ExtractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	tokenArr := strings.Split(bearerToken, " ")
	if len(tokenArr) == 2 {
		bearerCheck := strings.ToLower(tokenArr[0])
		if bearerCheck == "bearer" {
			return tokenArr[1]
		}
	}
	return ""
}

/*
CheckToken

Checks if token is valid

	Args:
		string: Token to check
	Returns:
		*jwt.Token: Parsed token
		error: Returns if token is invalid or there was an error inside jwt.Parse function
*/
func CheckToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		secret = configs.Secret
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		var err error
		if !ok {
			err = errors.New("invalid token")
		}
		return []byte(secret), err
	})
	return token, err
}
