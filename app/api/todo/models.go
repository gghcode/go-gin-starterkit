package todo

import (
	"time"

	"github.com/jinzhu/gorm"
)

// TodoModel is todo data model.
type TodoModel struct {
	gorm.Model
	ID        string `gorm:"type:uuid;primary_key;"`
	Title     string
	Contents  string
	CreatedAt time.Time
}
