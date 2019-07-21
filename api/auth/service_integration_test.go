package auth_test

import (
	"testing"
	"time"

	"github.com/gghcode/go-gin-starterkit/api/auth"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const existUserID = 10

type serviceIntegration struct {
	suite.Suite

	authService auth.Service
}

func (suite *serviceIntegration) SetupSuite() {
	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	dbConn, err := db.NewConn(conf)
	require.NoError(suite.T(), err)

	redisConn := db.NewRedisConn(conf)
	redisConn.Client().Set(
		auth.RefreshTokenRedisStorageKey(existUserID),
		"exist_refreshtoken",
		30*time.Second,
	)

	suite.authService = auth.NewService(
		conf,
		user.NewRepository(dbConn),
		service.NewPassport(),
		redisConn,
	)
}

func TestAuthServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(serviceIntegration))
}

func (suite *serviceIntegration) TestIssueRefreshToken() {
	testCases := []struct {
		description string
		argsUserID  int64
		expected    bool
	}{
		{
			description: "ShouldIssueRefreshToken",
			argsUserID:  1,
			expected:    true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			refreshToken, _ := suite.authService.IssueRefreshToken(tc.argsUserID)

			actual := suite.authService.VerifyRefreshToken(
				tc.argsUserID,
				refreshToken,
			)

			suite.Equal(tc.expected, actual)
		})
	}
}

func (suite *serviceIntegration) TestVerifyRefreshToken() {
	testCases := []struct {
		description      string
		argsUserID       int64
		argsRefreshToken string
		stub             func(authService auth.Service) (string, error)
		expected         bool
	}{
		{
			description:      "ShouldBeValid",
			argsUserID:       existUserID,
			argsRefreshToken: "exist_refreshtoken",
			expected:         true,
		},
		{
			description:      "ShouldBeInvalid_WhenInvalidRefreshToken",
			argsUserID:       existUserID,
			argsRefreshToken: "NOT_EXIST_REFRESH_TOKEN",
			expected:         false,
		},
		{
			description:      "ShouldBeInvalid_WhenNotExistsUserID",
			argsUserID:       -1,
			argsRefreshToken: "exist_refreshtoken",
			expected:         false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual := suite.authService.VerifyRefreshToken(
				tc.argsUserID,
				tc.argsRefreshToken,
			)

			suite.Equal(tc.expected, actual)
		})
	}
}
