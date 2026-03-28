package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessJWTSecret = []byte(os.Getenv("ACCESS_JWT_SECRET"))

type AccessJWTClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(ID, Role string) (string, error) {
	accessClaims := &AccessJWTClaims{
		ID:   ID,
		Role: Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return token.SignedString(accessJWTSecret)
}

func GenerateResetPasswordToken(ID string) (string, error) {
	resetClaims := &AccessJWTClaims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, resetClaims)
	return token.SignedString(accessJWTSecret)
}

func ParseAccessToken(tokenString string) (*AccessJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return accessJWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AccessJWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
