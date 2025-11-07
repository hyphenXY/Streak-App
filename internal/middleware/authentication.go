package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dataprovider "github.com/hyphenXY/Streak-App/internal/dataproviders"
	"github.com/hyphenXY/Streak-App/internal/models"
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

func IsUserClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		classID := c.Param("classID")
		classIDUint, err := strconv.ParseUint(classID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
			c.Abort()
			return
		}
		ifClassExists, err := dataprovider.IfClassExists(uint(classIDUint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check class existence"})
			c.Abort()
			return
		}
		if !ifClassExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			c.Abort()
			return
		}
		userIDVal, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}
		isEnrolled, err := dataprovider.IfAlreadyEnrolled(uint(userIDVal.(float64)), uint(classIDUint), &models.User_Classes{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment status"})
			c.Abort()
			return
		}
		if !isEnrolled {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not enrolled in this class"})
			c.Abort()
			return
		}
		c.Set("classID", uint(classIDUint))
		c.Next()
	}
}
