package middleware

import (
	"log/slog"
	"net/http"

	"olimotracker/internal/api"
	"olimotracker/pkg/jwttoken"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDKey    = "userID"
	APIKeyHeader = "X-API-Key"
)

type Middleware struct {
	TokenValidator *jwttoken.TokenGenerator
	l              *slog.Logger
	RepositoryUser api.Repository
}

func NewMiddlware(tokenValidator *jwttoken.TokenGenerator, l *slog.Logger, repositoryUser api.Repository) *Middleware {
	return &Middleware{TokenValidator: tokenValidator, l: l, RepositoryUser: repositoryUser}
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

func (m *Middleware) APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if len(apiKey) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			m.l.Error("API key is required", "header", apiKey)
			c.Abort()
			return
		}

		parse, err := uuid.Parse(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			m.l.Error("error parsing API key: ", "err", err)
			c.Abort()
			return
		}

		user, err := m.RepositoryUser.GetUserByAPIKey(c, &parse)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			m.l.Error("error getting user by API key: ", "err", err)
			c.Abort()
			return
		}

		c.Set(UserIDKey, user.ID)
		c.Next()
	}
}
