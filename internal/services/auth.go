package services

import (
	"errors"
	"go/token"
	"os/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID int
	Role string
	jwt.StandardClaims
}

func GenerateToken(userID int, role string) (string, error) {
	expirationTime := time.Now() .Add(24 * time.Hour)
	claims := &Claims {
		UserID: userID,
		Role: role,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func (token *jwt.Token) (interface{}, error)  {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("token parse error")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}