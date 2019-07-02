package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockRepository struct {
	mock.Mock
}

func (repo *mockRepository) CreateTodo(todo Todo) (Todo, error) {
	args := repo.Called(todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) GetTodos() ([]Todo, error) {
	args := repo.Called()
	return args.Get(0).([]Todo), args.Error(1)
}

func (repo *mockRepository) GetTodoByTodoID(todoID string) (Todo, error) {
	args := repo.Called(todoID)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) UpdateTodoByTodoID(todoID string, todo Todo) (Todo, error) {
	args := repo.Called(todoID, todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) RemoveTodoByTodoID(todoID string) (string, error) {
	args := repo.Called(todoID)
	return args.String(0), args.Error(1)
}

type mockTodoModelValidator struct {
	mock.Mock
}

func (validator *mockTodoModelValidator) Bind(c *gin.Context) error {
	args := validator.Called(c)
	return args.Error(0)
}

func (validator *mockTodoModelValidator) Todo() Todo {
	args := validator.Called()
	return args.Get(0).(Todo)
}

type controllerUnitTestSuite struct {
	suite.Suite

	ginEngine  *gin.Engine
	controller *Controller
	mockRepo   *mockRepository
}

func TestTodoControllerUnit(t *testing.T) {
	suite.Run(t, new(controllerUnitTestSuite))
}

func (suite *controllerUnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.ginEngine = gin.New()
	suite.mockRepo = &mockRepository{}
	suite.controller = NewController(suite.mockRepo)
	suite.controller.RegisterRoutes(suite.ginEngine)
}

func (suite *controllerUnitTestSuite) TestGetAllTodosExpectTodosFetched() {
	expectedCode := http.StatusOK
	expectedTodos := []Todo{
		Todo{
			ID:       uuid.NewV4(),
			Title:    "title",
			Contents: "contents",
		},
	}

	suite.mockRepo.
		On("GetTodos").
		Return(expectedTodos, nil)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodos, actual.Body)
}

func (suite *controllerUnitTestSuite) TestGetAllTodosExpectInternalErrReturn() {
	expectedCode := http.StatusInternalServerError

	suite.mockRepo.
		On("GetTodos").
		Return([]Todo{}, errors.New("MockError"))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerUnitTestSuite) TestGetTodoByIDExpectTodoFetched() {
	expectedCode := http.StatusOK
	expectedTodoID := uuid.NewV4()
	expectedTodo := Todo{
		ID:       expectedTodoID,
		Title:    "title",
		Contents: "contents",
	}

	suite.mockRepo.
		On("GetTodoByTodoID", expectedTodoID.String()).
		Return(expectedTodo, nil)

	req, err := http.NewRequest("GET", "/"+expectedTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodo, actual.Body)
}

func (suite *controllerUnitTestSuite) TestGetTodoByIDExpectNotFoundReturn() {
	expectedCode := http.StatusNotFound
	notExistsTodoID := uuid.NewV4()

	suite.mockRepo.
		On("GetTodoByTodoID", notExistsTodoID.String()).
		Return(Todo{}, common.ErrEntityNotFound)

	req, err := http.NewRequest("GET", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerUnitTestSuite) TestGetTodoByIDExpectInternalErrorReturn() {
	expectedCode := http.StatusInternalServerError
	todoID := uuid.NewV4()

	suite.mockRepo.
		On("GetTodoByTodoID", todoID.String()).
		Return(Todo{}, errors.New("Occurred error."))

	req, err := http.NewRequest("GET", "/"+todoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerUnitTestSuite) TestCreateTodoExpectTodoCreated() {
	expectedCode := http.StatusCreated
	expectedTodoID := uuid.NewV4()
	expectedTodo := Todo{
		Title:    "new title",
		Contents: "new contents",
	}

	suite.mockRepo.
		On("CreateTodo", expectedTodo).
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

func (suite *controllerUnitTestSuite) TestCreateTodoExpectBadRequestReturn() {
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

func (suite *controllerUnitTestSuite) TestUpdateTodoByIDExpectTodoUpdated() {
	expectedCode := http.StatusOK
	expectedTodoID := uuid.NewV4()
	expectedTodo := Todo{
		Title:    "updated title",
		Contents: "updated contents",
	}

	suite.mockRepo.
		On("UpdateTodoByTodoID", expectedTodoID.String(), expectedTodo).
		Return(expectedTodo, nil)

	expectedTodoJSON, err := json.Marshal(expectedTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+expectedTodoID.String(),
		bytes.NewBuffer(expectedTodoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
	assertEqualJSON(suite.T(), expectedTodo, actual.Body)
}

func (suite *controllerUnitTestSuite) TestUpdateTodoByIDExpectNotFoundReturn() {
	expectedCode := http.StatusNotFound
	notExistsTodo := Todo{
		ID:       uuid.Nil,
		Title:    "not exists title",
		Contents: "not exists contents",
	}

	suite.mockRepo.
		On("UpdateTodoByTodoID", notExistsTodo.ID.String(), notExistsTodo).
		Return(Todo{}, common.ErrEntityNotFound)

	notExistsTodoJSON, err := json.Marshal(notExistsTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+notExistsTodo.ID.String(),
		bytes.NewBuffer(notExistsTodoJSON))
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerUnitTestSuite) TestRemoveTodoByIDExpectTodoRemoved() {
	expectedCode := http.StatusOK
	expectedTodoID := uuid.NewV4()

	suite.mockRepo.
		On("RemoveTodoByTodoID", expectedTodoID.String()).
		Return(expectedTodoID.String(), nil)

	req, err := http.NewRequest("DELETE", "/"+expectedTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func (suite *controllerUnitTestSuite) TestRemoveTodoByIDExpectNotFoundReturn() {
	expectedCode := http.StatusNotFound
	notExistsTodoID := uuid.NewV4()

	suite.mockRepo.
		On("RemoveTodoByTodoID", notExistsTodoID.String()).
		Return("", common.ErrEntityNotFound)

	req, err := http.NewRequest("DELETE", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	actual := getActualResponse(suite, req)

	assert.Equal(suite.T(), expectedCode, actual.Code)
}

func getActualResponse(suite *controllerUnitTestSuite, req *http.Request) *httptest.ResponseRecorder {
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
