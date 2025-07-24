package utils

import (
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "time"
)

var jwtSecret = []byte("your-secret-key")

func ValidateJWT(tokenString string, requiredRole string) (map[string]interface{}, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // ✅ Check expiry
        if exp, ok := claims["exp"].(float64); ok {
            if int64(exp) < time.Now().Unix() {
                return nil, errors.New("token expired")
            }
        }

        // ✅ Check role
        if role, ok := claims["role"].(string); ok {
            if role != requiredRole {
                return nil, errors.New("forbidden: role mismatch")
            }
        } else {
            return nil, errors.New("role not found in token")
        }

        return claims, nil
    }

    return nil, errors.New("invalid token")
}
