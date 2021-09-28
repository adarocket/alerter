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

func IsValidToken(tokenStr string) (isnValid bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	// FIXME что значит isTr ?
	claims, isTr := token.Claims.(*jwt.StandardClaims)
	if !isTr {
		log.Println("invalid type claims")
		return false
	}

	// FIXME эта проверка не обязательна, можно просто token.Valid
	if claims.ExpiresAt < time.Now().Unix() {
		return isnValid
	}

	return !isnValid
}
