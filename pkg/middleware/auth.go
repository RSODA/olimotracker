package middleware

import (
	"log/slog"
	"olimotracker/pkg/jwttoken"
	"strings"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

type Middleware struct {
	TokenValidator *jwttoken.TokenGenerator
	l              *slog.Logger
}

func NewMiddlware(tokenValidator *jwttoken.TokenGenerator, l *slog.Logger) *Middleware {
	return &Middleware{TokenValidator: tokenValidator, l: l}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			m.l.Error("Authorization header is required", "header", token)
			c.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			c.JSON(401, gin.H{"error": "invalid authorization format"})
			m.l.Error("invalid authorization format: ", "token", token)
			c.Abort()
			return
		}

		userID, err := m.TokenValidator.VerifyToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			m.l.Error("error verifying token: ", "err", err)
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
