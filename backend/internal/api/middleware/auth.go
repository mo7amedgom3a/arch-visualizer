package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired is a placeholder middleware that checks for the Authorization header.
// In a real application, this would validate a JWT token.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Token validation logic would go here.
		// For now, we trust any bearer token.
		token := parts[1]

		// Extract user info from token (mock)
		// c.Set("user_id", "extracted-uuid")
		_ = token

		c.Next()
	}
}
