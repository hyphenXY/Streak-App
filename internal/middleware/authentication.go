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

        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := utils.ValidateJWT(token, "user")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Store user info for downstream handlers
        c.Set("user_id", claims["user_id"])
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

        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := utils.ValidateJWT(token, "admin")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Store user info for downstream handlers
        c.Set("user_id", claims["user_id"])
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

        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := utils.ValidateJWT(token, "root")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Store user info for downstream handlers
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}
