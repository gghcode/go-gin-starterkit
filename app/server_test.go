package app

import (
	"github.com/gin-gonic/gin"
	ht "github.com/gyuhwankim/go-gin-starterkit/app/http"
	"github.com/gyuhwankim/go-gin-starterkit/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPong(t *testing.T) {
	router := setupRouter()
	recorder := performRequest(router, "GET", "/ping")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
}

func TestRegisterResource(t *testing.T) {
	router := setupRouter()

	registerResource(router, "VALID/PATH", []ht.Route{
		ht.Route{
			Method: "GET",
			Path:   "1",
			Handler: func(ctx *gin.Context) {
				ctx.String(200, "VALID")
			},
		},
	})

	recorder := performRequest(router, "GET", "/VALID/PATH/1")

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "VALID", recorder.Body.String())
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
	return New(config.Configuration{}).core
}
