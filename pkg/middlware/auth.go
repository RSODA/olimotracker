package middlware

import (
	"olimotracker/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

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

		userID, err := m.TokenValidator.VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
