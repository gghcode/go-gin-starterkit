package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerIntegrationTestSuite struct {
	suite.Suite

	ginEngine  *gin.Engine
	dbConn     *db.Conn
	controller *Controller

	testTodos []Todo
}

func TestTodoControllerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(controllerIntegrationTestSuite))
}

func (suite *controllerIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	dbConn, err := db.NewConn(conf)
	require.NoError(suite.T(), err)

	todoRepo := NewRepository(dbConn)

	suite.ginEngine = gin.New()
	suite.dbConn = dbConn
	suite.controller = NewController(todoRepo)
	suite.controller.RegisterRoutes(suite.ginEngine)

	suite.testTodos, err = pushTestDataToDB(todoRepo)
	require.NoError(suite.T(), err)
}

func (suite *controllerIntegrationTestSuite) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *controllerIntegrationTestSuite) TestGetAllTodosExpectTodosFetched() {
	expectedStatus := http.StatusOK

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	var actualTodos []Todo

	err = json.NewDecoder(res.Body).Decode(&actualTodos)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectTodoFetched() {
	expectedStatus := http.StatusOK
	expectedTodo := suite.testTodos[WillFetchedTodoIdx]

	req, err := http.NewRequest("GET", "/"+expectedTodo.ID.String(), nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	var actualTodo Todo

	err = json.NewDecoder(res.Body).Decode(&actualTodo)
	require.NoError(suite.T(), err)

	expectedTodo.CreatedAt = actualTodo.CreatedAt

	assert.Equal(suite.T(), expectedStatus, actualStatus)
	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectNotFoundReturn() {
	expectedStatus := http.StatusNotFound
	notExistsTodoID := uuid.Nil

	req, err := http.NewRequest("GET", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func (suite *controllerIntegrationTestSuite) TestCreateTodoExpectTodoCreated() {
	expectedStatus := http.StatusCreated
	expectedTodo := Todo{
		Title:    "new title",
		Contents: "new contents",
	}

	todoJSONBytes, err := json.Marshal(expectedTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(todoJSONBytes))
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	var actualTodo Todo

	err = json.NewDecoder(res.Body).Decode(&actualTodo)
	require.NoError(suite.T(), err)

	expectedTodo.ID = actualTodo.ID
	expectedTodo.CreatedAt = actualTodo.CreatedAt

	assert.Equal(suite.T(), expectedStatus, actualStatus)
	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *controllerIntegrationTestSuite) TestCreateTodoExpectBadRequestReturn() {
	expectedStatus := http.StatusBadRequest
	invalidTodo := Todo{
		Title:    "",
		Contents: "new contents",
	}

	todoJSONBytes, err := json.Marshal(invalidTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(todoJSONBytes))
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func (suite *controllerIntegrationTestSuite) TestUpdateTodoByIDExpectTodoUpdated() {
	expectedStatus := http.StatusOK
	expectedTodo := suite.testTodos[WillUpdatedTodoIdx]
	expectedTodo.Title = "updated title"
	expectedTodo.Contents = "updated contents"

	todoJSONBytes, err := json.Marshal(expectedTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+expectedTodo.ID.String(),
		bytes.NewBuffer(todoJSONBytes))
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	var actualTodo Todo

	err = json.NewDecoder(res.Body).Decode(&actualTodo)
	require.NoError(suite.T(), err)

	expectedTodo.CreatedAt = actualTodo.CreatedAt

	assert.Equal(suite.T(), expectedStatus, actualStatus)
	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *controllerIntegrationTestSuite) TestUpdateTodoByIDExpectNotFoundReturn() {
	expectedStatus := http.StatusNotFound
	notExistsTodo := Todo{
		ID:       uuid.Nil,
		Title:    "not exists title",
		Contents: "not exists contents",
	}

	todoJSONBytes, err := json.Marshal(notExistsTodo)
	require.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/"+notExistsTodo.ID.String(),
		bytes.NewBuffer(todoJSONBytes))
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectTodoRemoved() {
	expectedStatus := http.StatusOK
	expectedTodo := suite.testTodos[WillRemovedTodoIdx]

	req, err := http.NewRequest("DELETE", "/"+expectedTodo.ID.String(), nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectNotFoundReturn() {
	expectedStatus := http.StatusNotFound
	notExistsTodoID := uuid.Nil

	req, err := http.NewRequest("DELETE", "/"+notExistsTodoID.String(), nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func getResponse(suite *controllerIntegrationTestSuite, req *http.Request) *http.Response {
	httpRecorder := httptest.NewRecorder()

	suite.ginEngine.ServeHTTP(httpRecorder, req)

	return httpRecorder.Result()
}
