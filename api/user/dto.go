package user

import "time"

// CreateUserRequest is dto that contains info that require to create user.
type CreateUserRequest struct {
	UserName string `json:"username" example:"<new username>" binding:"required,min=4,max=100"`
	Password string `json:"password" example:"<new password>" binding:"required,min=8,max=50"`
}

// UpdateUserRequest is dto that contains info that require to update user.
type UpdateUserRequest struct {
	UserName string `json:"user_name" example:"<new user name>" binding:"min=4,max=100"`
}

// UserResponse is user response model.
type UserResponse struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"create_at"`
}
