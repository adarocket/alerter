package web

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	timeTokenAccess = 15 * time.Minute
	secretKey       = "secretKey"
)

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeTokenAccess).Unix(),
			IssuedAt:  time.Now().Unix(),
		})

	return token.SignedString([]byte(secretKey))
}

func IsValidToken(tokenStr string) (isnValid bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	claims, isTr := token.Claims.(*jwt.StandardClaims)
	if !isTr {
		log.Println("invalid type claims")
		return false
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return isnValid
	}

	return !isnValid
}
