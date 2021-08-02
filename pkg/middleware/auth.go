package middleware

import (
	"errors"
	"os"
	"strings"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/configs"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// Auth Middleware.
//
// Checks if token from header is valid and extracts the id.
func Auth(c *gin.Context) {
	exceptionReturn := new(models.Exception)
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

		authModel := new(models.Auth)
		authModel.Id = userId

		c.Set("auth", authModel)
	}
	c.Next()
}

// Extracts token from header
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

// Checks if token is valid
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
