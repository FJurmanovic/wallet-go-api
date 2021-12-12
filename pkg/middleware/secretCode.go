package middleware

import (
	"net/http"
	"os"
	"wallet-api/pkg/models"
	"wallet-api/pkg/utl/configs"

	"github.com/gin-gonic/gin"
)

/*
SecretCode

Checks if secret code from body is valid.
   	Args:
   		*gin.Context: Gin Application Context.
*/
func SecretCode(c *gin.Context) {
	exceptionReturn := new(models.Exception)
	secretCode := ExtractCode(c)
	secret := os.Getenv("SECRET_CODE")
	if secret == "" {
		secret = configs.SecretCode
	}
	if secret != secretCode.SecretCode {
		exceptionReturn.ErrorCode = "401101"
		exceptionReturn.StatusCode = 401
		exceptionReturn.Message = "Invalid secret code"
		c.AbortWithStatusJSON(exceptionReturn.StatusCode, exceptionReturn)
	}
	c.Set("migrate", secretCode)
	c.Next()
}

/*
ExtractCode

Extracts the secret code from body.
   	Args:
   		*gin.Context: Gin Application Context.
*/
func ExtractCode(c *gin.Context) SecretCodeModel {
	secret := new(SecretCodeModel)
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return SecretCodeModel{}
	}
	return *secret
}

type SecretCodeModel struct {
	SecretCode string `json:"secretCode"`
	Version    string `json:"version"`
}
