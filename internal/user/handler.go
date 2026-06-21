package user

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
	r.GET("/user/me", h.mw.AuthMiddleware(), h.GetUserByID)
	r.PUT("/user/me/api", h.mw.AuthMiddleware(), h.UpdateAPIByID)
}

func (h *handler) GetUserByID(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid not found"})
		return
	}

	userid := v.(*uuid.UUID)
	user, err := h.service.GetUserByID(c.Request.Context(), userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) UpdateAPIByID(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid not found"})
		return
	}

	userid := v.(*uuid.UUID)
	err := h.service.UpdateAPIKeyByID(c.Request.Context(), userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
