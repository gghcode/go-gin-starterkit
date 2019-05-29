package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/config"
	"github.com/stretchr/testify/assert"
)

func TestHealthyHandler(t *testing.T) {
	router := setupRouter()
	recorder := performRequest(router, "GET", "/healthy")

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func performRequest(router http.Handler, method, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func setupRouter() *gin.Engine {
	return NewServer(config.Configuration{}).core
}
