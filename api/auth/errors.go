package auth

import (
	"errors"
)

var (
	// ErrInvalidPassword is occurred when invalid user password
	ErrInvalidPassword = errors.New("Invalid user password")

	// ErrInvalidRefreshToken is occurred when invalid refresh token
	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
)
