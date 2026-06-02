package jwttoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenGenerator struct {
	jwtSecret string
	ttl       time.Duration
}

func NewTokenGenerator(jwtSecret string, ttl time.Duration) *TokenGenerator {
	return &TokenGenerator{
		jwtSecret: jwtSecret,
		ttl:       ttl,
	}
}

func (s *TokenGenerator) GenerateToken(userID *uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.ttl).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *TokenGenerator) VerifyToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			return nil, errors.New("invalid user id")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user id")
		}

		return &userID, nil
	}
	return nil, errors.New("invalid token")
}
