package todo_test

import (
	"encoding/json"

	"io"
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/todo"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/internal/testutil"
	"github.com/gghcode/go-gin-starterkit/middleware"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerIntegration struct {
	suite.Suite

	ginEngine *gin.Engine
	dbConn    *db.Conn

	testTodos []todo.Todo
}

func TestTodoControllerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(controllerIntegration))
}

func (suite *controllerIntegration) SetupSuite() {
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

	todoRepo := todo.NewRepository(dbConn)
	todoController := todo.NewController(todoRepo)
	todoController.RegisterRoutes(suite.ginEngine)

	suite.testTodos, err = pushTestDataToDB(todoRepo)
	require.NoError(suite.T(), err)
}

func (suite *controllerIntegration) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *controllerIntegration) TestGetAllTodos() {
	testCases := []struct {
		description    string
		expectedStatus int
	}{
		{
			description:    "ShouldFetchTodos",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"GET",
				todo.APIPath,
				nil)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.NotNil(actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestGetTodoByID() {
	testCases := []struct {
		description    string
		argsTodoID     string
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldFetchTodo",
			argsTodoID:     suite.testTodos[WillFetchedTodoIdx].ID.String(),
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				suite.testTodos[WillFetchedTodoIdx].TodoResponse()),
		},
		{
			description:    "ShouldReturnNotFoundErr",
			argsTodoID:     uuid.Nil.String(),
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"GET",
				todo.APIPath+tc.argsTodoID,
				nil)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestCreateTodo() {
	testCases := []struct {
		description    string
		reqBody        io.Reader
		expectedStatus int
		expectedJSONFn func(string) string
	}{
		{
			description: "ShouldCreateTodo",
			reqBody: testutil.ReqBodyFromInterface(suite.T(), todo.CreateTodoRequest{
				Title:    "new title",
				Contents: "new contents",
			}),
			expectedStatus: http.StatusCreated,
			expectedJSONFn: func(actualJSON string) string {
				actualTodoRes := TodoResFromJSONString(suite.T(), actualJSON)
				expectedTodoRes := todo.TodoResponse{
					ID:        actualTodoRes.ID,
					Title:     "new title",
					Contents:  "new contents",
					CreatedAt: actualTodoRes.CreatedAt,
				}

				return testutil.JSONStringFromInterface(suite.T(), expectedTodoRes)
			},
		},
		{
			description: "ShouldReturnBadRequestErr_WhenNotExistTitle",
			reqBody: testutil.ReqBodyFromInterface(suite.T(), todo.CreateTodoRequest{
				Title:    "",
				Contents: "new contents",
			}),
			expectedStatus: http.StatusBadRequest,
			expectedJSONFn: func(actualJSON string) string {
				return actualJSON
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"POST",
				todo.APIPath,
				tc.reqBody,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedJSONFn(actualJSON), actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestUpdateTodoByID() {
	testCases := []struct {
		description    string
		argsTodoID     string
		reqBody        io.Reader
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldUpdateTodo",
			argsTodoID:  suite.testTodos[WillUpdatedTodoIdx].ID.String(),
			reqBody: testutil.ReqBodyFromInterface(suite.T(), todo.CreateTodoRequest{
				Title:    "updated title",
				Contents: "updated contents",
			}),
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), todo.TodoResponse{
				ID:        suite.testTodos[WillUpdatedTodoIdx].ID,
				Title:     "updated title",
				Contents:  "updated contents",
				CreatedAt: suite.testTodos[WillUpdatedTodoIdx].TodoResponse().CreatedAt,
			}),
		},
		{
			description: "ShouldReturnNotFoundErr",
			argsTodoID:  todo.EmptyTodo.ID.String(),
			reqBody: testutil.ReqBodyFromInterface(suite.T(), todo.CreateTodoRequest{
				Title:    "updated title",
				Contents: "updated contents",
			}),
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"PUT",
				todo.APIPath+tc.argsTodoID,
				tc.reqBody,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestRemoveTodoByID() {
	testCases := []struct {
		description    string
		argsTodoID     string
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldRemoveTodo",
			argsTodoID:     suite.testTodos[WillRemovedTodoIdx].ID.String(),
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				suite.testTodos[WillRemovedTodoIdx].TodoResponse()),
		},
		{
			description:    "ShouldReturnNotFoundErr",
			argsTodoID:     todo.EmptyTodo.ID.String(),
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"DELETE",
				todo.APIPath+tc.argsTodoID,
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func TodoResFromJSONString(t *testing.T, jsonString string) todo.TodoResponse {
	var result todo.TodoResponse

	err := json.Unmarshal([]byte(jsonString), &result)
	require.NoError(t, err)

	return result
}
