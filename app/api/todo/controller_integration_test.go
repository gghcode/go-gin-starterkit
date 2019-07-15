package todo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/middleware"

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

	ginEngine *gin.Engine
	dbConn    *db.Conn

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

	suite.ginEngine = gin.New()
	suite.ginEngine.Use(func(ctx *gin.Context) {
		var innerHandler gin.HandlerFunc = func(ctx *gin.Context) {}

		ctx.Set(middleware.VerifyHandlerKey, innerHandler)
		ctx.Next()
	})

	suite.dbConn = dbConn

	todoRepo := NewRepository(dbConn)
	todoController := NewController(todoRepo)
	todoController.RegisterRoutes(suite.ginEngine)

	suite.testTodos, err = pushTestDataToDB(todoRepo)
	require.NoError(suite.T(), err)
}

func (suite *controllerIntegrationTestSuite) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *controllerIntegrationTestSuite) TestGetAllTodosExpectTodosFetched() {
	expectedStatus := http.StatusOK

	actualRes := actualResponse(suite, "GET", "/", nil)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.NotNil(suite.T(), actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectTodoFetched() {
	willFetchTodoRes := suite.testTodos[WillFetchedTodoIdx].TodoResponse()

	expectedStatus := http.StatusOK
	expectedJSON := jsonStringFromInterface(suite.T(), willFetchTodoRes)

	actualRes := actualResponse(suite, "GET", "/"+willFetchTodoRes.ID.String(), nil)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectNotFoundReturn() {
	notExistsTodoID := uuid.Nil

	expectedStatus := http.StatusNotFound
	expectedErrRes := common.NewErrResp(common.ErrEntityNotFound)
	expectedJSON := jsonStringFromInterface(suite.T(), expectedErrRes)

	actualRes := actualResponse(suite, "GET", "/"+notExistsTodoID.String(), nil)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestCreateTodoExpectTodoCreated() {
	createTodoReq := CreateTodoRequest{
		Title:    "new title",
		Contents: "new contents",
	}

	expectedStatus := http.StatusCreated
	expectedTodoRes := TodoResponse{
		Title:    createTodoReq.Title,
		Contents: createTodoReq.Contents,
	}

	reqBody := reqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := actualResponse(suite, "POST", "/", reqBody)
	actualTodoRes := todoResFromResBody(suite.T(), actualRes.Body)
	actualJSON := jsonStringFromInterface(suite.T(), actualTodoRes)

	expectedTodoRes.ID = actualTodoRes.ID
	expectedTodoRes.CreatedAt = actualTodoRes.CreatedAt
	expectedJSON := jsonStringFromInterface(suite.T(), expectedTodoRes)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestCreateTodoExpectBadRequestReturn() {
	invalidCreateTodoReq := CreateTodoRequest{
		Title:    "",
		Contents: "new contents",
	}

	expectedStatus := http.StatusBadRequest

	reqBody := reqBodyFromInterface(suite.T(), invalidCreateTodoReq)

	actualRes := actualResponse(suite, "POST", "/", reqBody)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.NotNil(suite.T(), actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestUpdateTodoByIDExpectTodoUpdated() {
	createTodoReq := CreateTodoRequest{
		Title:    "updated title",
		Contents: "updated contents",
	}

	willUpdateTodoRes := suite.testTodos[WillUpdatedTodoIdx].TodoResponse()
	willUpdateTodoRes.Title = createTodoReq.Title
	willUpdateTodoRes.Contents = createTodoReq.Contents

	expectedStatus := http.StatusOK
	expectedJSON := jsonStringFromInterface(suite.T(), willUpdateTodoRes)

	reqBody := reqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := actualResponse(suite, "PUT", "/"+willUpdateTodoRes.ID.String(), reqBody)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestUpdateTodoByIDExpectNotFoundReturn() {
	notExistsTodoID := uuid.Nil
	createTodoReq := CreateTodoRequest{
		Title:    "updated title",
		Contents: "updated contents",
	}

	expectedStatus := http.StatusNotFound
	expectedErrRes := common.NewErrResp(common.ErrEntityNotFound)
	expectedJSON := jsonStringFromInterface(suite.T(), expectedErrRes)

	reqBody := reqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := actualResponse(suite, "PUT", "/"+notExistsTodoID.String(), reqBody)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectTodoRemoved() {
	willRemoveTodoRes := suite.testTodos[WillRemovedTodoIdx].TodoResponse()

	expectedStatus := http.StatusOK
	expectedJSON := jsonStringFromInterface(suite.T(), willRemoveTodoRes)

	actualRes := actualResponse(suite, "DELETE", "/"+willRemoveTodoRes.ID.String(), nil)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectNotFoundReturn() {
	notExistsTodoID := uuid.Nil

	expectedStatus := http.StatusNotFound
	expectedErrRes := common.NewErrResp(common.ErrEntityNotFound)
	expectedJSON := jsonStringFromInterface(suite.T(), expectedErrRes)

	actualRes := actualResponse(suite, "DELETE", "/"+notExistsTodoID.String(), nil)
	actualJSON := jsonStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func actualResponse(suite *controllerIntegrationTestSuite,
	method, url string, body io.Reader) *http.Response {
	httpRecorder := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, body)
	require.NoError(suite.T(), err)

	suite.ginEngine.ServeHTTP(httpRecorder, req)

	return httpRecorder.Result()
}

func reqBodyFromInterface(t *testing.T, body interface{}) *bytes.Buffer {
	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)

	return bytes.NewBuffer(jsonBytes)
}

func jsonStringFromInterface(t *testing.T, res interface{}) string {
	bytes, err := json.Marshal(res)
	require.NoError(t, err)

	return string(bytes)
}

func jsonStringFromResBody(t *testing.T, body io.Reader) string {
	bytes, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	return string(bytes)
}

func todoResFromResBody(t *testing.T, body io.Reader) TodoResponse {
	var result TodoResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
