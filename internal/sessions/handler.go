package sessions

import (
	"log/slog"
	"net/http"
	"olimotracker/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	s  Service
	l  *slog.Logger
	mw *middleware.Middleware
}

func NewHandler(s Service, l *slog.Logger, mw *middleware.Middleware) *Handler {
	return &Handler{s: s, l: l, mw: mw}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	sessions := r.Group("/sessions")
	sessions.Use(h.mw.AuthMiddleware())
	{
		sessions.POST("/", h.Create)
		sessions.GET("/", h.GetByUserID)
		sessions.GET("/:id", h.GetByID)
		sessions.GET("/by-category/:id", h.GetByCategoryID)
		sessions.PATCH("/:id", h.Update)
		sessions.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.s.Create(c.Request.Context(), userID.(*uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByID(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.s.GetByID(c.Request.Context(), &sessionID, userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.s.GetByUserID(c.Request.Context(), userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetByCategoryID(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.s.GetByCategoryID(c.Request.Context(), &categoryID, userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.s.Update(c.Request.Context(), &sessionID, userID.(*uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Delete(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	err = h.s.Delete(c.Request.Context(), &sessionID, userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted successfully"})
}
