package todo

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/app/api/testutil"
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

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "GET", "/", nil)
	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)

	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
	assert.NotNil(suite.T(), actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectTodoFetched() {
	willFetchTodoRes := suite.testTodos[WillFetchedTodoIdx].TodoResponse()

	expectedStatus := http.StatusOK
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), willFetchTodoRes)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "GET", "/"+willFetchTodoRes.ID.String(), nil)
	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)

	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestGetTodoByIDExpectNotFoundReturn() {
	notExistsTodoID := uuid.Nil

	expectedStatus := http.StatusNotFound
	expectedErrRes := common.NewErrResp(common.ErrEntityNotFound)
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), expectedErrRes)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "GET", "/"+notExistsTodoID.String(), nil)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

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

	reqBody := testutil.ReqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "POST", "/", reqBody)
	actualTodoRes := todoResFromResBody(suite.T(), actualRes.Body)
	actualJSON := testutil.JSONStringFromInterface(suite.T(), actualTodoRes)

	expectedTodoRes.ID = actualTodoRes.ID
	expectedTodoRes.CreatedAt = actualTodoRes.CreatedAt
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), expectedTodoRes)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestCreateTodoExpectBadRequestReturn() {
	invalidCreateTodoReq := CreateTodoRequest{
		Title:    "",
		Contents: "new contents",
	}

	expectedStatus := http.StatusBadRequest

	reqBody := testutil.ReqBodyFromInterface(suite.T(), invalidCreateTodoReq)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "POST", "/", reqBody)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

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
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), willUpdateTodoRes)

	reqBody := testutil.ReqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "PUT", "/"+willUpdateTodoRes.ID.String(), reqBody)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

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
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), expectedErrRes)

	reqBody := testutil.ReqBodyFromInterface(suite.T(), createTodoReq)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "PUT", "/"+notExistsTodoID.String(), reqBody)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectTodoRemoved() {
	willRemoveTodoRes := suite.testTodos[WillRemovedTodoIdx].TodoResponse()

	expectedStatus := http.StatusOK
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), willRemoveTodoRes)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "DELETE", "/"+willRemoveTodoRes.ID.String(), nil)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func (suite *controllerIntegrationTestSuite) TestRemoveTodoByIDExpectNotFoundReturn() {
	notExistsTodoID := uuid.Nil

	expectedStatus := http.StatusNotFound
	expectedErrRes := common.NewErrResp(common.ErrEntityNotFound)
	expectedJSON := testutil.JSONStringFromInterface(suite.T(), expectedErrRes)

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "DELETE", "/"+notExistsTodoID.String(), nil)
	actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
	assert.Equal(suite.T(), expectedJSON, actualJSON)
}

func todoResFromResBody(t *testing.T, body io.Reader) TodoResponse {
	var result TodoResponse

	err := json.NewDecoder(body).Decode(&result)
	require.NoError(t, err)

	return result
}
