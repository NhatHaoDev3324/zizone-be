package utils

import (
	"math/rand"
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

func GenerateRandomToken(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
