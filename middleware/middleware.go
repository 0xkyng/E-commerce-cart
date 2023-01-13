package middleware

import (
	"net/http"

	token "github.com/codekyng/E-commerce-cart.git/tokens"
	"github.com/gin-gonic/gin"
)

// Authentication
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token
		ClientToken := c.Request.Header.Get("token")
		// Check if token is empty
		if ClientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"No authorization header provided"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()

	}

}