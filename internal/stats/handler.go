package stats

import (
	"net/http"
	"olimotracker/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	s  Service
	mw *middleware.Middleware
}

func NewHandler(s Service, mw *middleware.Middleware) *handler {
	return &handler{s: s, mw: mw}
}

func (h *handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/stats/me", h.mw.AuthMiddleware(), h.getMyStats)
}

func (h *handler) RegisterAPIRoutes(r *gin.RouterGroup) {
	r.GET("/stats", h.getMyStats)
}

func (h *handler) getMyStats(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := v.(*uuid.UUID)
	stats, err := h.s.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
