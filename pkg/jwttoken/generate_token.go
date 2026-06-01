package jwttoken

import (
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
