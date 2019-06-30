package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type controllerUnitTestSuite struct {
	suite.Suite

	ginEngine  *gin.Engine
	controller *Controller
}

func TestCommonControllerUnit(t *testing.T) {
	suite.Run(t, new(controllerUnitTestSuite))
}

func (suite *controllerUnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.ginEngine = gin.New()
	suite.controller = NewController()
	suite.controller.RegisterRoutes(suite.ginEngine)
}

func (suite *controllerUnitTestSuite) TestHealthyExpectedStatusOK() {
	expectedStatus := http.StatusOK

	req, err := http.NewRequest("GET", "/healthy", nil)
	require.NoError(suite.T(), err)

	res := getResponse(suite, req)
	actualStatus := res.StatusCode

	assert.Equal(suite.T(), expectedStatus, actualStatus)
}

func getResponse(suite *controllerUnitTestSuite, req *http.Request) *http.Response {
	httpRecorder := httptest.NewRecorder()

	suite.ginEngine.ServeHTTP(httpRecorder, req)

	return httpRecorder.Result()
}
