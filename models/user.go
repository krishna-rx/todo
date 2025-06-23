package models

import (
	"github.com/volatiletech/null"
	"time"
)

type UserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type User struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type TodoRequest struct {
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	ArchivedAt  time.Time `json:"archivedAt" db:"archived_at"`
}
type TODO struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	PendingAt   time.Time `json:"pendingAt" db:"pending_at"`
	ArchivedAt  time.Time `json:"archivedAt" db:"archived_at"`
}

type SessionData struct {
	CreatedAt  time.Time   `db:"created_at"`
	ArchivedAt *time.Time  `db:"archived_at"`
	UserID     null.String `db:"user_id"` // pointer for nullable
}

type UpdateTodo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type Session struct {
	ID        string    `db:"id"`
	JWTToken  string    `db:"jwt_token"`
	ExpiresAt time.Time `db:"expires_at"`
}
