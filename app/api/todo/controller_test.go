package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gyuhwankim/go-gin-starterkit/app/api/common"
	"gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)


type controllerTestSuite struct {
	suite.Suite

	ginEngine  *gin.Engine
	controller *Controller
	mockRepo   *mockRepository
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(controllerTestSuite))
}

func (suite *controllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.ginEngine = gin.New()
	suite.mockRepo = &mockRepository{}
	suite.controller = NewController(suite.mockRepo)
	suite.controller.RegisterRoutes(suite.ginEngine)
}

func (suite *controllerTestSuite) TestShouldGetTodos() {
	expectedCode := http.StatusOK
	expectedTodos := []Todo{
		Todo{
			ID:       uuid.NewV4(),
			Title:    "title",
			Contents: "contents",
		},
	}

	suite.mockRepo.
		On("getTodos").
		Return(expectedTodos, nil)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodos, actual.Body)
}

func (suite *controllerTestSuite) TestShouldInternalErrorWhenGetTodos() {
	expectedCode := http.StatusInternalServerError

	suite.mockRepo.
		On("getTodos").
		Return([]Todo{}, errors.New("MockError"))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldGetTodoByTodoID() {
	expectedCode := http.StatusOK
	expectedTodoID := uuid.NewV4()
	expectedTodo := Todo{
		ID:       expectedTodoID,
		Title:    "title",
		Contents: "contents",
	}

	suite.mockRepo.
		On("getTodoByTodoID", expectedTodoID.String()).
		Return(expectedTodo, nil)

	req, err := http.NewRequest("GET", "/"+expectedTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodo, actual.Body)
}

func (suite *controllerTestSuite) TestShouldBeNotFoundWhenGetTodoByTodoID() {
	expectedCode := http.StatusNotFound
	notExistsTodoID := uuid.NewV4()

	suite.mockRepo.
		On("getTodoByTodoID", notExistsTodoID.String()).
		Return(Todo{}, common.ErrEntityNotFound)

	req, err := http.NewRequest("GET", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldBeInternalErrorWhenGetTodoByTodoID() {
	expectedCode := http.StatusInternalServerError
	todoID := uuid.NewV4()

	suite.mockRepo.
		On("getTodoByTodoID", todoID.String()).
		Return(Todo{}, errors.New("Occurred error."))

	req, err := http.NewRequest("GET", "/"+todoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldBeCreatedTodo() {
	expectedCode := http.StatusCreated
	expectedTodoID := uuid.NewV4()
	expectedTodo := Todo{
		Title:    "new title",
		Contents: "new contents",
	}

	suite.mockRepo.
		On("createTodo", expectedTodo).
		Return(Todo{
			ID:       expectedTodoID,
			Title:    expectedTodo.Title,
			Contents: expectedTodo.Contents,
		}, nil)

	expectedTodoJSON, err := json.Marshal(expectedTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(expectedTodoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)
	expectedTodo.ID = expectedTodoID

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodo, actual.Body)
}

func (suite *controllerTestSuite) TestShouldBeBadRequestWhenCreateTodo() {
	expectedCode := http.StatusBadRequest
	invalidTodo := Todo{
		Title:    "",
		Contents: "new contents",
	}

	mockTodoModelValidator := mockTodoModelValidator{}
	suite.controller.todoValidatorFactory = func() TodoModelValidator {
		return &mockTodoModelValidator
	}

	mockGinContext := gin.Context{}
	mockTodoModelValidator.
		On("Bind", &mockGinContext).
		Return(validator.ValidationErrors{})

	todoJSON, err := json.Marshal(invalidTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(todoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldBeUpdatedTodo() {
	expectedCode := http.StatusOK
	expectedTodo := Todo{
		ID:       uuid.NewV4(),
		Title:    "updated title",
		Contents: "updated contents",
	}

	suite.mockRepo.
		On("updateTodoByTodoID", expectedTodo.ID.String(), expectedTodo).
		Return(expectedTodo, nil)

	expectedTodoJSON, err := json.Marshal(expectedTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+expectedTodo.ID.String(),
		bytes.NewBuffer(expectedTodoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodo, actual.Body)
}

func (suite *controllerTestSuite) TestShouldBeNotFoundWhenUpdateTodoByTodoID() {
	expectedCode := http.StatusNotFound
	notExistsTodo := Todo{
		ID:       uuid.NewV4(),
		Title:    "not exists title",
		Contents: "not exists contents",
	}

	suite.mockRepo.
		On("updateTodoByTodoID", notExistsTodo.ID.String(), notExistsTodo).
		Return(Todo{}, common.ErrEntityNotFound)

	notExistsTodoJSON, err := json.Marshal(notExistsTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+notExistsTodo.ID.String(),
		bytes.NewBuffer(notExistsTodoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldBeRemovedTodo() {
	expectedCode := http.StatusNoContent
	expectedTodoID := uuid.NewV4()

	suite.mockRepo.
		On("removeTodoByTodoID", expectedTodoID.String()).
		Return(expectedTodoID.String(), nil)

	req, err := http.NewRequest("DELETE", "/"+expectedTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerTestSuite) TestShouldBeNotFoundWhenRemoveTodoByTodoID() {
	expectedCode := http.StatusNotFound
	notExistsTodoID := uuid.NewV4()

	suite.mockRepo.
		On("removeTodoByTodoID", notExistsTodoID.String()).
		Return("", common.ErrEntityNotFound)

	req, err := http.NewRequest("DELETE", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func getActualResponse(suite *controllerTestSuite, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	suite.ginEngine.ServeHTTP(recorder, req)

	return recorder
}

func assertEqualJSON(t *testing.T, data interface{}, buf *bytes.Buffer) {
	expected, err := json.Marshal(data)

	require.NoError(t, err)

	actual, err := ioutil.ReadAll(buf)

	require.NoError(t, err)
	assert.Equal(t, string(expected), string(actual))
}
