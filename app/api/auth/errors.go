package auth

import (
	"errors"
)

var (
	// ErrInvalidPassword is occurred when invalid user password
	ErrInvalidPassword = errors.New("Invalid user password")
)
