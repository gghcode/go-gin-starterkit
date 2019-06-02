package todo

import (
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerTestSuite struct {
	suite.Suite

	ginEngine  *gin.Engine
	controller *Controller
	mockRepo   *mockRepository

	requestHTTP func(router http.Handler, method, path string) *httptest.ResponseRecorder
}

type mockRepository struct {
	mock.Mock
}

func (repo *mockRepository) createTodo(todo Todo) (Todo, error) {
	args := repo.Called(todo)

	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) getTodos() ([]Todo, error) {
	args := repo.Called()

	return args.Get(0).([]Todo), args.Error(1)
}

func (repo *mockRepository) getTodoByTodoID(todoID string) (Todo, error) {
	panic("Not Implement")
}

func (repo *mockRepository) updateTodoByTodoID(todoID string, todo Todo) (Todo, error) {
	panic("Not Implement")
}

func (repo *mockRepository) removeTodoByTodoID(todoID string) (string, error) {
	panic("Not Implement")
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(controllerTestSuite))
}

func (suite *controllerTestSuite) SetupTest() {
	suite.ginEngine = gin.New()
	suite.mockRepo = &mockRepository{}
	suite.controller = NewController(suite.mockRepo)
	suite.controller.RegisterRoutes(suite.ginEngine)
	suite.requestHTTP = func(router http.Handler, method, path string) *httptest.ResponseRecorder {
		req, err := http.NewRequest(method, path, nil)
		if err != nil {
			panic(err)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		return recorder
	}
}

func (suite *controllerTestSuite) TestShouldGetTodos() {
	returnTodos := []Todo{
		Todo{
			ID:       uuid.NewV4(),
			Title:    "title",
			Contents: "contents",
		},
	}

	suite.mockRepo.
		On("getTodos").
		Return(returnTodos, nil)

	req, err := http.NewRequest("GET", "/", nil)

	require.NoError(suite.T(), err)

	recorder := httptest.NewRecorder()
	suite.ginEngine.ServeHTTP(recorder, req)

	expected, err := json.Marshal(returnTodos)
	actual, err := ioutil.ReadAll(recorder.Body)

	require.Equal(suite.T(), http.StatusOK, recorder.Code)
	require.Equal(suite.T(), string(expected), string(actual))
}

// func TestGetAllTodos(t *testing.T) {
// 	engine := setupTestcase()
// 	recorder := requestHTTP(engine, "GET", "/")

// 	assert.Equal(t, http.StatusOK, recorder.Code)

// 	bodyBytes, err := ioutil.ReadAll(recorder.Body)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	bodyString := string(bodyBytes)
// 	expected := "{\"fasdf\":\"fasdf\"}"

// 	assert.Equal(t, expected, bodyString)
// }
