package categories

import (
	"log/slog"
	"net/http"
	"olimotracker/pkg/middleware"
	"olimotracker/pkg/parse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler interface {
	Create(c *gin.Context)
	GetCategoriesByUserID(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RegisterRoutes(r *gin.Engine)
}

type handler struct {
	service Service
	l       *slog.Logger
	mw      *middleware.Middleware
}

func NewHandler(service Service, l *slog.Logger, mw *middleware.Middleware) Handler {
	return &handler{service: service, l: l, mw: mw}
}

func (h *handler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	category := &Category{
		Title: req.Title,
		Color: req.Color,
	}

	res, err := h.service.CreateCategory(c.Request.Context(), userID.(*uuid.UUID), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) GetCategoriesByUserID(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.service.GetCategoriesByUserID(c.Request.Context(), userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) GetCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")
	categoryUUID, err := parse.ParseUUID(categoryID, h.l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	res, err := h.service.GetCategoryByID(c.Request.Context(), &categoryUUID, userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) Update(c *gin.Context) {
	var req UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryID := c.Param("id")

	categoryUUID, err := parse.ParseUUID(categoryID, h.l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	category := &Category{
		ID:     categoryUUID,
		UserID: *userID.(*uuid.UUID),
		Title:  req.Title,
		Color:  req.Color,
	}

	res, err := h.service.UpdateCategory(c.Request.Context(), &categoryUUID, userID.(*uuid.UUID), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) Delete(c *gin.Context) {
	categoryID := c.Param("id")

	categoryUUID, err := parse.ParseUUID(categoryID, h.l)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}

	err = h.service.DeleteCategory(c.Request.Context(), &categoryUUID, userID.(*uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

func (h *handler) RegisterRoutes(r *gin.Engine) {
	catergories := r.Group("/categories")
	{
		catergories.POST("/", h.mw.AuthMiddleware(), h.Create)
		catergories.GET("/", h.mw.AuthMiddleware(), h.GetCategoriesByUserID)
		catergories.GET("/:id", h.mw.AuthMiddleware(), h.GetCategoryByID)
		catergories.PUT("/:id", h.mw.AuthMiddleware(), h.Update)
		catergories.DELETE("/:id", h.mw.AuthMiddleware(), h.Delete)
	}
}
