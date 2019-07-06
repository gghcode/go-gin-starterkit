package user

import (
	"time"
)

// EmptyUserID is empty of user id
const EmptyUserID = int64(0)

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    time.Time
}
