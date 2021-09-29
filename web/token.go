package web

import (
	"github.com/adarocket/alerter/config"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	timeTokenAccess = 15 * time.Minute
	//secretKey       = "secretKey"
)

func GenerateToken() (string, error) {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeTokenAccess).Unix(),
			IssuedAt:  time.Now().Unix(),
		})

	return token.SignedString([]byte(loadConfig.SecretKey))
}

func IsValidToken(tokenStr string) bool {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return false
	}

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(loadConfig.SecretKey), nil
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
