package middleware

import (
	"github.com/gin-gonic/gin"
	"jwt-with-go/helpers"
)

func authenticateHandlerFunc(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(403, gin.H{"error": "This user is not authorized"})
		c.Abort()
	}
	claims, err := helpers.ValidateToken(clientToken)
	if err != "" {
		c.JSON(403, gin.H{"err": err})
		c.Abort()
		return
	}
	c.Set("email", claims.Email)
	c.Set("firstName", claims.FirstName)
	c.Set("lastName", claims.LastName)
	c.Set("uid", claims.Uid)
	c.Set("userType", claims.UserType)
	c.Next()
}

func Authenticate() gin.HandlerFunc {
	return authenticateHandlerFunc
}
