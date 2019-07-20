package todo

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// EmptyTodo is empty todo model
var EmptyTodo = Todo{}

// Todo is todo data model.
type Todo struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Title     string
	Contents  string
	CreatedAt int64
}

// TodoResponse return instance of TodoResponse by Todo entity.
func (todo Todo) TodoResponse() TodoResponse {

	return TodoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Contents:  todo.Contents,
		CreatedAt: time.Unix(todo.CreatedAt, 0),
	}
}
