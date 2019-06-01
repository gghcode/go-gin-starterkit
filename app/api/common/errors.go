package common

import (
	"github.com/pkg/errors"
)

var (
	// ErrEntityNotFound was occurred when can't found entity.
	ErrEntityNotFound = errors.New("Entity was not found.")

	// ErrInvalidUUID was occurred when uuid format was invalid.
	ErrInvalidUUID = errors.New("UUID was invalid.")
)
