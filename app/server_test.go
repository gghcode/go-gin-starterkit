package app

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/config"
)

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
