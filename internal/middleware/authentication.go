package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/utils"
)

func AuthUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		claims, err := utils.ValidateJWT(token, "user")
		if err != nil {
			switch err.Error() {
			case "token expired":
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			case "forbidden: role mismatch":
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: role mismatch"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("userId", claims["userId"]) // FIXED: match your JWT claim
		c.Next()
	}
}

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		claims, err := utils.ValidateJWT(token, "admin")
		if err != nil {
			switch err.Error() {
			case "token expired":
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			case "forbidden: role mismatch":
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: role mismatch"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("userId", claims["userId"]) // FIXED: match your JWT claim
		c.Next()
	}
}

func AuthRootMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))

		claims, err := utils.ValidateJWT(token, "root")
		if err != nil {
			switch err.Error() {
			case "token expired":
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			case "forbidden: role mismatch":
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: role mismatch"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("userId", claims["userId"]) // FIXED: match your JWT claim
		c.Next()
	}
}
