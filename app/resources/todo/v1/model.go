package v1

import (
	"time"
)

// Todo is todo data model.
type Todo struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Contents        string    `json:"contents"`
	CreatedDateTime time.Time `json:"created_datetime"`
}
