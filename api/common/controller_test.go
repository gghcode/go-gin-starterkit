package common

import (
	"net/http"
	"testing"

	"github.com/gghcode/go-gin-starterkit/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

	actualRes := testutil.ActualResponse(suite.T(), suite.ginEngine, "GET", "/healthy", nil)

	assert.Equal(suite.T(), expectedStatus, actualRes.StatusCode)
}
