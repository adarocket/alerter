package web

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	timeTokenAccess = 15 * time.Minute
	secretKey       = "secretKey" // FIXME: почему не в конфиге?
)

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeTokenAccess).Unix(),
			IssuedAt:  time.Now().Unix(),
		})

	return token.SignedString([]byte(secretKey))
}

func IsValidToken(tokenStr string) bool {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	claims, isTypeExist := token.Claims.(*jwt.StandardClaims)
	if !isTypeExist {
		log.Println("invalid type claims")
		return false
	}

	if err := claims.Valid(); err != nil {
		return false
	}

	return true
}
