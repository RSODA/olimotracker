package galaxy

import (
	"net/http"
	"olimotracker/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service Service
	mw      *middleware.Middleware
}

func NewHandler(service Service, mw *middleware.Middleware) *handler {
	return &handler{service: service, mw: mw}
}

func (h *handler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/galaxy")
	group.Use(h.mw.AuthMiddleware())
	group.GET("/me", h.GetGalaxy)
}

func (h *handler) GetGalaxy(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found"})
		return
	}
	userID, ok := v.(*uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID is not a string"})
		return
	}
	galaxy, err := h.service.GetGalaxy(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, galaxy)
}
