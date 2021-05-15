package middleware

import (
	"fmt"
	"os"
	"strings"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/configs"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	exceptionReturn := new(models.ExceptionModel)
	tokenString := ExtractToken(c)
	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		secret = configs.Secret
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		println(ok)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		exceptionReturn.ErrorCode = "401001"
		exceptionReturn.StatusCode = 401
		exceptionReturn.Message = "Invalid token"
		c.AbortWithStatusJSON(exceptionReturn.StatusCode, exceptionReturn)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, _ := claims["id"].(string)

		authModel := new(models.AuthModel)
		authModel.Id = userId

		c.Set("auth", authModel)
	}
	c.Next()
}

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
