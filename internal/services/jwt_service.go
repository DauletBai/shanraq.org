package services

import (
	"go/token"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey string
}

type Claims struct {
	UserID int `json:"user_id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func (s *JWTService) GenerateToken(userID int, role string,) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	Claims := &Claims{
		UserID: userID,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(HS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *JWTService) VerifyToken(tokenString string) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error)  {
		return []byte(s.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}