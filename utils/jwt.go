package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SecretKey = "welcome to github.nrg.com httpserver"
)

type CustomClaims struct {
	User string `json:"uid"`
	jwt.StandardClaims
}

func CreateToken(userId, role string) (string, error) {
	// Create the Claims
	claims := CustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * time.Duration(1)).Unix()),
			IssuedAt:  int64(time.Now().Unix()),
			Issuer:    "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	//logger.Debug("signed token ", tokenString)

	return tokenString, err
}

func CheckTokenValid(token string) bool {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		//logger.Errorf("parse token string failed: %v", err)
		return false
	}

	if _, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return true
	} else {
		//logger.Errorf("invalid token")
		return false
	}
}
