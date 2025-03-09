package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required!",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token!",
			})
			c.Abort()
			return
		}

		//extract token
		tokenStrings := tokenParts[1]

		//Validate the token
		claims, err := utils.ValidateToken(tokenStrings)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token!",
			})
			c.Abort()
			return
		}

		// Store user email in the context
		c.Set("userEmail", claims.Email)

		// Continue to the next middleware/handler
		c.Next()

	}
}
