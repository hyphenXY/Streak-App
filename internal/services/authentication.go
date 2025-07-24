package authentication

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey") // ideally from env

func GenerateJWT(userID uint, username string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(), // expires in 72h
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
