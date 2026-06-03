package categories

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Color     string
	CreatedAt time.Time
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
	Title string `json:"title" binding:"required,min=1,max=35"`
	Color string `json:"color" binding:"required,len=7"`
}

type UpdateCategoryRequest struct {
	Title string `json:"title" binding:"omitempty,min=1,max=35"`
	Color string `json:"color" binding:"omitempty,len=7"`
}
