package app

import (
	"github.com/gyuhwankim/go-gin-starterkit/config"
	"github.com/stretchr/testify/assert"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPong(t *testing.T) {
	router := New(config.Configuration{}).core
	recorder := performRequest(router, "GET", "/ping")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
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
