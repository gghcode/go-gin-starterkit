package todo

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// TodoModel is todo data model.
type TodoModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Title     string
	Contents  string
	CreatedAt time.Time
}
