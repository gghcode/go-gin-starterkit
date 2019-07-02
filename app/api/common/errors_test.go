package common

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/pkg/errors"
)

func TestAddError(t *testing.T) {
	expectedErrLen := 2

	firstErr := errors.New("first error")
	secondErr := errors.New("second error")

	errResp := ErrorResponse{
		Errors: []APIError{
			APIError{
				Message: firstErr.Error(),
			},
		},
	}.AddError(secondErr)

	assert.Equal(t, expectedErrLen, len(errResp.Errors))
}

func TestNewErrResp(t *testing.T) {
	err := errors.New("fake error")
	errResp := NewErrResp(err)

	assert.Equal(t, errResp.Errors[0].Message, err.Error())
}
