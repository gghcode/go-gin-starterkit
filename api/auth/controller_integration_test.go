package auth_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/api/auth"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/internal/testutil"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerIntegration struct {
	suite.Suite

	ginEngine *gin.Engine
	dbConn    *db.Conn

	testUser user.User
}

func TestAuthControllerIntegration(t *testing.T) {
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
	suite.dbConn = dbConn

	userRepo := user.NewRepository(dbConn)
	passport := service.NewPassport()

	authController := auth.NewController(
		conf,
		userRepo,
		passport,
		db.NewRedisConn(conf),
	)
	authController.RegisterRoutes(suite.ginEngine)

	testUser := user.User{
		UserName: "username",
	}
	testUser.PasswordHash, _ = passport.HashPassword("password")

	suite.testUser, err = userRepo.CreateUser(testUser)
	suite.NoError(err)
}

func (suite *controllerIntegration) TestGetToken() {
	testCases := []struct {
		description    string
		reqPayload     *auth.CreateAccessTokenRequest
		expectedStatus int
	}{
		{
			description: "ShouldGenerateToken",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: suite.testUser.UserName,
				Password: "password",
			},
			expectedStatus: http.StatusOK,
		},
		{
			description:    "ShouldReturnBadReqestErr",
			reqPayload:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "ShouldReturnUnauthorizedErr",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: "NOT_EXISTS_USER_NAME",
				Password: "password",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			description: "ShouldReturnUnauthorizedErr_WhenIncorrectPassword",
			reqPayload: &auth.CreateAccessTokenRequest{
				UserName: suite.testUser.UserName,
				Password: "INCORRECT_PASSWORD",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			reqBody := testutil.ReqBodyFromInterface(suite.T(), tc.reqPayload)

			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"POST",
				auth.APIPath+"token",
				reqBody,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}
