package common

import (
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type controllerUnit struct {
	suite.Suite

	ginEngine  *gin.Engine
	controller *Controller
}

func TestCommonControllerUnit(t *testing.T) {
	suite.Run(t, new(controllerUnit))
}

func (suite *controllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.ginEngine = gin.New()
	suite.controller = NewController()
	suite.controller.RegisterRoutes(suite.ginEngine)
}

func (suite *controllerUnit) TestHealthy() {
	testCases := []struct {
		description    string
		expectedStatus int
	}{
		{
			description:    "ShouldReturnOK",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.ActualResponse(
				suite.T(),
				suite.ginEngine,
				"GET",
				"/healthy",
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}
