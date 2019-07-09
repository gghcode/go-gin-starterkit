package user

import (
	"time"
)

var emptyUser = User{}

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    time.Time
}