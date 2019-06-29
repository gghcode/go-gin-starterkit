package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gopkg.in/go-playground/validator.v8"
)

type validatorsTestSuite struct {
	suite.Suite
}

func TestTodoValidatorUnit(t *testing.T) {
	suite.Run(t, new(validatorsTestSuite))
}

func (suite *validatorsTestSuite) TestBind() {
	testCases := []struct {
		description     string
		todoModel       Todo
		expectedErrType reflect.Type
	}{
		{
			"ShouldBeValid",
			Todo{
				Title:    "new title",
				Contents: "new contents",
			},
			nil,
		},
		{
			"ShouldBeInvalidWhenMinTitle",
			Todo{
				Title:    "",
				Contents: "new contents",
			},
			reflect.TypeOf(validator.ValidationErrors{}),
		},
		{
			"ShouldBeInvalidWhenMinContents",
			Todo{
				Title:    "new title",
				Contents: "",
			},
			reflect.TypeOf(validator.ValidationErrors{}),
		},
	}

	t := suite.T()

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			argTodo := testCase.todoModel
			argTodoJSON, err := json.Marshal(argTodo)
			require.NoError(t, err)

			mockRequest, err := http.NewRequest("POST",
				"/", bytes.NewBuffer(argTodoJSON))
			require.NoError(t, err)

			modelValidator := newTodoModelValidator()
			mockGinContext := gin.Context{
				Request: mockRequest,
			}

			actualErr := modelValidator.Bind(&mockGinContext)
			actualErrType := reflect.TypeOf(actualErr)

			assert.Equal(t, testCase.expectedErrType, actualErrType)
			if testCase.expectedErrType == nil {
				assert.Equal(t, argTodo.Title, modelValidator.todo.Title)
				assert.Equal(t, argTodo.Contents, modelValidator.todo.Contents)
			}
		})
	}
}
