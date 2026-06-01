package middleware

import (
	"olimotracker/pkg/jwttoken"
	"strings"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

type Middleware struct {
	TokenValidator *jwttoken.TokenGenerator
}

func NewMiddlware(tokenValidator *jwttoken.TokenGenerator) *Middleware {
	return &Middleware{TokenValidator: tokenValidator}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(token, bearerPrefix) {
			c.JSON(401, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		userID, err := m.TokenValidator.VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}
