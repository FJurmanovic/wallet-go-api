package middleware

import (
	"os"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/configs"

	"github.com/gin-gonic/gin"
)

func SecretCode(c *gin.Context) {
	exceptionReturn := new(models.Exception)
	secretCode := ExtractCode(c)
	secret := os.Getenv("SECRET_CODE")
	if secret == "" {
		secret = configs.SecretCode
	}
	if secret != secretCode {
		exceptionReturn.ErrorCode = "401101"
		exceptionReturn.StatusCode = 401
		exceptionReturn.Message = "Invalid secret code"
		c.AbortWithStatusJSON(exceptionReturn.StatusCode, exceptionReturn)
	}
	c.Next()
}

func ExtractCode(c *gin.Context) string {
	secret := c.GetHeader("SECRET_CODE")
	return secret
}
