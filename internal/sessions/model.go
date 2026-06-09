package sessions

import (
	"time"

	"github.com/google/uuid"
)

type SessionResponse struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	CategoryTitle *string    `json:"category_title,omitempty"`
	CategoryColor *string    `json:"category_color,omitempty"`
	Duration      int        `json:"duration"`
	Note          *string    `json:"note,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CategoryMinutes struct {
	CategoryID    uuid.UUID
	CategoryTitle string
	CategoryColor string
	Minutes       int
}

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CategoryID *uuid.UUID
	Duration   int
	Note       *string
	CreatedAt  time.Time
}

type CreateSessionRequest struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Duration   int        `json:"duration" binding:"min=1"`
	Note       *string    `json:"note,omitempty"`
}

type UpdateSessionRequest struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Duration   *int       `json:"duration" binding:"min=1"`
	Note       *string    `json:"note,omitempty"`
}
