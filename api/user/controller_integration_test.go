package user_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gghcode/go-gin-starterkit/middleware"
	"github.com/gghcode/go-gin-starterkit/service"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/testutil"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerIntegration struct {
	suite.Suite

	ginEngine *gin.Engine
	dbConn    *db.Conn

	testUsers []user.User
}

func TestUserControllerIntegration(t *testing.T) {
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

	userRepo := user.NewRepository(dbConn)
	userController := user.NewController(userRepo, service.NewPassport())
	userController.RegisterRoutes(suite.ginEngine)

	suite.testUsers, err = pushTestDataToDB(userRepo, "controller")
	require.NoError(suite.T(), err)
}

func (suite *controllerIntegration) TestCreateUser() {
	testCases := []struct {
		description    string
		createUserReq  *user.CreateUserRequest
		expectedStatus int
		expectedJSON   func(string) string
	}{
		{
			description: "ShouldBeCreated",
			createUserReq: &user.CreateUserRequest{
				UserName: "New User",
				Password: "New Password",
			},
			expectedStatus: http.StatusCreated,
			expectedJSON: func(actualJSON string) string {
				actualUserRes := UserResFromJSONString(suite.T(), actualJSON)
				expectedUserRes := user.UserResponse{
					ID:        actualUserRes.ID,
					UserName:  "New User",
					CreatedAt: actualUserRes.CreatedAt,
				}

				return testutil.JSONStringFromInterface(suite.T(), expectedUserRes)
			},
		},
		{
			description: "ShouldBeConflictWhenAlreadyExistUser",
			createUserReq: &user.CreateUserRequest{
				UserName: suite.testUsers[WillFetchedEntityIdx].UserName,
				Password: "New Password",
			},
			expectedStatus: http.StatusConflict,
			expectedJSON: func(string) string {
				return testutil.JSONStringFromInterface(suite.T(),
					common.NewErrResp(common.ErrAlreadyExistsEntity))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			reqBody := testutil.ReqBodyFromInterface(suite.T(), tc.createUserReq)

			actualRes := testutil.ActualResponse(suite.T(),
				suite.ginEngine, "POST", "/users/", reqBody)
			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
			suite.JSONEq(tc.expectedJSON(actualJSON), actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestGetUserByName() {
	testCases := []struct {
		description    string
		username       string
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldBeOK",
			username:       suite.testUsers[WillFetchedEntityIdx].UserName,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				suite.testUsers[WillFetchedEntityIdx].Response()),
		},
		{
			description:    "ShouldBeNotFound",
			username:       "NOT_EXISTS_USER_NAME",
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine,
				"GET", "/users/"+tc.username, nil)
			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *controllerIntegration) TestUpdateUserByID() {
	testCases := []struct {
		description    string
		userID         string
		updateUserReq  *user.UpdateUserRequest
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldBeOK",
			userID:         strconv.FormatInt(suite.testUsers[WillUpdatedEntityIdx].ID, 10),
			updateUserReq:  &user.UpdateUserRequest{UserName: "updated_username"},
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), user.UserResponse{
				ID:        suite.testUsers[WillUpdatedEntityIdx].ID,
				UserName:  "updated_username",
				CreatedAt: suite.testUsers[WillUpdatedEntityIdx].Response().CreatedAt,
			}),
		},
		{
			description:    "ShouldBeBadRequest_WhenNotContainUserName",
			userID:         strconv.FormatInt(suite.testUsers[WillUpdatedEntityIdx].ID, 10),
			updateUserReq:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   "",
		},
		{
			description:    "ShouldBeBadRequest_WhenInvalidUserID",
			userID:         "NOT_INTEGER_USER_ID",
			updateUserReq:  &user.UpdateUserRequest{UserName: "username"},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrParsingFailed)),
		},
		{
			description:    "ShouldBeBadRequest_WhenLessUserNameLenMin4",
			userID:         strconv.FormatInt(suite.testUsers[WillUpdatedEntityIdx].ID, 10),
			updateUserReq:  &user.UpdateUserRequest{UserName: "use"},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   "",
		},
		{
			description:    "ShouldBeNotFound_WhenNotExistsEntity",
			userID:         strconv.FormatInt(user.EmptyUser.ID, 10),
			updateUserReq:  &user.UpdateUserRequest{UserName: "username100"},
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			reqBody := testutil.ReqBodyFromInterface(suite.T(), tc.updateUserReq)

			actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine,
				"PUT", "/users/"+tc.userID, reqBody)
			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
			if tc.expectedJSON != "" {
				suite.JSONEq(tc.expectedJSON, actualJSON)
			}
		})
	}
}

func (suite *controllerIntegration) TestRemoveUser() {
	testCases := []struct {
		description    string
		userID         string
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldBeOK",
			userID:         strconv.FormatInt(suite.testUsers[WillRemovedEntityIdx].ID, 10),
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				suite.testUsers[WillRemovedEntityIdx].Response()),
		},
		{
			description:    "ShouldBeBadRequest_WhenInvalidUserID",
			userID:         "INVALID_STRING_USER_ID",
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrParsingFailed)),
		},
		{
			description:    "ShouldBeNotFound_WhenNotExistEntity",
			userID:         strconv.FormatInt(user.EmptyUser.ID, 10),
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(common.ErrEntityNotFound)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine,
				"DELETE", "/users/"+tc.userID, nil)
			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}

func UserResFromJSONString(t *testing.T, jsonString string) user.UserResponse {
	var result user.UserResponse

	err := json.Unmarshal([]byte(jsonString), &result)
	require.NoError(t, err)

	return result
}
