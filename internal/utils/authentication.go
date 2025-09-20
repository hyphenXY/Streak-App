package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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

func GenerateJWT(claimsMap map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 15).Unix(), // expires in 15m
		"iat": time.Now().Unix(),
	}
	for k, v := range claimsMap {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 256-bit random
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10)) // random digit 0–9
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", n.Int64())
	}
	return otp, nil
}
