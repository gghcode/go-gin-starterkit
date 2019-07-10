package user

import "time"

// EmptyUser is empty user model
var EmptyUser = User{}

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    int64  `gorm:"not null;"`
}

// Response return new user response from user entity.
func (user User) Response() UserResponse {
	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		CreatedAt: time.Unix(user.CreatedAt, 0),
	}
}
