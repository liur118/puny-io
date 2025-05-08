package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

func (s *JWTService) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}

func (s *JWTService) ValidateToken(tokenStr string) (string, error) {
	if !strings.HasPrefix(tokenStr, "Bearer ") {
		return "", jwt.ErrSignatureInvalid
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	username, ok := claims["user"].(string)
	if !ok {
		return "", jwt.ErrSignatureInvalid
	}

	return username, nil
}
