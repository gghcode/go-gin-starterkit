package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// CreateTodoRequest is request model for creating todo.
type CreateTodoRequest struct {
	Title    string `json:"title" example:"<new title>" binding:"required"`
	Contents string `json:"contents" example:"<new contents>" binding:"required"`
}

// TodoResponse is todo response model.
type TodoResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `json:"create_at"`
}
